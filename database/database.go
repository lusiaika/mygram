package database

import (
	"context"
	"database/sql"
	"fmt"
	"mygram/entity"

	_ "github.com/denisenkom/go-mssqldb"
)

type DatabaseIface interface {
	CloseConnection()
	Login(ctx context.Context, userName string) (int64, string, error)
	GetUserByID(ctx context.Context, userid int64) (*entity.User, error)
	Register(ctx context.Context, user entity.UserRegister) (*entity.UserRegisterResp, error)
	UpdateUser(ctx context.Context, userid int64, email string, username string) (*entity.User, error)
	DeleteUser(ctx context.Context, userId int64) (string, error)

	GetPhotos(ctx context.Context) ([]entity.PhotoGetOutput, error)
	GetPhotoByID(ctx context.Context, id int64) (*entity.Photo, error)
	PostPhoto(ctx context.Context, userid int64, photo entity.PhotoPost) (*entity.Photo, error)
	UpdatePhoto(ctx context.Context, userid int64, id int64, photo entity.PhotoPost) (*entity.Photo, error)
	DeletePhoto(ctx context.Context, userid int64, id int64) (string, error)

	GetComments(ctx context.Context) ([]entity.CommentGetOutput, error)
	GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error)
	PostComment(ctx context.Context, userid int64, comment entity.CommentPost) (*entity.Comment, error)
	UpdateComment(ctx context.Context, userid int64, id int64, message string) (*entity.Comment, error)
	DeleteComment(ctx context.Context, userid int64, id int64) (string, error)

	GetSocialMedias(ctx context.Context) ([]entity.SocialMediaGetOutput, error)
	GetSocialMediaByID(ctx context.Context, id int64) (*entity.SocialMedia, error)
	PostSocialMedia(ctx context.Context, userid int64, socialmedia entity.SocialMediaPost) (*entity.SocialMedia, error)
	UpdateSocialMedia(ctx context.Context, userid int64, id int64, socialmedia entity.SocialMediaPost) (*entity.SocialMedia, error)
	DeleteSocialMedia(ctx context.Context, userid int64, id int64) (string, error)
}

type Database struct {
	SqlDb *sql.DB
}

var SqlDatabase DatabaseIface

func NewSqlConnection(connectionString string) DatabaseIface {

	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		fmt.Printf("[mssql] Error connecting to SQL Server: %v", err)
	}
	s := Database{}
	s.SqlDb = db
	s.SqlDb.SetMaxIdleConns(25)
	s.SqlDb.SetMaxOpenConns(25)

	return &s
}

func (d *Database) CloseConnection() {
	fmt.Println("connection closed!")
	d.SqlDb.Close()
}
