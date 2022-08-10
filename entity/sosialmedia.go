package entity

import "time"

type SocialMedia struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	SocialMediaURL  string    `json:"social_media_url"`
	ProfileImageURL *string   `json:"profile_image_url"`
	UserID          int64     `json:"user_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type SocialMediaPost struct {
	Name            string `json:"name" validate:"required"`
	SocialMediaURL  string `json:"social_media_url" validate:"required"`
	ProfileImageURL string `json:"profile_image_url"`
}

type SocialMediaPostOutput struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	SocialMediaURL  string    `json:"social_media_url"`
	ProfileImageURL *string   `json:"profile_image_url"`
	UserID          int64     `json:"user_id"`
	CreatedAt       time.Time `json:"created_at"`
}

func (s *SocialMedia) ToSocialMediaPostOutput() *SocialMediaPostOutput {
	out := &SocialMediaPostOutput{
		ID:              s.ID,
		Name:            s.Name,
		SocialMediaURL:  s.SocialMediaURL,
		ProfileImageURL: s.ProfileImageURL,
		UserID:          s.UserID,
		CreatedAt:       s.CreatedAt,
	}
	return out
}

type SocialMediaUpdateOutput struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	SocialMediaURL  string    `json:"social_media_url"`
	ProfileImageURL *string   `json:"profile_image_url"`
	UserID          int64     `json:"user_id"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (s *SocialMedia) ToSocialMediaUpdateOutput() *SocialMediaUpdateOutput {
	out := &SocialMediaUpdateOutput{
		ID:              s.ID,
		Name:            s.Name,
		SocialMediaURL:  s.SocialMediaURL,
		ProfileImageURL: s.ProfileImageURL,
		UserID:          s.UserID,
		UpdatedAt:       s.UpdatedAt,
	}
	return out
}

type SocialMediaGetOutput struct {
	SocialMedia
	User UserGetSocialMedia `json:"user"`
}
