package handler

import (
	"context"
	"encoding/json"
	"mygram/database"
	"mygram/entity"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type CommentHandler struct{}

func InstallCommentHandler(r *mux.Router) {
	api := CommentHandler{}
	r.HandleFunc("/comments/{id}", api.CommentsHandler)
	r.HandleFunc("/comments", api.CommentsHandler)
}

type CommentHandlerInterface interface {
	CommentsHandler(w http.ResponseWriter, r *http.Request)
}

func NewCommentHandler() CommentHandlerInterface {
	return &CommentHandler{}
}

func (h *CommentHandler) CommentsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	switch r.Method {
	case http.MethodGet:
		getCommentsHandler(w, r)
	case http.MethodPost:
		postCommentHandler(w, r)
	case http.MethodPut:
		updateCommentHandler(w, r, id)
	case http.MethodDelete:
		deleteCommentHandler(w, r, id)
	default:
		WriteJsonResp(w, ErrorNotFound, "PAGE NOT FOUND")
		return
	}
}

// getCommentsHandler
// Method: GET
// Example: localhost/comments
func getCommentsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	retVal, err := database.SqlDatabase.GetComments(ctx)
	if err != nil {
		WriteJsonResp(w, ErrorDataHandleError, err.Error())
		return
	}

	WriteJsonResp(w, Success, retVal)
}

// postCommentHandler
// Method: POST
// Example: localhost/comments
// JSON Body:
// {
// 	"message": "comment message",
// 	"photo_id": 1
// }
func postCommentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	validate := validator.New()
	decoder := json.NewDecoder(r.Body)
	var inp entity.CommentPost

	if err := decoder.Decode(&inp); err != nil {
		WriteJsonResp(w, ErrorDataHandleError, err.Error())
		return
	}
	err := validate.Struct(inp)
	if err != nil {
		WriteJsonResp(w, ErrorBadRequest, err.Error())
		return
	}
	c, err := database.SqlDatabase.PostComment(ctx, LogonUser.ID, inp)
	if err != nil {
		WriteJsonResp(w, ErrorDataHandleError, err.Error())
		return
	}

	retVal := c.ToCommentPostOutput()

	WriteJsonResp(w, Success201, retVal)
}

// updateCommentHandler
// Method: PUT
// Example: localhost/comments/1
// JSON Body:
// {
// 	"message": "comment message"
// }
func updateCommentHandler(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()

	if id != "" { // get by id
		if idInt, err := strconv.ParseInt(id, 10, 64); err == nil {
			validate := validator.New()
			decoder := json.NewDecoder(r.Body)
			var inp entity.CommentUpdate

			if err := decoder.Decode(&inp); err != nil {
				WriteJsonResp(w, ErrorDataHandleError, err.Error())
				return
			}
			err := validate.Struct(inp)
			if err != nil {
				WriteJsonResp(w, ErrorBadRequest, err.Error())
				return
			}
			c, err := database.SqlDatabase.GetCommentByID(ctx, idInt)
			if err != nil {
				WriteJsonResp(w, ErrorDataHandleError, err.Error())
				return
			}
			if c.UserID != LogonUser.ID {
				WriteJsonResp(w, ErrorUnauthorized, "UNAUTHORIZED")
				return
			}

			p, err := database.SqlDatabase.UpdateComment(ctx, LogonUser.ID, idInt, inp.Message)
			if err != nil {
				WriteJsonResp(w, ErrorDataHandleError, err.Error())
				return
			}
			retVal := p.ToCommentUpdateOutput()
			WriteJsonResp(w, Success, retVal)
		}
	}
}

// deleteCommentHandler
// Method: DELETE
// Example: localhost/comments/1
func deleteCommentHandler(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()
	if id != "" {
		if idInt, err := strconv.ParseInt(id, 10, 64); err == nil {
			c, err := database.SqlDatabase.GetCommentByID(ctx, idInt)
			if err != nil {
				WriteJsonResp(w, ErrorDataHandleError, err.Error())
				return
			}
			if c.UserID != LogonUser.ID {
				WriteJsonResp(w, ErrorUnauthorized, "UNAUTHORIZED")
				return
			}
			msg, err := database.SqlDatabase.DeleteComment(ctx, LogonUser.ID, idInt)
			if err != nil {
				WriteJsonResp(w, ErrorDataHandleError, err.Error())
				return
			}
			retVal := map[string]string{
				"message": msg,
			}
			WriteJsonResp(w, Success, retVal)

		}
	}
}
