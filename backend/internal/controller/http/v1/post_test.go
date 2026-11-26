package v1

import (
	"context"
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
	getErr        error
	recordErr     error
	recordedID    uint64
	recordedEvent string
	recordCalls   int
}

func (stub *postUseCaseStub) GetBySlug(context.Context, string) (*entity.Post, error) {
	return stub.post, stub.getErr
}

func (stub *postUseCaseStub) RecordView(_ context.Context, id uint64, eventID string) (int, error) {
	stub.recordedID, stub.recordedEvent = id, eventID
	stub.recordCalls++
	return 42, stub.recordErr
}

func TestGetBySlugNeverCountsViews(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &postUseCaseStub{post: &entity.Post{ID: 12, Slug: "test-post", ViewCount: 6}}
	router := gin.New()
	router.GET("/posts/:slug", NewPostHandler(stub, zerolog.Nop()).GetBySlug)
	// Metadata, page rendering and prefetches must all be pure reads.
	for range 3 {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/posts/test-post", nil))
		require.Equal(t, http.StatusOK, recorder.Code)
		require.Contains(t, recorder.Body.String(), `"view_count":6`)
	}
	require.Zero(t, stub.recordCalls)
}

func TestRecordViewReturnsDatabaseCount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &postUseCaseStub{post: &entity.Post{ID: 12, Slug: "test-post", ViewCount: 6}}
	router := gin.New()
	router.POST("/posts/:slug/views", NewPostHandler(stub, zerolog.Nop()).RecordView)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/posts/test-post/views",
		strings.NewReader(`{"event_id":"12345678-1234-4234-8234-123456789012"}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "no-store", recorder.Header().Get("Cache-Control"))
	require.Contains(t, recorder.Body.String(), `"view_count":42`)
	require.Equal(t, uint64(12), stub.recordedID)
	require.Equal(t, "12345678-1234-4234-8234-123456789012", stub.recordedEvent)
	require.Equal(t, 1, stub.recordCalls)
}

func TestRecordViewRejectsInvalidEvent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &postUseCaseStub{}
	router := gin.New()
	router.POST("/posts/:slug/views", NewPostHandler(stub, zerolog.Nop()).RecordView)
	for _, body := range []string{`{}`, `{"event_id":"not-a-uuid"}`, `{"event_id":123}`, `not-json`} {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/posts/test-post/views", strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(recorder, request)
		require.Equal(t, http.StatusBadRequest, recorder.Code)
	}
	require.Zero(t, stub.recordCalls)
}

func TestRecordViewErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, test := range []struct {
		name              string
		getErr, recordErr error
		status, calls     int
	}{
		{"missing or unpublished", errors.New("not found"), nil, http.StatusNotFound, 0},
		{"database unavailable", nil, errors.New("database unavailable"), http.StatusInternalServerError, 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			var logs strings.Builder
			stub := &postUseCaseStub{post: &entity.Post{ID: 12}, getErr: test.getErr, recordErr: test.recordErr}
			router := gin.New()
			router.POST("/posts/:slug/views", NewPostHandler(stub, zerolog.New(&logs)).RecordView)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/posts/test-post/views",
				strings.NewReader(`{"event_id":"12345678-1234-4234-8234-123456789012"}`))
			request.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(recorder, request)
			require.Equal(t, test.status, recorder.Code)
			require.Equal(t, test.calls, stub.recordCalls)
			if test.recordErr != nil {
				require.Contains(t, logs.String(), "failed to record post view")
			}
		})
	}
}
