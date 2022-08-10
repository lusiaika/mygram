package database

import (
	"context"
	"database/sql"
	"mygram/entity"
	"strings"
	"time"
)

func (s *Database) PostComment(ctx context.Context, userid int64, i entity.CommentPost) (*entity.Comment, error) {
	result := &entity.Comment{}
	qry := "insert into comments (message, photoid, userid, createdat, updatedat) values (@message, @photoid, @userid, @createdat, @updatedat); select id, message, photoid, userid, createdat from comments where id = SCOPE_IDENTITY()"
	now := time.Now()
	rows, err := s.SqlDb.QueryContext(ctx, qry,
		sql.Named("message", i.Message),
		sql.Named("photoid", i.PhotoID),
		sql.Named("userid", userid),
		sql.Named("createdat", now),
		sql.Named("updatedat", now))
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.Message,
			&result.PhotoID,
			&result.UserID,
			&result.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (s *Database) GetComments(ctx context.Context) ([]entity.CommentGetOutput, error) {
	var result []entity.CommentGetOutput
	var qry strings.Builder
	qry.WriteString("select c.id, c.message, c.photoid, c.userid, c.createdat, c.updatedat,")
	qry.WriteString(" p.title, p.caption, p.photourl,")
	qry.WriteString(" u.email, u.username from comments c")
	qry.WriteString(" join photos p on c.photoid=p.id")
	qry.WriteString(" join users u on c.userid=u.id")

	rows, err := s.SqlDb.QueryContext(ctx, qry.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var row entity.CommentGetOutput
		err := rows.Scan(
			&row.ID,
			&row.Message,
			&row.PhotoID,
			&row.UserID,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.Photo.Title,
			&row.Photo.Caption,
			&row.Photo.PhotoUrl,
			&row.User.Email,
			&row.User.Username,
		)
		if err != nil {
			return nil, err
		}
		row.User.ID = row.UserID
		row.Photo.UserID = row.UserID
		row.Photo.ID = row.PhotoID
		result = append(result, row)
	}
	return result, nil
}

func (s *Database) GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error) {
	result := &entity.Comment{}

	rows, err := s.SqlDb.QueryContext(ctx, "select id, message, photoid, userid, createdat, updatedat from comments where id = @ID",
		sql.Named("ID", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.Message,
			&result.PhotoID,
			&result.UserID,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (s *Database) UpdateComment(ctx context.Context, userid int64, id int64, message string) (*entity.Comment, error) {
	result := &entity.Comment{}
	now := time.Now()
	qry := "update comments set message=@message, updatedat=@updatedat where id = @ID and userid = userid; select id, userid, photoid, message, updatedat from comments where id = @ID"
	rows, err := s.SqlDb.QueryContext(ctx, qry,
		sql.Named("message", message),
		sql.Named("updatedat", now),
		sql.Named("userid", userid),
		sql.Named("ID", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.UserID,
			&result.PhotoID,
			&result.Message,
			&result.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (s *Database) DeleteComment(ctx context.Context, userid int64, id int64) (string, error) {
	var result string
	qry := "delete from comments where id=@id and userid = @userid"
	_, err := s.SqlDb.ExecContext(ctx, qry,
		sql.Named("userid", userid),
		sql.Named("id", id))
	if err != nil {
		return "", err
	}

	result = "Your photo has been successfully deleted"

	return result, nil
}
