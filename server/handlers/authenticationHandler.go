package handlers

import (
	"log"
	"net/http"
)

var (
	testUserName = "lidor"
	testPassword = "lidor"
)

type AuthenticationHandler struct {
	handlerLogger *log.Logger
}

func NewAuthenticationHandler(l *log.Logger) *AuthenticationHandler {
	return &AuthenticationHandler{handlerLogger: l}
}

func (crawlerHandler *AuthenticationHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

}
