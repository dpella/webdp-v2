package handlers

import (
	"net/http"
	errors "webdp/internal/api/http"
	"webdp/internal/api/http/entity"
	"webdp/internal/api/http/middlewares"
	"webdp/internal/api/http/response"
	"webdp/internal/api/http/services"
	"webdp/internal/api/http/utils"
)

type LoginHandler struct {
	userService  services.UserService
	tokenService services.TokenService
}

type JWTAccessToken = string
type TokenResponse = map[string]JWTAccessToken

func NewLoginHandler(userService services.UserService, tokenService services.TokenService) LoginHandler {
	return LoginHandler{userService: userService, tokenService: tokenService}
}

// assume that input is validated
// we could refactor some parsing for nices error handling
// Login godoc
// @Summary      Login User
// @Description  Login user with user/password credentials.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param		 loginRequest 	body	entity.LoginRequest	true  "Login Request"
// @Success      200  {object}  response.Token
// @Failure      400  {object}  response.Error
// @Failure      401  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v1/login [post]
// @Router       /v2/login [post]
func (lh LoginHandler) LoginRequestHandler(w http.ResponseWriter, r *http.Request) error {
	var loginRequest entity.LoginRequest
	// parse request
	if err := utils.ParseJsonRequestBody[entity.LoginRequest](r, &loginRequest); err != nil {
		return RenderError(w, err)
	}
	// validate
	if err := loginRequest.Valid(); err != nil {
		return RenderError(w, err)
	}

	// validate credentials
	if !lh.userService.CompareHashedPlain(loginRequest) {
		return RenderError(w, errors.ErrUnauthorized)
	}

	// retrieve user to make token
	user, err := lh.userService.GetUser(loginRequest.Username)
	if err != nil {
		return RenderError(w, err)
	}

	// make token
	token, err := lh.tokenService.IssueNewTokenFor(user.Handle, user.Roles)
	if err != nil {
		return RenderError(w, err)
	}
	tokenResponse := make(map[string]JWTAccessToken)
	tokenResponse["jwt"] = token

	// all ok - return token
	return RenderResponse(w, response.NewSuccess(http.StatusOK, tokenResponse))
}

// logout godoc
// @Summary      Logout User
// @Description  Logout user from session.
// @Tags         auth
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      json
// @Success      204
// @Failure      400  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v1/logout [post]
// @Router       /v2/logout [post]
func (lh LoginHandler) LogoutRequestHandler(w http.ResponseWriter, r *http.Request) error {
	var userToken services.JWTTokenClaims
	if _, err := middlewares.ExtracAuthnHeader[*services.JWTTokenClaims](r.Header, &userToken); err != nil {
		return RenderError(w, err)
	}
	if err := lh.tokenService.LogOffUser(userToken.Handle); err != nil {
		return RenderError(w, err)
	}
	return RenderResponse(w, response.NoContent())
}
