package migrations

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"sort"
	"strings"

	"gorm.io/gorm"
)

//go:embed *.up.sql
var migrationFiles embed.FS

type legacyRequirement struct {
	version    string
	relations  map[string][]string
	indexes    []string
	types      []string
	triggers   []string
	functions  []string
	extensions []string
}

var legacyRequirements = []legacyRequirement{
	{
		version:   "001_create_users.up.sql",
		relations: map[string][]string{"users": {"id", "username", "email", "password_hash", "avatar_url", "role", "status", "created_at", "updated_at"}},
		indexes:   []string{"idx_users_email", "idx_users_username"}, types: []string{"user_status", "post_status"},
		triggers: []string{"trg_users_updated_at"}, functions: []string{"update_updated_at_column"}, extensions: []string{"pg_trgm"},
	},
	{
		version:   "002_create_posts.up.sql",
		relations: map[string][]string{"posts": {"id", "title", "slug", "excerpt", "content", "cover_image", "status", "view_count", "is_pinned", "author_id", "published_at", "created_at", "updated_at"}},
		indexes:   []string{"idx_posts_slug", "idx_posts_status_published_at", "idx_posts_author", "idx_posts_search"}, triggers: []string{"trg_posts_updated_at"},
	},
	{
		version: "003_create_tags.up.sql",
		relations: map[string][]string{
			"tags":      {"id", "name", "slug", "created_at"},
			"post_tags": {"post_id", "tag_id"},
		},
		indexes: []string{"idx_tags_slug", "idx_tags_name", "post_tags_pkey", "idx_post_tags_tag"},
	},
	{
		version:   "004_create_projects.up.sql",
		relations: map[string][]string{"projects": {"id", "title", "description", "repo_url", "demo_url", "image_url", "tech_stack", "is_featured", "sort_order", "created_at", "updated_at"}},
		indexes:   []string{"idx_projects_featured"}, triggers: []string{"trg_projects_updated_at"},
	},
}

func exists(ctx context.Context, conn *sql.Conn, query string, value string) (bool, error) {
	var found bool
	if err := conn.QueryRowContext(ctx, query, value).Scan(&found); err != nil {
		return false, err
	}
	return found, nil
}

func baselineLegacySchema(ctx context.Context, conn *sql.Conn) error {
	for _, requirement := range legacyRequirements {
		complete := true
		for relation, columns := range requirement.relations {
			found, err := exists(ctx, conn, `SELECT to_regclass($1) IS NOT NULL`, "public."+relation)
			if err != nil {
				return err
			}
			complete = complete && found
			for _, column := range columns {
				var columnFound bool
				if err := conn.QueryRowContext(ctx, `SELECT EXISTS (
					SELECT 1 FROM information_schema.columns
					WHERE table_schema = 'public' AND table_name = $1 AND column_name = $2
				)`, relation, column).Scan(&columnFound); err != nil {
					return err
				}
				complete = complete && columnFound
			}
		}
		for _, index := range requirement.indexes {
			found, err := exists(ctx, conn, `SELECT to_regclass($1) IS NOT NULL`, "public."+index)
			if err != nil {
				return err
			}
			complete = complete && found
		}
		for _, typeName := range requirement.types {
			found, err := exists(ctx, conn, `SELECT EXISTS (
				SELECT 1 FROM pg_type t JOIN pg_namespace n ON n.oid = t.typnamespace
				WHERE n.nspname = 'public' AND t.typname = $1
			)`, typeName)
			if err != nil {
				return err
			}
			complete = complete && found
		}
		for _, trigger := range requirement.triggers {
			found, err := exists(ctx, conn, `SELECT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = $1 AND NOT tgisinternal)`, trigger)
			if err != nil {
				return err
			}
			complete = complete && found
		}
		for _, function := range requirement.functions {
			found, err := exists(ctx, conn, `SELECT EXISTS (
				SELECT 1 FROM pg_proc p JOIN pg_namespace n ON n.oid = p.pronamespace
				WHERE n.nspname = 'public' AND p.proname = $1
			)`, function)
			if err != nil {
				return err
			}
			complete = complete && found
		}
		for _, extension := range requirement.extensions {
			found, err := exists(ctx, conn, `SELECT EXISTS (SELECT 1 FROM pg_extension WHERE extname = $1)`, extension)
			if err != nil {
				return err
			}
			complete = complete && found
		}
		if complete {
			if _, err := conn.ExecContext(ctx, `INSERT INTO schema_migrations (version, checksum) VALUES ($1, '') ON CONFLICT DO NOTHING`, requirement.version); err != nil {
				return err
			}
		}
	}
	return nil
}

// Up applies embedded migrations exactly once. A PostgreSQL advisory lock
// makes startup safe when more than one application replica starts at once.
func Up(ctx context.Context, db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get database handle: %w", err)
	}
	conn, err := sqlDB.Conn(ctx)
	if err != nil {
		return fmt.Errorf("reserve migration connection: %w", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	if _, err := conn.ExecContext(ctx, `SELECT pg_advisory_lock(7411052026)`); err != nil {
		return fmt.Errorf("acquire migration lock: %w", err)
	}
	defer func() { _, _ = conn.ExecContext(context.Background(), `SELECT pg_advisory_unlock(7411052026)`) }()

	if _, err := conn.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			checksum VARCHAR(64) NOT NULL DEFAULT '',
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		ALTER TABLE schema_migrations ADD COLUMN IF NOT EXISTS checksum VARCHAR(64) NOT NULL DEFAULT ''
	`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}
	if err := baselineLegacySchema(ctx, conn); err != nil {
		return fmt.Errorf("baseline legacy schema: %w", err)
	}

	entries, err := migrationFiles.ReadDir(".")
	if err != nil {
		return fmt.Errorf("read embedded migrations: %w", err)
	}
	var names []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)

	for _, name := range names {
		contents, err := migrationFiles.ReadFile(name)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}
		checksum := fmt.Sprintf("%x", sha256.Sum256(contents))
		var storedChecksum string
		err = conn.QueryRowContext(ctx, `SELECT checksum FROM schema_migrations WHERE version = $1`, name).Scan(&storedChecksum)
		if err == nil {
			if storedChecksum == "" {
				if _, err := conn.ExecContext(ctx, `UPDATE schema_migrations SET checksum = $2 WHERE version = $1`, name, checksum); err != nil {
					return fmt.Errorf("record baseline checksum %s: %w", name, err)
				}
			} else if storedChecksum != checksum {
				return fmt.Errorf("migration %s checksum changed after it was applied", name)
			}
			continue
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("check migration %s: %w", name, err)
		}
		tx, err := conn.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin migration %s: %w", name, err)
		}
		if _, err := tx.ExecContext(ctx, string(contents)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("apply migration %s: %w", name, err)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations (version, checksum) VALUES ($1, $2)`, name, checksum); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration %s: %w", name, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", name, err)
		}
	}
	return nil
}
