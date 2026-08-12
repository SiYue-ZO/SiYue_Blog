package homepage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"siyue-blog/internal/entity"
)

type homePageRepoStub struct {
	content *entity.HomePageContent
	saved   entity.HomePageConfig
}

func (stub *homePageRepoStub) Get(context.Context) (*entity.HomePageContent, error) {
	return stub.content, nil
}
func (stub *homePageRepoStub) SaveDraft(_ context.Context, config entity.HomePageConfig, _ uint64) (*entity.HomePageContent, error) {
	stub.saved = config
	stub.content.DraftConfig = config
	return stub.content, nil
}
func (stub *homePageRepoStub) Publish(context.Context, uint64) (*entity.HomePageContent, error) {
	stub.content.PublishedConfig = stub.content.DraftConfig
	return stub.content, nil
}
func (stub *homePageRepoStub) ResetDraft(context.Context, uint64) (*entity.HomePageContent, error) {
	stub.content.DraftConfig = stub.content.PublishedConfig
	return stub.content, nil
}

func validConfig() entity.HomePageConfig {
	return entity.HomePageConfig{
		SchemaVersion: 1,
		Order:         []string{"hero", "intro", "tech_stack", "featured_projects", "latest_posts", "contact_cta"},
		Sections: entity.HomeSections{
			Hero:             entity.HomeHeroSection{Enabled: true, TitlePrefix: "你好，我是", PrimaryButton: entity.HomeButton{Label: "博客", Href: "/blog"}, SecondaryButton: entity.HomeButton{Label: "项目", Href: "/projects"}},
			Intro:            entity.HomeIntroSection{Enabled: true, Title: "关于", Paragraphs: []string{"介绍"}},
			TechStack:        entity.HomeTechStackSection{Enabled: true, Title: "技术栈", Items: []string{"Go"}},
			FeaturedProjects: entity.HomeFeaturedProjectsSection{Enabled: true, Title: "项目", Limit: 3},
			LatestPosts:      entity.HomeLatestPostsSection{Enabled: true, Title: "文章", Limit: 3, Source: "latest"},
			ContactCTA:       entity.HomeContactCTASection{Enabled: true, Title: "联系", Description: "欢迎联系", Button: entity.HomeButton{Label: "联系我", Href: "/contact"}},
		},
	}
}

func TestSaveDraftRejectsIncompleteOrder(t *testing.T) {
	config := validConfig()
	config.Order = config.Order[:5]
	stub := &homePageRepoStub{content: &entity.HomePageContent{DraftConfig: validConfig()}}
	uc := NewUseCase(stub)

	_, err := uc.SaveDraft(context.Background(), config, 1)

	require.ErrorContains(t, err, "order must contain")
}

func TestSaveDraftRejectsUnsafeLink(t *testing.T) {
	config := validConfig()
	config.Sections.ContactCTA.Button.Href = "javascript:alert(1)"
	stub := &homePageRepoStub{content: &entity.HomePageContent{DraftConfig: validConfig()}}
	uc := NewUseCase(stub)

	_, err := uc.SaveDraft(context.Background(), config, 1)

	require.ErrorContains(t, err, "internal path")
}

func TestPublishCopiesValidatedDraft(t *testing.T) {
	config := validConfig()
	config.Sections.Hero.Highlight = "SiYue"
	stub := &homePageRepoStub{content: &entity.HomePageContent{DraftConfig: config}}
	uc := NewUseCase(stub)

	content, err := uc.Publish(context.Background(), 1)

	require.NoError(t, err)
	require.Equal(t, "SiYue", content.PublishedConfig.Sections.Hero.Highlight)
}
