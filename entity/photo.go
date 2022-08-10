package entity

import "time"

type Photo struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoUrl  string    `json:"photo_url"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PhotoGetComment struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
	UserID   int64  `json:"user_id"`
}

func (p *Photo) ToPhotoGetComment() *PhotoGetComment {
	out := &PhotoGetComment{
		ID:       p.ID,
		Title:    p.Title,
		Caption:  p.Caption,
		PhotoUrl: p.PhotoUrl,
		UserID:   p.UserID,
	}
	return out
}

type PhotoPost struct {
	Title    string `json:"title" validate:"required"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url" validate:"required"`
}

type PhotoPostOutput struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoUrl  string    `json:"photo_url"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (p *Photo) ToPhotoPostOutput() *PhotoPostOutput {
	out := &PhotoPostOutput{
		ID:        p.ID,
		Title:     p.Title,
		Caption:   p.Caption,
		PhotoUrl:  p.PhotoUrl,
		UserID:    p.UserID,
		CreatedAt: p.CreatedAt,
	}
	return out
}

type PhotoUpdateOutput struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoUrl  string    `json:"photo_url"`
	UserID    int64     `json:"user_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Photo) ToPhotoUpdateOutput() *PhotoUpdateOutput {
	out := &PhotoUpdateOutput{
		ID:        p.ID,
		Title:     p.Title,
		Caption:   p.Caption,
		PhotoUrl:  p.PhotoUrl,
		UserID:    p.UserID,
		UpdatedAt: p.UpdatedAt,
	}
	return out
}

type PhotoGetOutput struct {
	Photo
	User UserUpdate `json:"user"`
}
