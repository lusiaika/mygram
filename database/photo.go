package database

import (
	"context"
	"database/sql"
	"mygram/entity"
	"strings"
	"time"
)

func (s *Database) PostPhoto(ctx context.Context, u int64, i entity.PhotoPost) (*entity.Photo, error) {
	result := &entity.Photo{}
	qry := "insert into photos (title, caption, photourl, userid, createdat, updatedat) values (@title, @caption, @photourl, @userid, @createdat, @updatedat); select id, title, caption, photourl, userid, createdat from photos where id = SCOPE_IDENTITY()"
	now := time.Now()

	rows, err := s.SqlDb.QueryContext(ctx, qry,
		sql.Named("title", i.Title),
		sql.Named("caption", i.Caption),
		sql.Named("photourl", i.PhotoUrl),
		sql.Named("userid", u),
		sql.Named("createdat", now),
		sql.Named("updatedat", now))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.Title,
			&result.Caption,
			&result.PhotoUrl,
			&result.UserID,
			&result.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (s *Database) GetPhotos(ctx context.Context) ([]entity.PhotoGetOutput, error) {
	var result []entity.PhotoGetOutput
	var qry strings.Builder
	qry.WriteString("select p.id, p.title, p.caption, p.photourl, p.userid, p.createdat, p.updatedat, u.email, u.username from photos p")
	qry.WriteString(" join users u on p.userid=u.id")
	rows, err := s.SqlDb.QueryContext(ctx, qry.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var row entity.PhotoGetOutput
		err := rows.Scan(
			&row.ID,
			&row.Title,
			&row.Caption,
			&row.PhotoUrl,
			&row.UserID,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.User.Email,
			&row.User.Username,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	return result, nil
}

func (s *Database) GetPhotoByID(ctx context.Context, id int64) (*entity.Photo, error) {
	result := &entity.Photo{}
	qry := "select id, title, caption, photourl, userid, createdat, updatedat from photos where id = @ID"
	rows, err := s.SqlDb.QueryContext(ctx, qry,
		sql.Named("ID", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.Title,
			&result.Caption,
			&result.PhotoUrl,
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

func (s *Database) UpdatePhoto(ctx context.Context, userid int64, id int64, i entity.PhotoPost) (*entity.Photo, error) {
	result := &entity.Photo{}
	now := time.Now()
	qry := "update photos set title=@title, caption=@caption, photourl=@photourl, updatedat=@updatedat where id = @ID and userid = @userid; select id, title, caption, photourl, userid, updatedat from photos where id = @ID"
	rows, err := s.SqlDb.QueryContext(ctx, qry,
		sql.Named("title", i.Title),
		sql.Named("caption", i.Caption),
		sql.Named("photourl", i.PhotoUrl),
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
			&result.Title,
			&result.Caption,
			&result.PhotoUrl,
			&result.UserID,
			&result.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (s *Database) DeletePhoto(ctx context.Context, userid int64, id int64) (string, error) {
	var result string
	qry := "delete from comments where photoid=@id and userid=@userid; delete from photos where id=@id and userid=@userid"
	_, err := s.SqlDb.ExecContext(ctx, qry,
		sql.Named("userid", userid),
		sql.Named("id", id))
	if err != nil {
		return "", err
	}

	result = "Your photo has been successfully deleted"

	return result, nil
}
