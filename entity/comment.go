package entity

import "time"

type Comment struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	PhotoID   int64     `json:"photo_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CommentPost struct {
	PhotoID int    `json:"photo_id" validate:"required"`
	Message string `json:"message" validate:"required"`
}

type CommentUpdate struct {
	Message string `json:"message" validate:"required"`
}

type CommentPostOutput struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	PhotoID   int64     `json:"photo_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *Comment) ToCommentPostOutput() *CommentPostOutput {
	out := &CommentPostOutput{
		ID:        c.ID,
		UserID:    c.UserID,
		PhotoID:   c.PhotoID,
		Message:   c.Message,
		CreatedAt: c.CreatedAt,
	}
	return out
}

type CommentUpdateOutput struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	PhotoID   int64     `json:"photo_id"`
	Message   string    `json:"message"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Comment) ToCommentUpdateOutput() *CommentUpdateOutput {
	out := &CommentUpdateOutput{
		ID:        c.ID,
		UserID:    c.UserID,
		PhotoID:   c.PhotoID,
		Message:   c.Message,
		UpdatedAt: c.UpdatedAt,
	}
	return out
}

type CommentGetOutput struct {
	Comment
	User  UserGetComment  `json:"user"`
	Photo PhotoGetComment `json:"photo"`
}
