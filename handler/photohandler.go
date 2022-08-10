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

type PhotoHandler struct{}

func InstallPhotosHandler(r *mux.Router) {
	api := PhotoHandler{}
	r.HandleFunc("/photos/{id}", api.PhotosHandler)
	r.HandleFunc("/photos", api.PhotosHandler)
}

type PhotoHandlerInterface interface {
	UsersHandler(w http.ResponseWriter, r *http.Request)
}

func NewPhotoHandler() PhotoHandlerInterface {
	return &UserHandler{}
}

func (h *PhotoHandler) PhotosHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	switch r.Method {
	case http.MethodGet:
		getPhotosHandler(w, r)
	case http.MethodPost:
		postPhotoHandler(w, r)
	case http.MethodPut:
		updatePhotoHandler(w, r, id)
	case http.MethodDelete:
		deletePhotoHandler(w, r, id)
	default:
		WriteJsonResp(w, ErrorNotFound, "PAGE NOT FOUND")
		return
	}
}

// getPhotosHandler
// Method: GET
// Example: localhost/photos
func getPhotosHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	retVal, err := database.SqlDatabase.GetPhotos(ctx)
	if err != nil {
		WriteJsonResp(w, ErrorDataHandleError, err.Error())
		return
	}

	WriteJsonResp(w, Success, retVal)
}

// postPhotoHandler
// Method: POST
// Example: localhost/photos
// JSON Body:
// {
// 	"title": "title photo",
// 	"caption": "caption photo",
// 	"photo_url": "https://photo.domain.com"
// }
func postPhotoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	validate := validator.New()
	decoder := json.NewDecoder(r.Body)
	var inp entity.PhotoPost
	if err := decoder.Decode(&inp); err != nil {
		WriteJsonResp(w, ErrorDataHandleError, err.Error())
		return
	}
	err := validate.Struct(inp)
	if err != nil {
		WriteJsonResp(w, ErrorBadRequest, err.Error())
		return
	}
	p, err := database.SqlDatabase.PostPhoto(ctx, LogonUser.ID, inp)
	if err != nil {
		WriteJsonResp(w, ErrorDataHandleError, err.Error())
		return
	}

	retVal := p.ToPhotoPostOutput()
	WriteJsonResp(w, Success201, retVal)
}

// updatePhotoHandler
// Method: PUT
// Example: localhost/photos/1
// JSON Body:
// {
// 	"title": "title photo",
// 	"caption": "caption photo",
// 	"photo_url": "https://photo.domain.com"
// }
func updatePhotoHandler(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()

	if id != "" { // get by id
		if idInt, err := strconv.ParseInt(id, 10, 64); err == nil {
			validate := validator.New()
			decoder := json.NewDecoder(r.Body)
			var inp entity.PhotoPost
			if err := decoder.Decode(&inp); err != nil {
				WriteJsonResp(w, ErrorDataHandleError, err.Error())
				return
			}

			err := validate.Struct(inp)
			if err != nil {
				WriteJsonResp(w, ErrorBadRequest, err.Error())
				return
			}
			c, err := database.SqlDatabase.GetPhotoByID(ctx, idInt)
			if err != nil {
				WriteJsonResp(w, ErrorDataHandleError, err.Error())
				return
			}

			if c.UserID != LogonUser.ID {
				WriteJsonResp(w, ErrorUnauthorized, "UNAUTHORIZED")
				return
			}

			p, err := database.SqlDatabase.UpdatePhoto(ctx, LogonUser.ID, idInt, inp)
			if err != nil {
				WriteJsonResp(w, ErrorDataHandleError, err.Error())
				return
			}
			retVal := p.ToPhotoUpdateOutput()
			WriteJsonResp(w, Success, retVal)
		}
	}
}

// deletePhotoHandler
// Method: DELETE
// Example: localhost/photos/1
func deletePhotoHandler(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()
	if id != "" {
		if idInt, err := strconv.ParseInt(id, 10, 64); err == nil {
			c, err := database.SqlDatabase.GetPhotoByID(ctx, idInt)
			if err != nil {
				WriteJsonResp(w, ErrorDataHandleError, err.Error())
				return
			}
			if c.UserID != LogonUser.ID {
				WriteJsonResp(w, ErrorUnauthorized, "UNAUTHORIZED")
				return
			}
			msg, err := database.SqlDatabase.DeletePhoto(ctx, LogonUser.ID, idInt)
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
