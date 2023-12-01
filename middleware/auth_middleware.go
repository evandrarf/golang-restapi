package middleware

import (
	"net/http"
	"restful_api/helper"
	"restful_api/model/web"
)

type AuthMiddleware struct {
	Handler http.Handler
}

func NewAuthMiddleware(handler http.Handler) *AuthMiddleware {
	return &AuthMiddleware{Handler: handler}
}

func (middleware *AuthMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Header.Get("X-API-KEY") == "SECRET"{
		middleware.Handler.ServeHTTP(writer, request)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusUnauthorized)

	webResponse := web.WebResponse{
		Code:   http.StatusUnauthorized,
		Status: "UNAUTHORIZED",
	}

	helper.WriteToResponseBody(writer, webResponse)
}	