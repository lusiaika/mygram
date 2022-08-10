package database

import (
	"context"
	"database/sql"
	"mygram/entity"
	"strings"
	"time"
)

func (s *Database) PostSocialMedia(ctx context.Context, userid int64, i entity.SocialMediaPost) (*entity.SocialMedia, error) {
	result := &entity.SocialMedia{}
	qry := "insert into socialmedias (name, socialmediaurl, profileimageurl, userid, createdat, updatedat) values (@name, @socialmediaurl, @profileimageurl, @userid, @createdat, @updatedat); select id, name, socialmediaurl, profileimageurl, userid, createdat,updatedat from socialmedias"
	now := time.Now()
	rows, err := s.SqlDb.QueryContext(ctx, qry,
		sql.Named("name", i.Name),
		sql.Named("socialmediaurl", i.SocialMediaURL),
		sql.Named("profileimageurl", i.ProfileImageURL),
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
			&result.Name,
			&result.SocialMediaURL,
			&result.ProfileImageURL,
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

func (s *Database) GetSocialMedias(ctx context.Context) ([]entity.SocialMediaGetOutput, error) {
	var result []entity.SocialMediaGetOutput
	var qry strings.Builder
	qry.WriteString("select s.id, s.name, s.socialmediaurl, s.userid, s.createdat, s.updatedat,")
	qry.WriteString(" u.username, s.profileimageurl")
	qry.WriteString(" from socialmedias s join users u on s.userid=u.id")
	rows, err := s.SqlDb.QueryContext(ctx, qry.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var row entity.SocialMediaGetOutput
		err := rows.Scan(
			&row.ID,
			&row.Name,
			&row.SocialMediaURL,
			&row.UserID,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.User.Username,
			&row.User.ProfileImageURL,
		)
		if err != nil {
			return nil, err
		}
		row.User.ID = row.UserID
		result = append(result, row)
	}
	return result, nil
}

func (s *Database) GetSocialMediaByID(ctx context.Context, id int64) (*entity.SocialMedia, error) {
	result := &entity.SocialMedia{}

	rows, err := s.SqlDb.QueryContext(ctx, "select id, name, socialmediaurl, userid, createdat, updatedat from socialmedias where id = @ID",
		sql.Named("ID", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.Name,
			&result.SocialMediaURL,
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

func (s *Database) UpdateSocialMedia(ctx context.Context, userid int64, id int64, i entity.SocialMediaPost) (*entity.SocialMedia, error) {
	result := &entity.SocialMedia{}
	now := time.Now()
	qry := "update socialmedias set name=@name, socialmediaurl=@socialmediaurl, profileimageurl=@profileimageurl, updatedat=@updatedat where id = @ID; select id, name, socialmediaurl, profileimageurl, userid, updatedat from socialmedias where id = @ID"
	rows, err := s.SqlDb.QueryContext(ctx, qry,
		sql.Named("name", i.Name),
		sql.Named("socialmediaurl", i.SocialMediaURL),
		sql.Named("profileimageurl", i.ProfileImageURL),
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
			&result.Name,
			&result.SocialMediaURL,
			&result.ProfileImageURL,
			&result.UserID,
			&result.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (s *Database) DeleteSocialMedia(ctx context.Context, userid int64, id int64) (string, error) {
	var result string
	qry := "delete from socialmedias where id=@id and userid=@userid"
	_, err := s.SqlDb.ExecContext(ctx, qry,
		sql.Named("userid", userid),
		sql.Named("id", id))
	if err != nil {
		return "", err
	}

	result = "Your social media has been successfully deleted"

	return result, nil
}
