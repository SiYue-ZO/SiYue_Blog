package entity

import "time"

const (
	HomeSectionHero             = "hero"
	HomeSectionIntro            = "intro"
	HomeSectionTechStack        = "tech_stack"
	HomeSectionFeaturedProjects = "featured_projects"
	HomeSectionLatestPosts      = "latest_posts"
	HomeSectionContactCTA       = "contact_cta"
)

type HomeButton struct {
	Label string `json:"label"`
	Href  string `json:"href"`
}

type HomeHeroSection struct {
	Enabled         bool       `json:"enabled"`
	TitlePrefix     string     `json:"title_prefix"`
	Highlight       string     `json:"highlight"`
	Subtitle        string     `json:"subtitle"`
	PrimaryButton   HomeButton `json:"primary_button"`
	SecondaryButton HomeButton `json:"secondary_button"`
	ShowSocialLinks bool       `json:"show_social_links"`
}

type HomeIntroSection struct {
	Enabled    bool     `json:"enabled"`
	Title      string   `json:"title"`
	ImageURL   string   `json:"image_url"`
	Paragraphs []string `json:"paragraphs"`
}

type HomeTechStackSection struct {
	Enabled  bool     `json:"enabled"`
	Title    string   `json:"title"`
	Subtitle string   `json:"subtitle"`
	Items    []string `json:"items"`
}

type HomeFeaturedProjectsSection struct {
	Enabled  bool   `json:"enabled"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Limit    int    `json:"limit"`
}

type HomeLatestPostsSection struct {
	Enabled  bool   `json:"enabled"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Limit    int    `json:"limit"`
	Source   string `json:"source"`
}

type HomeContactCTASection struct {
	Enabled     bool       `json:"enabled"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Button      HomeButton `json:"button"`
}

type HomeSections struct {
	Hero             HomeHeroSection             `json:"hero"`
	Intro            HomeIntroSection            `json:"intro"`
	TechStack        HomeTechStackSection        `json:"tech_stack"`
	FeaturedProjects HomeFeaturedProjectsSection `json:"featured_projects"`
	LatestPosts      HomeLatestPostsSection      `json:"latest_posts"`
	ContactCTA       HomeContactCTASection       `json:"contact_cta"`
}

type HomePageConfig struct {
	SchemaVersion int          `json:"schema_version"`
	Order         []string     `json:"order"`
	Sections      HomeSections `json:"sections"`
}

type HomePageContent struct {
	DraftConfig     HomePageConfig `json:"draft_config"`
	PublishedConfig HomePageConfig `json:"published_config"`
	ConfigVersion   int            `json:"config_version"`
	UpdatedBy       *uint64        `json:"updated_by,omitempty"`
	PublishedAt     time.Time      `json:"published_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}
