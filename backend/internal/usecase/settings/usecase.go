package settings

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"

	"siyue-blog/internal/entity"
	"siyue-blog/internal/repo"
	"siyue-blog/internal/usecase"
)

type useCase struct {
	settingsRepo repo.SettingsRepo
}

func NewUseCase(settingsRepo repo.SettingsRepo) usecase.SettingsUseCase {
	return &useCase{settingsRepo: settingsRepo}
}

func (uc *useCase) Get(ctx context.Context) (*entity.SiteSettings, error) {
	settings, err := uc.settingsRepo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("get site settings: %w", err)
	}
	return settings, nil
}

var allowedPlatforms = map[string]bool{
	"github": true, "email": true, "bilibili": true, "twitter": true,
	"linkedin": true, "youtube": true, "weibo": true, "website": true, "other": true,
}

func validateSocialLinks(links []entity.SocialLink) error {
	if len(links) > 20 {
		return fmt.Errorf("social_links may contain at most 20 items")
	}
	for index, link := range links {
		link.Platform = strings.ToLower(strings.TrimSpace(link.Platform))
		link.Label = strings.TrimSpace(link.Label)
		link.Href = strings.TrimSpace(link.Href)
		if !allowedPlatforms[link.Platform] {
			return fmt.Errorf("social_links[%d].platform is unsupported", index)
		}
		if length := utf8.RuneCountInString(link.Label); length == 0 || length > 40 {
			return fmt.Errorf("social_links[%d].label must contain 1 to 40 characters", index)
		}
		if utf8.RuneCountInString(link.Href) > 512 {
			return fmt.Errorf("social_links[%d].href must not exceed 512 characters", index)
		}
		parsed, err := url.ParseRequestURI(link.Href)
		if err != nil || (parsed.Scheme != "https" && parsed.Scheme != "mailto") {
			return fmt.Errorf("social_links[%d].href must be an HTTPS or mailto URL", index)
		}
		links[index] = link
	}
	return nil
}

func validateFooter(footer *entity.FooterConfig) error {
	footer.Brand = strings.TrimSpace(footer.Brand)
	footer.Description = strings.TrimSpace(footer.Description)
	footer.Copyright = strings.TrimSpace(footer.Copyright)
	footer.FilingIconURL = strings.TrimSpace(footer.FilingIconURL)
	footer.ICPNumber = strings.TrimSpace(footer.ICPNumber)
	footer.ICPURL = strings.TrimSpace(footer.ICPURL)
	footer.PublicSecurityNumber = strings.TrimSpace(footer.PublicSecurityNumber)
	footer.PublicSecurityURL = strings.TrimSpace(footer.PublicSecurityURL)
	footer.FilingInformation = strings.TrimSpace(footer.FilingInformation)
	if utf8.RuneCountInString(footer.Brand) > 80 {
		return fmt.Errorf("footer.brand must not exceed 80 characters")
	}
	if utf8.RuneCountInString(footer.Description) > 300 {
		return fmt.Errorf("footer.description must not exceed 300 characters")
	}
	if utf8.RuneCountInString(footer.Copyright) > 200 {
		return fmt.Errorf("footer.copyright must not exceed 200 characters")
	}
	if utf8.RuneCountInString(footer.ICPNumber) > 100 {
		return fmt.Errorf("footer.icp_number must not exceed 100 characters")
	}
	if utf8.RuneCountInString(footer.PublicSecurityNumber) > 100 {
		return fmt.Errorf("footer.public_security_number must not exceed 100 characters")
	}
	if utf8.RuneCountInString(footer.FilingInformation) > 300 {
		return fmt.Errorf("footer.filing_information must not exceed 300 characters")
	}
	for field, value := range map[string]string{
		"filing_icon_url":     footer.FilingIconURL,
		"icp_url":             footer.ICPURL,
		"public_security_url": footer.PublicSecurityURL,
	} {
		if value == "" {
			continue
		}
		if len(value) > 512 {
			return fmt.Errorf("footer.%s must not exceed 512 characters", field)
		}
		parsed, err := url.ParseRequestURI(value)
		isAllowedSitePath := field == "filing_icon_url" && strings.HasPrefix(value, "/") && !strings.HasPrefix(value, "//")
		if err != nil || (!isAllowedSitePath && parsed.Scheme != "https") {
			requirement := "an HTTPS URL"
			if field == "filing_icon_url" {
				requirement = "an HTTPS URL or site path"
			}
			return fmt.Errorf("footer.%s must be %s", field, requirement)
		}
	}
	return nil
}

func (uc *useCase) Update(ctx context.Context, projectPageVisible *bool, socialLinks *[]entity.SocialLink, footer *entity.FooterConfig) (*entity.SiteSettings, error) {
	if projectPageVisible == nil && socialLinks == nil && footer == nil {
		return nil, fmt.Errorf("at least one setting is required")
	}
	if socialLinks != nil {
		if err := validateSocialLinks(*socialLinks); err != nil {
			return nil, err
		}
	}
	if footer != nil {
		if err := validateFooter(footer); err != nil {
			return nil, err
		}
	}
	settings, err := uc.settingsRepo.Update(ctx, projectPageVisible, socialLinks, footer)
	if err != nil {
		return nil, fmt.Errorf("update site settings: %w", err)
	}
	return settings, nil
}
