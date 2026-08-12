package entity

import "time"

type SiteSettings struct {
	ID                 uint          `json:"-" gorm:"primaryKey"`
	ProjectPageVisible bool          `json:"project_page_visible" gorm:"not null"`
	SocialLinks        *[]SocialLink `json:"social_links" gorm:"serializer:json;type:jsonb"`
	Footer             *FooterConfig `json:"footer" gorm:"serializer:json;type:jsonb"`
	UpdatedAt          time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
}

type SocialLink struct {
	Platform string `json:"platform"`
	Label    string `json:"label"`
	Href     string `json:"href"`
}

type FooterConfig struct {
	Enabled              bool   `json:"enabled"`
	Brand                string `json:"brand"`
	Description          string `json:"description"`
	Copyright            string `json:"copyright"`
	ShowSocialLinks      bool   `json:"show_social_links"`
	FilingIconURL        string `json:"filing_icon_url"`
	ICPNumber            string `json:"icp_number"`
	ICPURL               string `json:"icp_url"`
	PublicSecurityNumber string `json:"public_security_number"`
	PublicSecurityURL    string `json:"public_security_url"`
	FilingInformation    string `json:"filing_information"`
}

func (SiteSettings) TableName() string {
	return "site_settings"
}
