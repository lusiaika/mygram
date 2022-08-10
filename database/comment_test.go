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

func TestDatabase_PostComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	db, mock := NewMock()
	defer db.Close()
	dbtes := Database{
		SqlDb: db,
	}
	qry := "insert into comments (message, photoid, userid, createdat, updatedat) values (@message, @photoid, @userid, @createdat, @updatedat); select id, message, photoid, userid, createdat from comments where id = SCOPE_IDENTITY()"

	inp := entity.CommentPost{
		Message: "Foto Kopi",
		PhotoID: 1,
	}

	t.Run("postcomment database down", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs(inp.Message, inp.PhotoID, int64(1), time.Now(), time.Now()).
			WillReturnError(errors.New("db down"))
		out, err := dbtes.PostComment(ctx, int64(1), inp)
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "db down", err.Error())
	})

	t.Run("postcomment required userid", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs(inp.Message, inp.PhotoID, int64(0), time.Now(), time.Now()).
			WillReturnError(errors.New("required userid"))
		out, err := dbtes.PostComment(ctx, int64(0), inp)
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "required userid", err.Error())
	})

	t.Run("postcomment success", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "message", "photoid", "userid", "createdat"}).
			AddRow(1, "Message nya apa", 1, 1, time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs(inp.Message, inp.PhotoID, int64(1), time.Now(), time.Now()).
			WillReturnRows(rows)
		out, err := dbtes.PostComment(ctx, int64(1), inp)
		assert.NotNil(t, out)
		assert.NoError(t, err)
	})
}

func TestDatabase_GetComments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	db, mock := NewMock()
	defer db.Close()
	dbtes := Database{
		SqlDb: db,
	}
	var qry strings.Builder
	qry.WriteString("select c.id, c.message, c.photoid, c.userid, c.createdat, c.updatedat,")
	qry.WriteString(" p.title, p.caption, p.photourl,")
	qry.WriteString(" u.email, u.username from comments c")
	qry.WriteString(" join photos p on c.photoid=p.id")
	qry.WriteString(" join users u on c.userid=u.id")
	t.Run("getcomments database down", func(t *testing.T) {
		mock.ExpectQuery(qry.String()).
			WillReturnError(errors.New("db down"))
		out, err := dbtes.GetComments(ctx)
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "db down", err.Error())
	})

	t.Run("getcomments success", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "message", "photoid", "userid", "createdat", "updatedat", "title", "caption", "photourl", "email", "username"}).
			AddRow(1, "Message nya apa", 1, 1, time.Now(), time.Now(), "Title photo", "Caption Photoo", "http://photourl.com/photourl.jpg", "deadapeipit@email.com", "deadapeipit")

		mock.ExpectQuery(qry.String()).WillReturnRows(rows)
		out, err := dbtes.GetComments(ctx)
		assert.NotNil(t, out)
		assert.NoError(t, err)
	})
}

func TestDatabase_GetCommentByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	db, mock := NewMock()
	defer db.Close()
	dbtes := Database{
		SqlDb: db,
	}
	qry := "select id, message, photoid, userid, createdat, updatedat from comments where id = @ID"
	t.Run("getcommentbyid database down", func(t *testing.T) {
		mock.ExpectQuery(qry).
			WithArgs(int64(1)).
			WillReturnError(errors.New("db down"))
		out, err := dbtes.GetCommentByID(ctx, int64(1))
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "db down", err.Error())
	})

	t.Run("getcommentbyid required id", func(t *testing.T) {
		mock.ExpectQuery(qry).
			WithArgs(int64(0)).
			WillReturnError(errors.New("required id"))
		out, err := dbtes.GetCommentByID(ctx, int64(0))
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "required id", err.Error())
	})

	t.Run("getcommentbyid success", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "message", "photoid", "userid", "createdat", "updatedat"}).
			AddRow(1, "Message nya apa", 1, 1, time.Now(), time.Now())

		mock.ExpectQuery(qry).
			WithArgs(int64(1)).
			WillReturnRows(rows)
		out, err := dbtes.GetCommentByID(ctx, int64(1))
		assert.NotNil(t, out)
		assert.NoError(t, err)
	})
}

func TestDatabase_UpdateComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	db, mock := NewMock()
	defer db.Close()
	dbtes := Database{
		SqlDb: db,
	}
	qry := "update comments set message=@message, updatedat=@updatedat where id = @ID and userid = userid; select id, userid, photoid, message, updatedat from comments where id = @ID"
	inp := entity.CommentPost{
		Message: "Foto Kopi",
		PhotoID: 1,
	}
	t.Run("updatecomment database down", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs(inp.Message, time.Now(), int64(1), int64(1)).
			WillReturnError(errors.New("db down"))
		out, err := dbtes.UpdateComment(ctx, int64(1), int64(1), inp.Message)
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "db down", err.Error())
	})

	t.Run("updatecomment required id", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs(inp.Message, time.Now(), int64(1), int64(0)).
			WillReturnError(errors.New("required id"))
		out, err := dbtes.UpdateComment(ctx, int64(1), int64(0), inp.Message)
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "required id", err.Error())
	})

	t.Run("updatecomment required userid", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs(inp.Message, time.Now(), int64(0), int64(1)).
			WillReturnError(errors.New("required userid"))
		out, err := dbtes.UpdateComment(ctx, int64(0), int64(1), inp.Message)
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "required userid", err.Error())
	})

	t.Run("updatecomment success", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "userid", "photoid", "message", "updatedat"}).
			AddRow(1, 1, 1, "Foto kopi doang beneran cuk", time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs(inp.Message, time.Now(), int64(1), int64(1)).
			WillReturnRows(rows)
		out, err := dbtes.UpdateComment(ctx, int64(1), int64(1), inp.Message)
		assert.NotNil(t, out)
		assert.NoError(t, err)
	})
}

func TestDatabase_DeleteComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	db, mock := NewMock()
	defer db.Close()
	dbtes := Database{
		SqlDb: db,
	}
	qry := "delete from comments where id=@id and userid = @userid"
	t.Run("deletecomment database down", func(t *testing.T) {
		mock.ExpectExec(qry).
			WithArgs(int64(1), int64(1)).
			WillReturnError(errors.New("db down"))
		out, err := dbtes.DeleteComment(ctx, int64(1), int64(1))
		assert.Error(t, err)
		assert.Equal(t, "", out)
		assert.Equal(t, "db down", err.Error())
	})

	t.Run("deletecomment required userid", func(t *testing.T) {
		mock.ExpectExec(qry).
			WithArgs(int64(0), int64(1)).
			WillReturnError(errors.New("required userid"))
		out, err := dbtes.DeleteComment(ctx, int64(0), int64(1))
		assert.Error(t, err)
		assert.Equal(t, "", out)
		assert.Equal(t, "required userid", err.Error())
	})

	t.Run("deletecomment required id", func(t *testing.T) {
		mock.ExpectExec(qry).
			WithArgs(int64(1), int64(0)).
			WillReturnError(errors.New("required id"))
		out, err := dbtes.DeleteComment(ctx, int64(1), int64(0))
		assert.Error(t, err)
		assert.Equal(t, "", out)
		assert.Equal(t, "required id", err.Error())
	})

	t.Run("deletecomment success", func(t *testing.T) {
		mock.ExpectExec(qry).
			WithArgs(int64(1), int64(1)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		out, err := dbtes.DeleteComment(ctx, int64(1), int64(1))
		assert.NotNil(t, out)
		assert.NoError(t, err)
	})
}
