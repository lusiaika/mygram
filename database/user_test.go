package database

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"mygram/entity"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestDatabase_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	db, mock := NewMock()
	defer db.Close()
	dbtes := Database{
		SqlDb: db,
	}
	qry := "select id, password from users where username = @username"
	t.Run("login database down", func(t *testing.T) {
		mock.ExpectQuery(qry).
			WithArgs("deadapeipit").
			WillReturnError(errors.New("db down"))
		id, pass, err := dbtes.Login(ctx, "deadapeipit")
		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
		assert.Equal(t, "", pass)
		assert.Equal(t, "db down", err.Error())
	})

	t.Run("login success", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "password"}).
			AddRow(1, "deadapeipit")

		mock.ExpectQuery(qry).
			WithArgs("deadapeipit").
			WillReturnRows(rows)
		id, pass, err := dbtes.Login(ctx, "deadapeipit")
		assert.NotEqual(t, int64(0), id)
		assert.NotEqual(t, "", pass)
		assert.NoError(t, err)
	})
}

func TestDatabase_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	db, mock := NewMock()
	defer db.Close()
	dbtes := Database{
		SqlDb: db,
	}
	qry := "select id, username, email, password, age, createdat, updatedat from users where id = @ID"
	t.Run("getuserbyid database down", func(t *testing.T) {
		mock.ExpectQuery(qry).
			WithArgs(int64(1)).
			WillReturnError(errors.New("db down"))
		out, err := dbtes.GetUserByID(ctx, int64(1))
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "db down", err.Error())
	})

	t.Run("getuserbyid required userid", func(t *testing.T) {
		mock.ExpectQuery(qry).
			WithArgs(int64(0)).
			WillReturnError(errors.New("required userid"))
		out, err := dbtes.GetUserByID(ctx, int64(0))
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "required userid", err.Error())
	})

	t.Run("getuserbyid success", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "username", "email", "password", "age", "createdat", "updatedat"}).
			AddRow(1, "deadapeipit", "deadapeipit@github.com", "$2a$10$rcIrmHvODKlw91zkIVeEGeAomU47EBbAveY8//HCvYK7cqrd23gx2", 22, time.Now(), time.Now())

		mock.ExpectQuery(qry).
			WithArgs(int64(1)).
			WillReturnRows(rows)
		out, err := dbtes.GetUserByID(ctx, int64(1))
		assert.NotNil(t, out)
		assert.NoError(t, err)
	})
}

func TestDatabase_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	db, mock := NewMock()
	defer db.Close()
	dbtes := Database{
		SqlDb: db,
	}
	qry := "update users set email=@email, username=@username, updatedat=@updatedat where id = @ID; select ID, email, username, age, updatedat from users where id = @ID"
	t.Run("updateuser database down", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs("deadapeipit@github.com", "deadapeipit", time.Now(), int64(1)).
			WillReturnError(errors.New("db down"))
		out, err := dbtes.UpdateUser(ctx, int64(1), "deadapeipit@github.com", "deadapeipit")
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "db down", err.Error())
	})

	t.Run("updateuser required userid", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs("deadapeipit@github.com", "deadapeipit", time.Now(), int64(0)).
			WillReturnError(errors.New("required userid"))
		out, err := dbtes.UpdateUser(ctx, int64(0), "deadapeipit@github.com", "deadapeipit")
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "required userid", err.Error())
	})

	t.Run("updateuser success", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "username", "email", "age", "updatedat"}).
			AddRow(1, "deadapeipit", "deadapeipit@github.com", 22, time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs("deadapeipit@github.com", "deadapeipit", time.Now(), int64(1)).
			WillReturnRows(rows)
		out, err := dbtes.UpdateUser(ctx, int64(1), "deadapeipit@github.com", "deadapeipit")
		assert.NotNil(t, out)
		assert.NoError(t, err)
	})
}

func TestDatabase_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	db, mock := NewMock()
	defer db.Close()
	dbtes := Database{
		SqlDb: db,
	}
	inp := entity.UserRegister{
		Username: "deadapeipit",
		Email:    "deadapeipit@email.com",
		Password: "passdeadapeipit",
		Age:      29,
	}
	qry := "insert into users (username, email, password, age, createdat, updatedat) values (@username, @email, @password, @age, @createdat, @updatedat)"
	t.Run("register database down", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs(inp.Username, inp.Email, inp.Password, inp.Age, time.Now(), time.Now()).
			WillReturnError(errors.New("db down"))
		out, err := dbtes.Register(ctx, inp)
		assert.Error(t, err)
		assert.Nil(t, out)
		assert.Equal(t, "db down", err.Error())
	})

	t.Run("register success", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "username", "email", "age"}).
			AddRow(1, "deadapeipit", "deadapeipit@github.com", 22)

		mock.ExpectQuery(regexp.QuoteMeta(qry)).
			WithArgs(inp.Username, inp.Email, inp.Password, inp.Age, time.Now(), time.Now()).
			WillReturnRows(rows)
		out, err := dbtes.Register(ctx, inp)
		assert.NotNil(t, out)
		assert.NoError(t, err)
	})
}

func TestDatabase_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	db, mock := NewMock()
	defer db.Close()
	dbtes := Database{
		SqlDb: db,
	}
	qry := "delete from socialmedias where userid=@id; delete from photos where userid=@id; delete from comments where userid=@id; delete from users where id=@id"
	t.Run("deleteuser database down", func(t *testing.T) {
		mock.ExpectExec(qry).
			WithArgs(int64(1)).
			WillReturnError(errors.New("db down"))
		out, err := dbtes.DeleteUser(ctx, int64(1))
		assert.Error(t, err)
		assert.Equal(t, "", out)
		assert.Equal(t, "db down", err.Error())
	})

	t.Run("deleteuser required userid", func(t *testing.T) {
		mock.ExpectExec(qry).
			WithArgs(int64(0)).
			WillReturnError(errors.New("required userid"))
		out, err := dbtes.DeleteUser(ctx, int64(0))
		assert.Error(t, err)
		assert.Equal(t, "", out)
		assert.Equal(t, "required userid", err.Error())
	})

	t.Run("deleteuser success", func(t *testing.T) {
		mock.ExpectExec(qry).
			WithArgs(int64(1)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		out, err := dbtes.DeleteUser(ctx, int64(1))
		assert.NotNil(t, out)
		assert.NoError(t, err)
	})
}
