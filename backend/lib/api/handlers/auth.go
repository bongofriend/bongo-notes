package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bongofriend/bongo-notes/backend/lib/api/models"
	"github.com/bongofriend/bongo-notes/backend/lib/api/services"
	httputils "github.com/bongofriend/bongo-notes/backend/lib/api/utils"
)

type authHandler struct {
	authService services.AuthService
}

func (a authHandler) Register(m *ApiMux) {
	m.HandleFunc("POST /login", a.Login)
	m.AuthenticatedHandlerFunc("/greet", a.Greet)
}

func NewAuthHandler(srvContainer services.ServicesContainer) ApiHandler {
	return authHandler{
		authService: srvContainer.AuthService(),
	}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginSuccessResponse struct {
	Token string `json:"token"`
}

// Login godoc
//
//	@Summary	Login with user credentials
//	@Tags		auth
//	@Router		/login [post]
//	@Param		logindata	formData	handlers.loginRequest	true	"User login"
//
//	@Success	200			{object}	handlers.loginSuccessResponse
//
//	@Failure	400
func (a authHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		httputils.BadRequestError(w)
		return
	}
	l, ok := extractLoginDetails(r)
	if !ok {
		httputils.BadRequestError(w)
		return
	}
	token, err := a.authService.GenerateToken(l.Username, l.Password)
	if err != nil {
		log.Println(err)
		httputils.InternalServerError(w)
		return
	}
	rsp := loginSuccessResponse{
		Token: token,
	}
	httputils.WriteJson(w, http.StatusAccepted, rsp)
}

func extractLoginDetails(r *http.Request) (loginRequest, bool) {
	username := r.FormValue("username")
	if len(username) == 0 {
		return loginRequest{}, false
	}
	password := r.FormValue("password")
	if len(password) == 0 {
		return loginRequest{}, false
	}
	return loginRequest{
		Username: username,
		Password: password,
	}, true
}

// Greet godoc
//
//	@Summary	Greet the user
//	@Tags		auth
//	@Router		/greet [get]
//	@Security	BearerAuth
func (a authHandler) Greet(user models.User, w http.ResponseWriter, r *http.Request) {
	rsp := map[string]string{
		"msg": fmt.Sprintf("Hello, %s", user.Username),
	}
	httputils.WriteJson(w, http.StatusAccepted, rsp)
}
