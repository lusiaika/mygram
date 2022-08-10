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

func TestDatabase_PostPhoto(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	db, mock := NewMock()
	defer db.Close()
	dbtes := Database{
		SqlDb: db,
	}
	qry := "insert into photos (title, caption, photourl, userid, createdat, updatedat) values (@title, @caption, @photourl, @userid, @createdat, @updatedat); select id, title, caption, photourl, userid, createdat from photos where id = SCOPE_IDENTITY()"

	inp := entity.PhotoPost{
		Title:    "Foto Kopi",
		Caption:  "Foto kopi doang beneran",
		PhotoUrl: "https://imageurl.com/fotokopi.jpg",
	}

	t.Run("postphoto database down", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs(inp.Title, inp.Caption, inp.PhotoUrl, int64(1), time.Now(), time.Now()).
			WillReturnError(errors.New("db down"))
		out, err := dbtes.PostPhoto(ctx, int64(1), inp)
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "db down", err.Error())
	})

	t.Run("postphoto required userid", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs(inp.Title, inp.Caption, inp.PhotoUrl, int64(0), time.Now(), time.Now()).
			WillReturnError(errors.New("required userid"))
		out, err := dbtes.PostPhoto(ctx, int64(0), inp)
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "required userid", err.Error())
	})

	t.Run("postphoto success", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "title", "caption", "photourl", "userid", "createdat"}).
			AddRow(1, "Foto Kopi", "Foto kopi doang beneran", "http://imageurl.com/fotokopi.jpg", 1, time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs(inp.Title, inp.Caption, inp.PhotoUrl, int64(1), time.Now(), time.Now()).
			WillReturnRows(rows)
		out, err := dbtes.PostPhoto(ctx, int64(1), inp)
		assert.NotNil(t, out)
		assert.NoError(t, err)
	})
}

func TestDatabase_GetPhotos(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	db, mock := NewMock()
	defer db.Close()
	dbtes := Database{
		SqlDb: db,
	}
	var qry strings.Builder
	qry.WriteString("select p.id, p.title, p.caption, p.photourl, p.userid, p.createdat, p.updatedat, u.email, u.username from photos p")
	qry.WriteString(" join users u on p.userid=u.id")
	t.Run("getphotos database down", func(t *testing.T) {
		mock.ExpectQuery(qry.String()).
			WillReturnError(errors.New("db down"))
		out, err := dbtes.GetPhotos(ctx)
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "db down", err.Error())
	})

	t.Run("getphotos success", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "title", "caption", "photourl", "userid", "createdat", "updatedat", "email", "username"}).
			AddRow(1, "Foto Kopi", "Foto kopi doang beneran", "http://imageurl.com/fotokopi.jpg", 1, time.Now(), time.Now(), "deadapeipit@email.com", "deadapeipit")

		mock.ExpectQuery(qry.String()).WillReturnRows(rows)
		out, err := dbtes.GetPhotos(ctx)
		assert.NotNil(t, out)
		assert.NoError(t, err)
	})
}

func TestDatabase_UpdatePhoto(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	db, mock := NewMock()
	defer db.Close()
	dbtes := Database{
		SqlDb: db,
	}
	qry := "update photos set title=@title, caption=@caption, photourl=@photourl, updatedat=@updatedat where id = @ID and userid = @userid; select id, title, caption, photourl, userid, updatedat from photos where id = @ID"
	inp := entity.PhotoPost{
		Title:    "Foto Kopi",
		Caption:  "Foto kopi doang beneran",
		PhotoUrl: "https://imageurl.com/fotokopi.jpg",
	}

	t.Run("updatephoto database down", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs(inp.Title, inp.Caption, inp.PhotoUrl, time.Now(), int64(1), int64(1)).
			WillReturnError(errors.New("db down"))
		out, err := dbtes.UpdatePhoto(ctx, int64(1), int64(1), inp)
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "db down", err.Error())
	})

	t.Run("updatephoto required id", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs(inp.Title, inp.Caption, inp.PhotoUrl, time.Now(), int64(1), int64(0)).
			WillReturnError(errors.New("required id"))
		out, err := dbtes.UpdatePhoto(ctx, int64(1), int64(0), inp)
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "required id", err.Error())
	})

	t.Run("updatephoto required userid", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs(inp.Title, inp.Caption, inp.PhotoUrl, time.Now(), int64(0), int64(1)).
			WillReturnError(errors.New("required userid"))
		out, err := dbtes.UpdatePhoto(ctx, int64(0), int64(1), inp)
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "required userid", err.Error())
	})

	t.Run("updatephoto success", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "title", "caption", "photourl", "userid", "updatedat"}).
			AddRow(1, "Foto Kopi", "Foto kopi doang beneran", "http://imageurl.com/fotokopi.jpg", 1, time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs(inp.Title, inp.Caption, inp.PhotoUrl, time.Now(), int64(1), int64(1)).
			WillReturnRows(rows)
		out, err := dbtes.UpdatePhoto(ctx, int64(1), int64(1), inp)
		assert.NotNil(t, out)
		assert.NoError(t, err)
	})
}

func TestDatabase_DeletePhoto(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	db, mock := NewMock()
	defer db.Close()
	dbtes := Database{
		SqlDb: db,
	}
	qry := "delete from comments where photoid=@id and userid=@userid; delete from photos where id=@id and userid=@userid"
	t.Run("deletephoto database down", func(t *testing.T) {
		mock.ExpectExec(qry).
			WithArgs(int64(1), int64(1)).
			WillReturnError(errors.New("db down"))
		out, err := dbtes.DeletePhoto(ctx, int64(1), int64(1))
		assert.Error(t, err)
		assert.Equal(t, "", out)
		assert.Equal(t, "db down", err.Error())
	})

	t.Run("deletephoto required userid", func(t *testing.T) {
		mock.ExpectExec(qry).
			WithArgs(int64(0), int64(1)).
			WillReturnError(errors.New("required userid"))
		out, err := dbtes.DeletePhoto(ctx, int64(0), int64(1))
		assert.Error(t, err)
		assert.Equal(t, "", out)
		assert.Equal(t, "required userid", err.Error())
	})

	t.Run("deletephoto required id", func(t *testing.T) {
		mock.ExpectExec(qry).
			WithArgs(int64(1), int64(0)).
			WillReturnError(errors.New("required id"))
		out, err := dbtes.DeletePhoto(ctx, int64(1), int64(0))
		assert.Error(t, err)
		assert.Equal(t, "", out)
		assert.Equal(t, "required id", err.Error())
	})

	t.Run("deletephoto success", func(t *testing.T) {
		mock.ExpectExec(qry).
			WithArgs(int64(1), int64(1)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		out, err := dbtes.DeletePhoto(ctx, int64(1), int64(1))
		assert.NotNil(t, out)
		assert.NoError(t, err)
	})
}
