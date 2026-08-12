package v1

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"siyue-blog/internal/entity"
)

type homePageUseCaseStub struct {
	content *entity.HomePageContent
	userID  uint64
}

func (stub *homePageUseCaseStub) GetPublished(context.Context) (*entity.HomePageConfig, error) {
	return &stub.content.PublishedConfig, nil
}
func (stub *homePageUseCaseStub) GetAdmin(context.Context) (*entity.HomePageContent, error) {
	return stub.content, nil
}
func (stub *homePageUseCaseStub) SaveDraft(_ context.Context, config entity.HomePageConfig, userID uint64) (*entity.HomePageContent, error) {
	stub.userID = userID
	stub.content.DraftConfig = config
	return stub.content, nil
}
func (stub *homePageUseCaseStub) Publish(_ context.Context, userID uint64) (*entity.HomePageContent, error) {
	stub.userID = userID
	return stub.content, nil
}
func (stub *homePageUseCaseStub) ResetDraft(_ context.Context, userID uint64) (*entity.HomePageContent, error) {
	stub.userID = userID
	return stub.content, nil
}

func TestHomePagePublishedDoesNotExposeDraft(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &homePageUseCaseStub{content: &entity.HomePageContent{
		DraftConfig:     entity.HomePageConfig{SchemaVersion: 99},
		PublishedConfig: entity.HomePageConfig{SchemaVersion: 1},
	}}
	router := gin.New()
	router.GET("/home", NewHomePageHandler(stub).GetPublished)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/home", nil))

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Contains(t, recorder.Body.String(), `"schema_version":1`)
	require.NotContains(t, recorder.Body.String(), `"draft_config"`)
}

func TestSaveHomePageDraftUsesAuthenticatedUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &homePageUseCaseStub{content: &entity.HomePageContent{}}
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("user_id", uint64(42)) })
	router.PUT("/home/draft", NewHomePageHandler(stub).SaveDraft)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/home/draft", bytes.NewBufferString(`{"schema_version":1,"order":[],"sections":{}}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, uint64(42), stub.userID)
}
