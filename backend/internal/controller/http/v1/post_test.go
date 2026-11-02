package v1

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"siyue-blog/internal/entity"
	"siyue-blog/internal/usecase"
)

type postUseCaseStub struct {
	usecase.PostUseCase
	post          *entity.Post
	incrementErr  error
	incrementedID uint64
}

func (stub *postUseCaseStub) GetBySlug(context.Context, string) (*entity.Post, error) {
	return stub.post, nil
}

func (stub *postUseCaseStub) IncrementViewCount(_ context.Context, id uint64) (int, error) {
	stub.incrementedID = id
	if stub.incrementErr != nil {
		return 0, stub.incrementErr
	}
	return stub.post.ViewCount + 1, nil
}

func TestGetBySlugReturnsIncrementedViewCount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &postUseCaseStub{post: &entity.Post{ID: 12, Slug: "test-post", ViewCount: 6}}
	router := gin.New()
	router.GET("/posts/:slug", NewPostHandler(stub, zerolog.Nop()).GetBySlug)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/posts/test-post", nil))

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, uint64(12), stub.incrementedID)
	var response struct {
		Data publicPost `json:"data"`
	}
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	require.Equal(t, 7, response.Data.ViewCount)
}

func TestGetBySlugKeepsReturningPostWhenViewCountUpdateFails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var logs strings.Builder
	stub := &postUseCaseStub{
		post:         &entity.Post{ID: 12, Slug: "test-post", ViewCount: 6},
		incrementErr: errors.New("database unavailable"),
	}
	router := gin.New()
	router.GET("/posts/:slug", NewPostHandler(stub, zerolog.New(&logs)).GetBySlug)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/posts/test-post", nil))

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Contains(t, recorder.Body.String(), `"view_count":6`)
	require.Contains(t, logs.String(), "failed to increment post view count")
	require.Contains(t, logs.String(), "database unavailable")
}
