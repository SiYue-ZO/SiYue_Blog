package persistence

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"siyue-blog/internal/entity"
	"siyue-blog/migrations"
)

func TestPostViewCountingPostgres(t *testing.T) {
	dsn := os.Getenv("SIYUE_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("set SIYUE_TEST_POSTGRES_DSN to a dedicated PostgreSQL test database")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	pool, err := db.DB()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, pool.Close()) })
	ctx := context.Background()
	require.NoError(t, migrations.Up(ctx, db))
	name := fmt.Sprintf("view-test-%d", time.Now().UnixNano())
	var authorID uint64
	require.NoError(t, db.Raw(`INSERT INTO users (username, email, password_hash)
		VALUES (?, ?, 'test') RETURNING id`, name, name+"@example.test").Scan(&authorID).Error)
	t.Cleanup(func() { require.NoError(t, db.Exec("DELETE FROM users WHERE id = ?", authorID).Error) })
	now := time.Now().Add(-time.Hour)
	post := &entity.Post{
		Title: name, Slug: name, Content: "body", Status: entity.PostStatusPublished,
		AuthorID: authorID, PublishedAt: &now, ViewCount: 10,
	}
	repo := NewPostRepo(db)
	require.NoError(t, repo.Create(ctx, post))
	t.Cleanup(func() { require.NoError(t, db.Exec("DELETE FROM posts WHERE id = ?", post.ID).Error) })
	initial, err := repo.GetByID(ctx, post.ID)
	require.NoError(t, err)
	eventID := "12345678-1234-4234-8234-123456789012"

	countConcurrently := func(t *testing.T, ids []string) {
		t.Helper()
		var wg sync.WaitGroup
		errors := make(chan error, len(ids))
		for _, id := range ids {
			wg.Go(func() {
				_, err := repo.RecordView(ctx, post.ID, id)
				errors <- err
			})
		}
		wg.Wait()
		close(errors)
		for err := range errors {
			require.NoError(t, err)
		}
	}

	t.Run("duplicate concurrent events count only once", func(t *testing.T) {
		ids := make([]string, 32)
		for i := range ids {
			ids[i] = eventID
		}
		countConcurrently(t, ids)
		count, err := NewPostRepo(db).RecordView(ctx, post.ID, eventID)
		require.NoError(t, err)
		require.Equal(t, 11, count)
		current, err := repo.GetByID(ctx, post.ID)
		require.NoError(t, err)
		require.True(t, initial.UpdatedAt.Equal(current.UpdatedAt), "reading must not change updated_at")
	})

	stale, err := repo.GetByID(ctx, post.ID)
	require.NoError(t, err)
	t.Run("different concurrent events each count once", func(t *testing.T) {
		ids := make([]string, 24)
		for i := range ids {
			ids[i] = fmt.Sprintf("%08x-1234-4234-8234-123456789012", i)
		}
		countConcurrently(t, ids)
		current, err := repo.GetByID(ctx, post.ID)
		require.NoError(t, err)
		require.Equal(t, 35, current.ViewCount)
	})

	t.Run("editing with tags cannot overwrite new views", func(t *testing.T) {
		tag := entity.Tag{Name: name, Slug: name}
		require.NoError(t, db.Create(&tag).Error)
		t.Cleanup(func() { require.NoError(t, db.Exec("DELETE FROM tags WHERE id = ?", tag.ID).Error) })
		stale.Title = "edited title"
		require.NoError(t, repo.UpdateWithTags(ctx, stale, []entity.Tag{tag}, true))
		current, err := repo.GetByID(ctx, post.ID)
		require.NoError(t, err)
		require.Equal(t, 35, current.ViewCount)
		require.Equal(t, "edited title", current.Title)
		require.Len(t, current.Tags, 1)
	})

	t.Run("failed increment rolls back event and permits retry", func(t *testing.T) {
		id := "ffffffff-1234-4234-8234-123456789012"
		require.NoError(t, db.Exec("UPDATE posts SET view_count = 2147483647 WHERE id = ?", post.ID).Error)
		_, err := repo.RecordView(ctx, post.ID, id)
		require.Error(t, err)
		var count int64
		require.NoError(t, db.Table("post_view_events").Where("post_id = ? AND event_id = ?", post.ID, id).Count(&count).Error)
		require.Zero(t, count)
		require.NoError(t, db.Exec("UPDATE posts SET view_count = 35 WHERE id = ?", post.ID).Error)
		value, err := repo.RecordView(ctx, post.ID, id)
		require.NoError(t, err)
		require.Equal(t, 36, value)
	})

	t.Run("unpublished and future articles cannot count", func(t *testing.T) {
		require.NoError(t, db.Exec("UPDATE posts SET status = 'draft' WHERE id = ?", post.ID).Error)
		_, err := repo.RecordView(ctx, post.ID, eventID)
		require.ErrorIs(t, err, gorm.ErrRecordNotFound)
		require.NoError(t, db.Exec("UPDATE posts SET status = 'published', published_at = NOW() + INTERVAL '1 day' WHERE id = ?", post.ID).Error)
		_, err = repo.RecordView(ctx, post.ID, eventID)
		require.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("deleting post also removes events", func(t *testing.T) {
		require.NoError(t, repo.Delete(ctx, post.ID))
		var count int64
		require.NoError(t, db.Table("post_view_events").Where("post_id = ?", post.ID).Count(&count).Error)
		require.Zero(t, count)
	})
}
