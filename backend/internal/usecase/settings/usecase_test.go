package settings

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"siyue-blog/internal/entity"
)

type settingsRepoStub struct {
	links  *[]entity.SocialLink
	footer *entity.FooterConfig
}

func (stub *settingsRepoStub) Get(context.Context) (*entity.SiteSettings, error) {
	return &entity.SiteSettings{ID: 1, SocialLinks: stub.links}, nil
}

func (stub *settingsRepoStub) Update(_ context.Context, _ *bool, links *[]entity.SocialLink, footer *entity.FooterConfig) (*entity.SiteSettings, error) {
	stub.links = links
	stub.footer = footer
	return &entity.SiteSettings{ID: 1, SocialLinks: links, Footer: footer}, nil
}

func TestUpdateNormalizesSocialLinks(t *testing.T) {
	repo := &settingsRepoStub{}
	uc := NewUseCase(repo)
	links := []entity.SocialLink{{Platform: " GitHub ", Label: " 代码 ", Href: " https://github.com/example "}}

	settings, err := uc.Update(context.Background(), nil, &links, nil)

	require.NoError(t, err)
	require.Equal(t, "github", (*settings.SocialLinks)[0].Platform)
	require.Equal(t, "代码", (*settings.SocialLinks)[0].Label)
}

func TestUpdateRejectsUnsafeSocialLink(t *testing.T) {
	repo := &settingsRepoStub{}
	uc := NewUseCase(repo)
	links := []entity.SocialLink{{Platform: "other", Label: "危险链接", Href: "javascript:alert(1)"}}

	_, err := uc.Update(context.Background(), nil, &links, nil)

	require.ErrorContains(t, err, "HTTPS or mailto")
}

func TestUpdateNormalizesAndValidatesFooter(t *testing.T) {
	repo := &settingsRepoStub{}
	uc := NewUseCase(repo)
	footer := entity.FooterConfig{Enabled: true, Brand: "  My Blog  ", Copyright: " © {year} Me "}

	settings, err := uc.Update(context.Background(), nil, nil, &footer)

	require.NoError(t, err)
	require.Equal(t, "My Blog", settings.Footer.Brand)
	require.Equal(t, "© {year} Me", settings.Footer.Copyright)

	footer.Description = strings.Repeat("文", 301)
	_, err = uc.Update(context.Background(), nil, nil, &footer)
	require.ErrorContains(t, err, "footer.description")
}

func TestUpdateValidatesFilingSettings(t *testing.T) {
	repo := &settingsRepoStub{}
	uc := NewUseCase(repo)
	footer := entity.FooterConfig{
		Enabled:              true,
		FilingIconURL:        " /images/beian.png ",
		ICPNumber:            " 粤ICP备12345678号 ",
		ICPURL:               "https://beian.miit.gov.cn/",
		PublicSecurityNumber: "粤公网安备 44000000000000号",
	}

	settings, err := uc.Update(context.Background(), nil, nil, &footer)

	require.NoError(t, err)
	require.Equal(t, "/images/beian.png", settings.Footer.FilingIconURL)
	require.Equal(t, "粤ICP备12345678号", settings.Footer.ICPNumber)

	footer.FilingIconURL = "javascript:alert(1)"
	_, err = uc.Update(context.Background(), nil, nil, &footer)
	require.ErrorContains(t, err, "footer.filing_icon_url")

	footer.FilingIconURL = "/images/beian.png"
	footer.PublicSecurityURL = "/unsafe-relative-link"
	_, err = uc.Update(context.Background(), nil, nil, &footer)
	require.ErrorContains(t, err, "footer.public_security_url")
}
