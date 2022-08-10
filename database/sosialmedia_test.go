package database

import (
	"context"
	"errors"
	"mygram/entity"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDatabase_PostSocialMedia(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	db, mock := NewMock()
	defer db.Close()
	dbtes := Database{
		SqlDb: db,
	}
	qry := "insert into socialmedias (name, socialmediaurl, profileimageurl, userid, createdat, updatedat) values (@name, @socialmediaurl, @profileimageurl, @userid, @createdat, @updatedat); select id, name, socialmediaurl, profileimageurl, userid, createdat,updatedat from socialmedias"

	inp := entity.SocialMediaPost{
		Name:            "socialmedia orang ganteng",
		SocialMediaURL:  "https://socialmediaurl.com/socialmediaurl.jpg",
		ProfileImageURL: "https://profileimageurl.com/profileimageurl.jpg",
	}

	t.Run("postsocialmedia database down", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs(inp.Name, inp.SocialMediaURL, inp.ProfileImageURL, int64(1), time.Now(), time.Now()).
			WillReturnError(errors.New("db down"))
		out, err := dbtes.PostSocialMedia(ctx, int64(1), inp)
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "db down", err.Error())
	})

	t.Run("postsocialmedia required userid", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs(inp.Name, inp.SocialMediaURL, inp.ProfileImageURL, int64(0), time.Now(), time.Now()).
			WillReturnError(errors.New("required userid"))
		out, err := dbtes.PostSocialMedia(ctx, int64(0), inp)
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "required userid", err.Error())
	})

	t.Run("postsocialmedia success", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "name", "socialmediaurl", "profileimageurl", "userid", "createdat", "updatedat"}).
			AddRow(1, "SocialMedia Name", "http://socialmediaurl.com/socialmediaurl.jpg", "http://profileimageurl.com/profileimageurl.jpg", 1, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs(inp.Name, inp.SocialMediaURL, inp.ProfileImageURL, int64(1), time.Now(), time.Now()).
			WillReturnRows(rows)
		out, err := dbtes.PostSocialMedia(ctx, int64(1), inp)
		assert.NotNil(t, out)
		assert.NoError(t, err)
	})
}

func TestDatabase_GetSocialMedias(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	db, mock := NewMock()
	defer db.Close()
	dbtes := Database{
		SqlDb: db,
	}
	var qry strings.Builder
	qry.WriteString("select s.id, s.name, s.socialmediaurl, s.userid, s.createdat, s.updatedat,")
	qry.WriteString(" u.username, s.profileimageurl")
	qry.WriteString(" from socialmedias s join users u on s.userid=u.id")
	t.Run("getsocialmedias database down", func(t *testing.T) {
		mock.ExpectQuery(qry.String()).
			WillReturnError(errors.New("db down"))
		out, err := dbtes.GetSocialMedias(ctx)
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "db down", err.Error())
	})

	t.Run("getsocialmedias success", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "name", "socialmediaurl", "userid", "createdat", "updatedat", "username", "profileimageurl"}).
			AddRow(1, "SocialMedia Name", "http://socialmediaurl.com/socialmediaurl.jpg", 1, time.Now(), time.Now(), "User Name", "http://profileimageurl/profile.jpg")

		mock.ExpectQuery(qry.String()).WillReturnRows(rows)
		out, err := dbtes.GetSocialMedias(ctx)
		assert.NotNil(t, out)
		assert.NoError(t, err)
	})
}

func TestDatabase_GetSocialMediaByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	db, mock := NewMock()
	defer db.Close()
	dbtes := Database{
		SqlDb: db,
	}

	qry := "select id, name, socialmediaurl, userid, createdat, updatedat from socialmedias where id = @ID"
	t.Run("getsocialmediabyid database down", func(t *testing.T) {
		mock.ExpectQuery(qry).
			WithArgs(int64(1)).
			WillReturnError(errors.New("db down"))
		out, err := dbtes.GetSocialMediaByID(ctx, int64(1))
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "db down", err.Error())
	})

	t.Run("getsocialmediabyid required userid", func(t *testing.T) {
		mock.ExpectQuery(qry).
			WithArgs(int64(0)).
			WillReturnError(errors.New("required userid"))
		out, err := dbtes.GetSocialMediaByID(ctx, int64(0))
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "required userid", err.Error())
	})

	t.Run("getsocialmediabyid success", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "name", "socialmediaurl", "userid", "createdat", "updatedat"}).
			AddRow(1, "SocialMedia Name", "http://socialmediaurl.com/socialmediaurl.jpg", 1, time.Now(), time.Now())

		mock.ExpectQuery(qry).
			WithArgs(int64(1)).
			WillReturnRows(rows)
		out, err := dbtes.GetSocialMediaByID(ctx, int64(1))
		assert.NotNil(t, out)
		assert.NoError(t, err)
	})
}

