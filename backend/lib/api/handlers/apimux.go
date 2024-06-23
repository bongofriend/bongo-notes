package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bongofriend/bongo-notes/backend/lib/api/models"
	"github.com/bongofriend/bongo-notes/backend/lib/api/services"
	httputils "github.com/bongofriend/bongo-notes/backend/lib/api/utils"
	"github.com/bongofriend/bongo-notes/backend/lib/config"
)

type ApiMux struct {
	http.ServeMux
	config      config.Config
	authService services.AuthService
}

func NewApiMux(c config.Config, a services.AuthService) *ApiMux {
	return &ApiMux{
		config:      c,
		authService: a,
	}
}

type AuthenticatedHttpHandlerFunc func(user models.User, w http.ResponseWriter, r *http.Request)

func (a *ApiMux) AuthenticatedHandlerFunc(pattern string, h AuthenticatedHttpHandlerFunc) {
	a.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		user, ok := a.authService.Authenticate(r)
		if !ok {
			httputils.NotAuthenticatedError(w)
			return
		}
		h(user, w, r)
	})
}

type ServiceResponseHandlerFunc func(r *http.Request) ServiceResponse

func (a *ApiMux) ServiceResponseHandlerFunc(pattern string, handler ServiceResponseHandlerFunc) {
	a.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		serviceError := handler(r)
		serviceError.WriteResponse(w)
	})
}

type AuthenticatedServiceHandlerFunc func(user models.User, r *http.Request) ServiceResponse

func (a *ApiMux) AuthenticatedServiceResponseHandlerFunc(pattern string, h AuthenticatedServiceHandlerFunc) {
	a.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		user, ok := a.authService.Authenticate(r)
		if !ok {
			Unauthorized(nil).WriteResponse(w)
		} else {
			serviceResponse := h(user, r)
			serviceResponse.WriteResponse(w)
		}
	})
}

type ApiHandler interface {
	Register(*ApiMux)
}

type ServiceResponse interface {
	WriteResponse(w http.ResponseWriter)
}

func ServiceErrorWithMessage(statusCode int, innerError error, message string) ServiceResponse {
	return serviceErrorMessageResponse{
		serviceBaseError: serviceBaseError{
			StatusCode: statusCode,
			InnerError: innerError,
		},
		message: message,
	}
}

func ServiceErrorWithBody[T any](statusCode int, innerError error, body T) ServiceResponse {
	return serviceErrorBodyResponse[T]{
		serviceBaseError: serviceBaseError{
			StatusCode: statusCode,
			InnerError: innerError,
		},
		Body: body,
	}
}

type serviceBaseError struct {
	StatusCode int
	InnerError error
}

type serviceErrorBodyResponse[T any] struct {
	serviceBaseError
	Body T
}

func (s serviceErrorBodyResponse[T]) WriteResponse(w http.ResponseWriter) {
	if s.InnerError != nil {
		log.Println(s.InnerError)
	}
	data, err := json.Marshal(s.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Sever Error"))
		return
	}
	w.WriteHeader(s.StatusCode)
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

type serviceErrorMessageResponse struct {
	serviceBaseError
	message string
}

func (s serviceErrorMessageResponse) WriteResponse(w http.ResponseWriter) {
	log.Println(s.InnerError)
	w.WriteHeader(s.StatusCode)
	w.Write([]byte(s.message))
}

type successMessageResponse struct {
	statusCode int
	message    string
}

func (s successMessageResponse) WriteResponse(w http.ResponseWriter) {
	w.WriteHeader(s.statusCode)
	w.Write([]byte(s.message))
}

func ServiceMessageSuccessResponse(stattusCode int, message string) ServiceResponse {
	return &successMessageResponse{
		statusCode: stattusCode,
		message:    message,
	}
}

type successBodyResponse[T any] struct {
	statusCode int
	body       T
}

func (s successBodyResponse[T]) WriteResponse(w http.ResponseWriter) {
	data, err := json.Marshal(s.body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	} else {
		w.WriteHeader(s.statusCode)
		w.Header().Add("Content-Type", "application/json")
		w.Write(data)
	}
}

func ServiceSuccessBodyResponse[T any](statusCode int, body T) ServiceResponse {
	return successBodyResponse[T]{
		statusCode: statusCode,
		body:       body,
	}
}

func Unauthorized(err error) ServiceResponse {
	return ServiceErrorWithMessage(http.StatusUnauthorized, err, "Not Authenticated")
}

func NotFound(err error) ServiceResponse {
	return ServiceErrorWithMessage(http.StatusNotFound, err, "Not Found")
}

func InternalServerError(err error) ServiceResponse {
	return ServiceErrorWithMessage(http.StatusInternalServerError, err, "Internal Server Error")
}

func BadRequest(err error) ServiceResponse {
	return ServiceErrorWithMessage(http.StatusBadRequest, err, "Bad Request")
}

func Success[T any](statusCode int, data T) ServiceResponse {
	return ServiceSuccessBodyResponse(statusCode, data)
}

func Ok() ServiceResponse {
	return ServiceMessageSuccessResponse(http.StatusOK, "OK")
}

func Accepted() ServiceResponse {
	return ServiceMessageSuccessResponse(http.StatusAccepted, "Accepted")
}
