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

type SocialMediaHandler struct{}

func InstallSocialMediaHandler(r *mux.Router) {
	api := SocialMediaHandler{}
	r.HandleFunc("/socialmedias/{id}", api.SocialMediasHandler)
	r.HandleFunc("/socialmedias", api.SocialMediasHandler)
}

type SocialMediaHandlerInterface interface {
	SocialMediasHandler(w http.ResponseWriter, r *http.Request)
}

func NewSocialMediaHandler() SocialMediaHandlerInterface {
	return &SocialMediaHandler{}
}

func (h *SocialMediaHandler) SocialMediasHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	switch r.Method {
	case http.MethodGet:
		getSocialMediasHandler(w, r)
	case http.MethodPost:
		postSocialMediaHandler(w, r)
	case http.MethodPut:
		updateSocialMediaHandler(w, r, id)
	case http.MethodDelete:
		deleteSocialMediaHandler(w, r, id)
	default:
		WriteJsonResp(w, ErrorNotFound, "PAGE NOT FOUND")
		return
	}
}

// getSocialMediasHandler
// Method: GET
// Example: localhost/socialmedias
func getSocialMediasHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	retVal, err := database.SqlDatabase.GetSocialMedias(ctx)
	if err != nil {
		WriteJsonResp(w, ErrorDataHandleError, err.Error())
		return
	}

	WriteJsonResp(w, Success, retVal)
}

// postSocialMediaHandler
// Method: POST
// Example: localhost/socialmedias
// JSON Body:
// {
// 	"name": "social media name",
// 	"social_media_url": "https://domainsocialmedia.com/user",
// 	"profile_image_url": "https://domainsocialmedia.com/userimage.jpg"
// }
func postSocialMediaHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	validate := validator.New()
	decoder := json.NewDecoder(r.Body)
	var inp entity.SocialMediaPost
	if err := decoder.Decode(&inp); err != nil {
		WriteJsonResp(w, ErrorDataHandleError, err.Error())
		return
	}
	err := validate.Struct(inp)
	if err != nil {
		WriteJsonResp(w, ErrorBadRequest, err.Error())
		return
	}
	p, err := database.SqlDatabase.PostSocialMedia(ctx, LogonUser.ID, inp)
	if err != nil {
		WriteJsonResp(w, ErrorDataHandleError, err.Error())
		return
	}

	retVal := p.ToSocialMediaPostOutput()

	WriteJsonResp(w, Success201, retVal)
}

// updateSocialMediaHandler
// Method: PUT
// Example: localhost/socialmedias/1
// JSON Body:
// {
// 	"name": "social media name",
// 	"social_media_url": "https://domainsocialmedia.com/user",
// 	"profile_image_url": "https://domainsocialmedia.com/userimage.jpg"
// }
func updateSocialMediaHandler(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()

	if id != "" { // get by id
		if idInt, err := strconv.ParseInt(id, 10, 64); err == nil {
			validate := validator.New()
			decoder := json.NewDecoder(r.Body)
			var inp entity.SocialMediaPost
			if err := decoder.Decode(&inp); err != nil {
				WriteJsonResp(w, ErrorDataHandleError, err.Error())
				return
			}
			err := validate.Struct(inp)
			if err != nil {
				WriteJsonResp(w, ErrorBadRequest, err.Error())
				return
			}
			c, err := database.SqlDatabase.GetSocialMediaByID(ctx, idInt)
			if err != nil {
				WriteJsonResp(w, ErrorDataHandleError, err.Error())
				return
			}
			if c.UserID != LogonUser.ID {
				WriteJsonResp(w, ErrorUnauthorized, "UNAUTHORIZED")
				return
			}

			p, err := database.SqlDatabase.UpdateSocialMedia(ctx, LogonUser.ID, idInt, inp)
			if err != nil {
				WriteJsonResp(w, ErrorDataHandleError, err.Error())
				return
			}
			retVal := p.ToSocialMediaUpdateOutput()
			WriteJsonResp(w, Success, retVal)
		}
	}
}

// deleteSocialMediaHandler
// Method: DELETE
// Example: localhost/socialmedias/1
func deleteSocialMediaHandler(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()
	if id != "" {
		if idInt, err := strconv.ParseInt(id, 10, 64); err == nil {
			c, err := database.SqlDatabase.GetSocialMediaByID(ctx, idInt)
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
