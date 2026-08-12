package homepage

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

var sectionKeys = []string{
	entity.HomeSectionHero,
	entity.HomeSectionIntro,
	entity.HomeSectionTechStack,
	entity.HomeSectionFeaturedProjects,
	entity.HomeSectionLatestPosts,
	entity.HomeSectionContactCTA,
}

type useCase struct {
	repo repo.HomePageRepo
}

func NewUseCase(homePageRepo repo.HomePageRepo) usecase.HomePageUseCase {
	return &useCase{repo: homePageRepo}
}

func validateText(label, value string, max int, required bool) error {
	length := utf8.RuneCountInString(strings.TrimSpace(value))
	if required && length == 0 {
		return fmt.Errorf("%s is required", label)
	}
	if length > max {
		return fmt.Errorf("%s must not exceed %d characters", label, max)
	}
	return nil
}

func ensureDistinct(label string, values []string) error {
	seen := make(map[string]bool, len(values))
	for _, value := range values {
		key := strings.ToLower(strings.TrimSpace(value))
		if seen[key] {
			return fmt.Errorf("%s contains duplicate value %q", label, value)
		}
		seen[key] = true
	}
	return nil
}

func validateHref(label, value string) error {
	if err := validateText(label, value, 512, true); err != nil {
		return err
	}
	if strings.HasPrefix(value, "/") && !strings.HasPrefix(value, "//") {
		return nil
	}
	parsed, err := url.ParseRequestURI(value)
	if err != nil || (parsed.Scheme != "https" && parsed.Scheme != "mailto") {
		return fmt.Errorf("%s must be an internal path, HTTPS URL, or mailto URL", label)
	}
	return nil
}

func validateButton(label string, button entity.HomeButton) error {
	if err := validateText(label+" label", button.Label, 40, true); err != nil {
		return err
	}
	return validateHref(label+" link", button.Href)
}

func validateConfig(config entity.HomePageConfig) error {
	if config.SchemaVersion != 1 {
		return fmt.Errorf("unsupported schema_version")
	}
	if len(config.Order) != len(sectionKeys) {
		return fmt.Errorf("order must contain all home page sections")
	}
	seen := make(map[string]bool, len(config.Order))
	for _, key := range config.Order {
		if seen[key] {
			return fmt.Errorf("order contains duplicate section %q", key)
		}
		seen[key] = true
	}
	for _, key := range sectionKeys {
		if !seen[key] {
			return fmt.Errorf("order is missing section %q", key)
		}
	}

	h := config.Sections.Hero
	for label, value := range map[string]string{"hero title prefix": h.TitlePrefix, "hero highlight": h.Highlight, "hero subtitle": h.Subtitle} {
		if err := validateText(label, value, 200, false); err != nil {
			return err
		}
	}
	if err := validateButton("primary button", h.PrimaryButton); err != nil {
		return err
	}
	if err := validateButton("secondary button", h.SecondaryButton); err != nil {
		return err
	}

	i := config.Sections.Intro
	if err := validateText("intro title", i.Title, 100, true); err != nil {
		return err
	}
	if len(i.Paragraphs) > 6 {
		return fmt.Errorf("intro may contain at most 6 paragraphs")
	}
	for _, paragraph := range i.Paragraphs {
		if err := validateText("intro paragraph", paragraph, 1000, true); err != nil {
			return err
		}
	}
	if i.ImageURL != "" {
		if strings.HasPrefix(i.ImageURL, "/") && !strings.HasPrefix(i.ImageURL, "//") {
			// Local public image.
		} else if parsed, err := url.ParseRequestURI(i.ImageURL); err != nil || parsed.Scheme != "https" || parsed.Host == "" {
			return fmt.Errorf("intro image must be a local path or HTTPS URL")
		}
	}

	t := config.Sections.TechStack
	if err := validateText("tech stack title", t.Title, 100, true); err != nil {
		return err
	}
	if err := validateText("tech stack subtitle", t.Subtitle, 200, false); err != nil {
		return err
	}
	if len(t.Items) > 20 {
		return fmt.Errorf("tech stack may contain at most 20 items")
	}
	for _, item := range t.Items {
		if err := validateText("tech stack item", item, 50, true); err != nil {
			return err
		}
	}
	if err := ensureDistinct("tech stack", t.Items); err != nil {
		return err
	}

	p := config.Sections.FeaturedProjects
	if err := validateText("featured projects title", p.Title, 100, true); err != nil {
		return err
	}
	if err := validateText("featured projects subtitle", p.Subtitle, 200, false); err != nil {
		return err
	}
	if p.Limit < 1 || p.Limit > 12 {
		return fmt.Errorf("featured projects limit must be between 1 and 12")
	}

	l := config.Sections.LatestPosts
	if err := validateText("latest posts title", l.Title, 100, true); err != nil {
		return err
	}
	if err := validateText("latest posts subtitle", l.Subtitle, 200, false); err != nil {
		return err
	}
	if l.Limit < 1 || l.Limit > 12 {
		return fmt.Errorf("latest posts limit must be between 1 and 12")
	}
	if l.Source != "latest" && l.Source != "pinned" {
		return fmt.Errorf("latest posts source must be latest or pinned")
	}

	c := config.Sections.ContactCTA
	if err := validateText("contact title", c.Title, 100, true); err != nil {
		return err
	}
	if err := validateText("contact description", c.Description, 500, true); err != nil {
		return err
	}
	return validateButton("contact button", c.Button)
}

func (uc *useCase) GetPublished(ctx context.Context) (*entity.HomePageConfig, error) {
	content, err := uc.repo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("get published home page: %w", err)
	}
	return &content.PublishedConfig, nil
}

func (uc *useCase) GetAdmin(ctx context.Context) (*entity.HomePageContent, error) {
	content, err := uc.repo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("get home page: %w", err)
	}
	return content, nil
}

func (uc *useCase) SaveDraft(ctx context.Context, config entity.HomePageConfig, userID uint64) (*entity.HomePageContent, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}
	content, err := uc.repo.SaveDraft(ctx, config, userID)
	if err != nil {
		return nil, fmt.Errorf("save home page draft: %w", err)
	}
	return content, nil
}

func (uc *useCase) Publish(ctx context.Context, userID uint64) (*entity.HomePageContent, error) {
	content, err := uc.repo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("get home page draft: %w", err)
	}
	if err := validateConfig(content.DraftConfig); err != nil {
		return nil, err
	}
	content, err = uc.repo.Publish(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("publish home page: %w", err)
	}
	return content, nil
}

func (uc *useCase) ResetDraft(ctx context.Context, userID uint64) (*entity.HomePageContent, error) {
	content, err := uc.repo.ResetDraft(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("reset home page draft: %w", err)
	}
	return content, nil
}
