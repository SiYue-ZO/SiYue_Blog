package v1

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"siyue-blog/internal/entity"
)

type settingsUseCaseStub struct {
	updatedVisible bool
	updatedLinks   []entity.SocialLink
	updatedFooter  *entity.FooterConfig
}

func (stub *settingsUseCaseStub) Get(context.Context) (*entity.SiteSettings, error) {
	return &entity.SiteSettings{ID: 1, ProjectPageVisible: true}, nil
}

func (stub *settingsUseCaseStub) Update(_ context.Context, visible *bool, links *[]entity.SocialLink, footer *entity.FooterConfig) (*entity.SiteSettings, error) {
	if visible != nil {
		stub.updatedVisible = *visible
	}
	if links != nil {
		stub.updatedLinks = *links
	}
	stub.updatedFooter = footer
	return &entity.SiteSettings{ID: 1, ProjectPageVisible: stub.updatedVisible, SocialLinks: links, Footer: footer}, nil
}

func TestUpdateSettingsAcceptsFooter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	settingsUC := &settingsUseCaseStub{}
	handler := NewSettingsHandler(settingsUC)
	router := gin.New()
	router.PATCH("/settings", handler.Update)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPatch, "/settings", bytes.NewBufferString(`{"footer":{"enabled":true,"brand":"我的博客","description":"简介","copyright":"© {year}","show_social_links":true,"filing_icon_url":"/images/beian.png","icp_number":"粤ICP备12345678号","icp_url":"https://beian.miit.gov.cn/","public_security_number":"粤公网安备 44000000000000号","public_security_url":"https://www.beian.gov.cn/","filing_information":"备案信息"}}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if settingsUC.updatedFooter == nil || settingsUC.updatedFooter.Brand != "我的博客" {
		t.Fatalf("expected footer to be updated, got %#v", settingsUC.updatedFooter)
	}
	if settingsUC.updatedFooter.ICPNumber != "粤ICP备12345678号" {
		t.Fatalf("expected filing settings to be updated, got %#v", settingsUC.updatedFooter)
	}
}

func TestUpdateSettingsAcceptsSocialLinks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	settingsUC := &settingsUseCaseStub{}
	handler := NewSettingsHandler(settingsUC)
	router := gin.New()
	router.PATCH("/settings", handler.Update)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPatch, "/settings", bytes.NewBufferString(`{"social_links":[{"platform":"bilibili","label":"哔哩哔哩","href":"https://space.bilibili.com/1"}]}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if len(settingsUC.updatedLinks) != 1 || settingsUC.updatedLinks[0].Platform != "bilibili" {
		t.Fatalf("expected social links to be updated, got %#v", settingsUC.updatedLinks)
	}
}

func TestUpdateSettingsAcceptsFalse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	settingsUC := &settingsUseCaseStub{updatedVisible: true}
	handler := NewSettingsHandler(settingsUC)
	router := gin.New()
	router.PATCH("/settings", handler.Update)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPatch, "/settings", bytes.NewBufferString(`{"project_page_visible":false}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if settingsUC.updatedVisible {
		t.Fatal("expected project page visibility to be updated to false")
	}
}
