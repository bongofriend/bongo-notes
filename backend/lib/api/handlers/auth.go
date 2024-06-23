package handlers

import (
	"net/http"

	"github.com/bongofriend/bongo-notes/backend/lib/api/services"
)

type authHandler struct {
	authService services.AuthService
}

func (a authHandler) Register(m *ApiMux) {
	m.ServiceResponseHandlerFunc("POST /login", a.Login)
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
//
// @Failure 500
func (a authHandler) Login(r *http.Request) ServiceResponse {
	if err := r.ParseForm(); err != nil {
		return BadRequest(err)
	}
	l, ok := extractLoginDetails(r)
	if !ok {
		return BadRequest(nil)
	}
	token, err := a.authService.GenerateToken(l.Username, l.Password)
	if err != nil {
		return BadRequest(err)
	}
	rsp := loginSuccessResponse{
		Token: token,
	}
	return Success(http.StatusAccepted, rsp)
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
