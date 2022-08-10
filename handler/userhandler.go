package handler

import (
	"context"
	"encoding/json"
	"errors"
	"mygram/database"
	"mygram/entity"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct{}

func InstallUsersHandler(r *mux.Router) {
	api := UserHandler{}
	r.HandleFunc("/users/{action}", api.UsersHandler)
	r.HandleFunc("/users", api.UsersHandler).Queries("userId", "{userId}").Methods("PUT")
	r.HandleFunc("/users", api.UsersHandler)
}

type UserHandlerInterface interface {
	UsersHandler(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler() UserHandlerInterface {
	return &UserHandler{}
}

func (h *UserHandler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	action := params["action"]

	switch r.Method {
	case http.MethodPost:
		if action == "login" {
			loginUserHandler(w, r)
		} else if action == "register" {
			registerUsersHandler(w, r)
		}
	case http.MethodPut:
		updateUserHandler(w, r)
	case http.MethodDelete:
		deleteUserHandler(w, r)
	default:
		WriteJsonResp(w, ErrorNotFound, "PAGE NOT FOUND")
		return
	}
}

// loginUserHandler
// Method: POST
// Example: localhost/login
// JSON Body:
// {
// 	"email": "user@email.com",
// 	"password": "password"
// }
func loginUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	decoder := json.NewDecoder(r.Body)
	var inp entity.UserLogin
	if err := decoder.Decode(&inp); err != nil {
		w.Write([]byte("error decoding json body"))
		return
	}
	id, pw, err := database.SqlDatabase.Login(ctx, inp.Email)
	if err != nil {
		WriteJsonResp(w, ErrorDataHandleError, err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(pw), []byte(inp.Password))
	if err != nil {
		WriteJsonResp(w, ErrorUnauthorized, "UNAUTHORIZED")
		return
	}
	claims := entity.MyClaims{
		Iat: int(time.Now().UnixMilli()),
		Exp: int(time.Now().Add(time.Second * time.Duration(60)).UnixMilli()),
		Uid: id,
	}

	token := jwt.NewWithClaims(
		JWT_SIGNING_METHOD,
		claims,
	)

	tokenVal, err := token.SignedString([]byte(Config.SecretKey))
	if err != nil {
		WriteJsonResp(w, ErrorBadRequest, "BAD_REQUEST")
		return
	}
	retVal := map[string]string{
		"token": tokenVal,
	}
	WriteJsonResp(w, Success, retVal)
}

// registerUsersHandler
// Method: POST
// Example: localhost/register
// JSON Body:
// {
//		"username": "user1",
//		"email": "user@email.com",
//		"password": "password1",
//		"age": 22
// }
func registerUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	validate := validator.New()
	decoder := json.NewDecoder(r.Body)
	var inp entity.UserRegister
	if err := decoder.Decode(&inp); err != nil {
		WriteJsonResp(w, ErrorDataHandleError, err.Error())
		return
	}
	err := validate.Struct(inp)
	if err != nil {
		WriteJsonResp(w, ErrorBadRequest, err.Error())
		return
	}

	encrtedPwd, err := EncryptPassword(inp.Password)
	if err != nil {
		WriteJsonResp(w, ErrorDataHandleError, err.Error())
		return
	}
	inp.Password = encrtedPwd

	users, err := database.SqlDatabase.Register(ctx, inp)
	if err != nil {
		WriteJsonResp(w, ErrorDataHandleError, err.Error())
		return
	}
	WriteJsonResp(w, Success201, users)
}

// updateUserHandler
// Method: PUT
// Example: localhost/users
// JSON Body:
// {
//		"username": "user1",
//		"email": "user@email.com"
// }
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)
	userid := vars["userId"]
	id, err := strconv.ParseInt(userid, 10, 64)
	if err != nil {
		WriteJsonResp(w, ErrorBadRequest, err.Error())
		return
	}
	if id != LogonUser.ID {
		WriteJsonResp(w, ErrorBadRequest, errors.New("wrong ID").Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	validate := validator.New()
	var inp entity.UserUpdate
	if err := decoder.Decode(&inp); err != nil {
		WriteJsonResp(w, ErrorDataHandleError, err.Error())
		return
	}
	err = validate.Struct(inp)
	if err != nil {
		WriteJsonResp(w, ErrorBadRequest, err.Error())
		return
	}
	users, err := database.SqlDatabase.UpdateUser(ctx, id, inp.Email, inp.Username)
	if err != nil {
		WriteJsonResp(w, ErrorDataHandleError, err.Error())
		return
	}
	retVal := users.ToUserUpdateOutput()
	WriteJsonResp(w, Success, retVal)

}

// deleteUserHandler
// Method: DELETE
// Example: localhost/users
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id := LogonUser.ID
	users, err := database.SqlDatabase.DeleteUser(ctx, id)
	if err != nil {
		WriteJsonResp(w, ErrorDataHandleError, err.Error())
		return
	}
	retVal := map[string]string{
		"message": users,
	}
	WriteJsonResp(w, Success, retVal)

}
