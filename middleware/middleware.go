package middleware

import (
	"context"
	"fmt"
	"mygram/database"
	h "mygram/handler"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func SecureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/login") ||
			strings.Contains(r.URL.Path, "/register") {
			next.ServeHTTP(w, r)
			return
		}

		// Check to see if this request can go thru
		auth := r.Header.Get("Authorization")
		if len(auth) == 0 {
			h.WriteJsonResp(w, h.ErrorForbidden, "FORBIDDEN")
			return
		}

		splitToken := strings.Split(auth, "Bearer ")
		if len(splitToken) != 2 {
			h.WriteJsonResp(w, h.ErrorForbidden, "FORBIDDEN")
			return
		}

		accessToken := splitToken[1]
		if len(accessToken) == 0 {
			h.WriteJsonResp(w, h.ErrorForbidden, "FORBIDDEN")
			return
		}

		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("signing method invalid")
			} else if method != h.JWT_SIGNING_METHOD {
				return nil, fmt.Errorf("signing method invalid")
			}

			return []byte(h.Config.SecretKey), nil
		})
		if err != nil {
			e, ok := err.(*jwt.ValidationError)
			if !ok || ok && e.Errors&jwt.ValidationErrorIssuedAt == 0 { // Don't report error that token used before issued.
				h.WriteJsonResp(w, h.ErrorBadRequest, "BAD REQUEST")
				return
			}
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok { //|| !token.Valid {
			h.WriteJsonResp(w, h.ErrorBadRequest, "BAD REQUEST")
			return
		}

		uid := claims["uid"].(float64)

		userID := int64(uid)

		l, err := database.SqlDatabase.GetUserByID(context.Background(), userID)
		if err != nil {
			h.WriteJsonResp(w, h.ErrorDataHandleError, err)
			return
		}
		//Set logonuser
		h.LogonUser = l
		//fmt.Println(uid)

		next.ServeHTTP(w, r)
	})
}
