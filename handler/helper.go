package handler

import (
	"encoding/json"
	"fmt"
	"mygram/entity"
	"net/http"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var JWT_SIGNING_METHOD = jwt.SigningMethodHS256
var LogonUser *entity.User

type response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

const (
	Success              int = 200
	Success201           int = 201
	ErrorBadRequest      int = 400
	ErrorUnauthorized    int = 401
	ErrorForbidden       int = 403
	ErrorNotFound        int = 404
	ErrorDataHandleError int = 500
)

func EncryptPassword(pwd string) (string, error) {
	// Hashing the password with the default cost of 10
	securePassword, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(securePassword), nil
}

func WriteJsonResp(w http.ResponseWriter, status int, obj interface{}) {
	resp := response{
		Status: status,
		Data:   obj,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

/*func encryptToken(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], text)
	return ciphertext, nil
}

func decryptToken(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	return text, nil
}*/

func GetConnectionString() string {
	connString := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		Config.ConnectionString.Sqluser, Config.ConnectionString.Sqlpassword, Config.ConnectionString.Sqlserver, Config.ConnectionString.Sqlport, Config.ConnectionString.SqldbName)
	return connString
}