func TestDatabase_UpdateSocialMedia(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	db, mock := NewMock()
	defer db.Close()
	dbtes := Database{
		SqlDb: db,
	}

	qry := "update socialmedias set name=@name, socialmediaurl=@socialmediaurl, profileimageurl=@profileimageurl, updatedat=@updatedat where id = @ID; select id, name, socialmediaurl, profileimageurl, userid, updatedat from socialmedias where id = @ID"
	inp := entity.SocialMediaPost{
		Name:            "socialmedia orang ganteng",
		SocialMediaURL:  "https://socialmediaurl.com/socialmediaurl.jpg",
		ProfileImageURL: "https://profileimageurl.com/profileimageurl.jpg",
	}
	t.Run("updatesocialmedia database down", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs(inp.Name, inp.SocialMediaURL, inp.ProfileImageURL, time.Now(), int64(1), int64(1)).
			WillReturnError(errors.New("db down"))
		out, err := dbtes.UpdateSocialMedia(ctx, int64(1), int64(1), inp)
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "db down", err.Error())
	})

	t.Run("updatesocialmedia required userid", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs(inp.Name, inp.SocialMediaURL, inp.ProfileImageURL, time.Now(), int64(0), int64(1)).
			WillReturnError(errors.New("required userid"))
		out, err := dbtes.UpdateSocialMedia(ctx, int64(0), int64(1), inp)
		assert.Nil(t, out)
		assert.Equal(t, "required userid", err.Error())
	})

	t.Run("updatesocialmedia required id", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs(inp.Name, inp.SocialMediaURL, inp.ProfileImageURL, time.Now(), int64(1), int64(0)).
			WillReturnError(errors.New("required id"))
		out, err := dbtes.UpdateSocialMedia(ctx, int64(1), int64(0), inp)
		assert.Nil(t, out)
		assert.Equal(t, "required id", err.Error())
	})

	t.Run("updatesocialmedia success", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "name", "socialmediaurl", "profileimageurl", "userid", "updatedat"}).
			AddRow(1, "SocialMedia Name", "http://socialmediaurl.com/socialmediaurl.jpg", "http://profileimageurl.com/profileimage.jpg", 1, time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs(inp.Name, inp.SocialMediaURL, inp.ProfileImageURL, time.Now(), int64(1), int64(1)).
			WillReturnRows(rows)
		out, err := dbtes.UpdateSocialMedia(ctx, int64(1), int64(1), inp)
		assert.NotNil(t, out)
		assert.NoError(t, err)
	})
}

func TestDatabase_DeleteSocialMedia(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	db, mock := NewMock()
	defer db.Close()
	dbtes := Database{
		SqlDb: db,
	}
	qry := "delete from socialmedias where id=@id and userid=@userid"
	t.Run("deletesocialmedia database down", func(t *testing.T) {
		mock.ExpectExec(qry).
			WithArgs(int64(1), int64(1)).
			WillReturnError(errors.New("db down"))
		out, err := dbtes.DeleteSocialMedia(ctx, int64(1), int64(1))
		assert.Error(t, err)
		assert.Equal(t, "", out)
		assert.Equal(t, "db down", err.Error())
	})
	t.Run("deletesocialmedia required userid", func(t *testing.T) {
		mock.ExpectExec(qry).
			WithArgs(int64(0), int64(1)).
			WillReturnError(errors.New("required userid"))
		out, err := dbtes.DeleteSocialMedia(ctx, int64(0), int64(1))
		assert.Error(t, err)
		assert.Equal(t, "", out)
		assert.Equal(t, "required userid", err.Error())
	})

	t.Run("deletesocialmedia required id", func(t *testing.T) {
		mock.ExpectExec(qry).
			WithArgs(int64(1), int64(0)).
			WillReturnError(errors.New("required id"))
		out, err := dbtes.DeleteSocialMedia(ctx, int64(1), int64(0))
		assert.Error(t, err)
		assert.Equal(t, "", out)
		assert.Equal(t, "required id", err.Error())
	})

	t.Run("deletesocialmedia success", func(t *testing.T) {
		mock.ExpectExec(qry).
			WithArgs(int64(1), int64(1)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		out, err := dbtes.DeleteSocialMedia(ctx, int64(1), int64(1))
		assert.NotNil(t, out)
		assert.NoError(t, err)
	})
}
