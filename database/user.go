package database

import (
	"context"
	"database/sql"
	"mygram/entity"
	"time"
)

func (s *Database) Login(ctx context.Context, email string) (userid int64, resultpassword string, err error) {
	qry := "select id, password from users where email = @email"
	rows, err := s.SqlDb.QueryContext(ctx, qry,
		sql.Named("email", email))
	if err != nil {
		return 0, "", err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&userid,
			&resultpassword,
		)
		if err != nil {
			return 0, "", err
		}
	}
	return userid, resultpassword, nil
}

func (s *Database) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	result := &entity.User{}

	rows, err := s.SqlDb.QueryContext(ctx, "select id, username, email, password, age, createdat, updatedat from users where id = @ID",
		sql.Named("ID", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.Username,
			&result.Email,
			&result.Password,
			&result.Age,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (s *Database) UpdateUser(ctx context.Context, id int64, email string, username string) (*entity.User, error) {
	result := &entity.User{}
	now := time.Now()
	qry := "update users set email=@email, username=@username, updatedat=@updatedat where id = @ID; select ID, email, username, age, updatedat from users where id = @ID"
	rows, err := s.SqlDb.QueryContext(ctx, qry,
		sql.Named("email", email),
		sql.Named("username", username),
		sql.Named("updatedat", now),
		sql.Named("ID", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&result.ID,
			&result.Email,
			&result.Username,
			&result.Age,
			&result.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (s *Database) Register(ctx context.Context, i entity.UserRegister) (*entity.UserRegisterResp, error) {
	result := &entity.UserRegisterResp{}
	qry := "insert into users (username, email, password, age, createdat, updatedat) values (@username, @email, @password, @age, @createdat, @updatedat); select top 1 age,email,id,username from users where email = @email order by id desc"
	now := time.Now()
	rows, err := s.SqlDb.QueryContext(ctx, qry,
		sql.Named("username", i.Username),
		sql.Named("email", i.Email),
		sql.Named("password", i.Password),
		sql.Named("age", i.Age),
		sql.Named("createdat", now),
		sql.Named("updatedat", now))
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&result.Age,
			&result.Email,
			&result.ID,
			&result.Username,
		)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (s *Database) DeleteUser(ctx context.Context, id int64) (string, error) {
	var result string
	qry := "delete from socialmedias where userid=@id; delete from photos where userid=@id; delete from comments where userid=@id; delete from users where id=@id"
	_, err := s.SqlDb.ExecContext(ctx, qry,
		sql.Named("id", id))
	if err != nil {
		return "", err
	}

	result = "Your account has been successfully deleted"

	return result, nil
}
