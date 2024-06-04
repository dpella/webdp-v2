package handlers

import (
	"net/http"
	"webdp/internal/api/http/entity"
	"webdp/internal/api/http/middlewares"
	"webdp/internal/api/http/response"
	"webdp/internal/api/http/services"
	"webdp/internal/api/http/utils"

	"github.com/gorilla/mux"

	errors "webdp/internal/api/http"
)

type UserHandler struct {
	userService  services.UserService
	tokenService services.TokenService
}

type UserListNew = []entity.UserResponse

func NewUserHandler(us services.UserService, ts services.TokenService) UserHandler {
	return UserHandler{userService: us, tokenService: ts}
}

/*
Get all users. Requester needs admin or curator roles.
*/
// GetUsers godoc
// @Summary      Get all users.
// @Description  Requester needs admin or curator roles.
// @Tags         users
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  []entity.UserResponse
// @Failure      403  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v1/users [get]
// @Router       /v2/users [get]
func (h UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) error {
	if err := middlewares.ValidateRoles(r, []string{entity.ADMIN, entity.CURATOR}); err != nil {
		return RenderError(w, err)
	}

	res, err := h.userService.GetAllUsers()
	if err != nil {
		return RenderError(w, err)
	}

	if middlewares.IsRootRequestor(r) {
		return RenderResponse(w, response.NewSuccess(http.StatusOK, res))
	}

	return RenderResponse(w, response.NewSuccess(http.StatusOK, purgeRoot(res)))
}

/*
Get a user. Requester needs admin or curator roles.
*/
// GetUser godoc
// @Summary      Get a user.
// @Description  Requester needs admin or curator roles.
// @Tags         users
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      json
// @Param		 userHandle	path string true "User Handle"
// @Success      200  {object}  entity.UserResponse
// @Failure      403  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v1/user/{userHandle} [get]
// @Router       /v2/users/{userHandle} [get]
func (h UserHandler) GetUser(w http.ResponseWriter, r *http.Request) error {
	if err := middlewares.ValidateRoles(r, []string{entity.ADMIN, entity.CURATOR}); err != nil {
		if err := middlewares.ValidateSelfRequest(r); err != nil {
			return RenderError(w, err)
		}
	}

	vars := mux.Vars(r)

	if isRoot(vars["userHandle"]) {
		if !middlewares.IsRootRequestor(r) {
			return RenderError(w, errors.ErrForbidden)
		}
	}

	res, err := h.userService.GetUser(vars["userHandle"])
	if err != nil {
		return RenderError(w, err)
	}

	return RenderResponse(w, response.NewSuccess(http.StatusOK, res))
}

/*
Create new user. Requester needs admin role.
*/
// PostUser godoc
// @Summary      Create new user.
// @Description  Requester needs admin role.
// @Tags         users
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      json
// @Param		 userRequest	body entity.UserPost true "User Request"
// @Success      201
// @Failure		 400  {object}  response.Error
// @Failure      403  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v1/users [post]
// @Router       /v2/users [post]
func (h UserHandler) PostUsers(w http.ResponseWriter, r *http.Request) error {
	if err := middlewares.ValidateRoles(r, []string{entity.ADMIN}); err != nil {
		return RenderError(w, err)
	}

	var createUser entity.UserPost
	if err := utils.ParseJsonRequestBody[entity.UserPost](r, &createUser); err != nil {
		return RenderError(w, err)
	}
	if err := createUser.Valid(); err != nil {
		return RenderError(w, err)
	}

	if _, err := h.userService.CreateUser(createUser); err != nil {
		return RenderError(w, err)
	}

	return RenderResponse(w, response.EmptyResponse(http.StatusCreated))
}

/*
Update a user. Requester needs admin role.
*/
// PatchUser godoc
// @Summary      Update a user.
// @Description  Update name, password and roles of a user.
// @Description  Requester needs admin role.
// @Tags         users
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      json
// @Param		 userHandle	path string true "User Handle"
// @Success      204
// @Failure		 400  {object}  response.Error
// @Failure      403  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v1/user/{userHandle} [patch]
// @Router       /v2/users/{userHandle} [patch]
func (h UserHandler) PatchUser(w http.ResponseWriter, r *http.Request) error {
	if err := middlewares.ValidateRoles(r, []string{entity.ADMIN}); err != nil {
		return RenderError(w, err)
	}

	vars := mux.Vars(r)
	var patch entity.UserPatch
	if err := utils.ParseJsonRequestBody[entity.UserPatch](r, &patch); err != nil {
		return RenderError(w, err)
	}

	// kick user whose data is updated (invalidate session)
	// optimally, user would get a notification about this.
	if patch.PWD != "" || patch.Roles != nil {
		if err := h.tokenService.LogOffUser(vars["userHandle"]); err != nil {
			return RenderError(w, err)
		}
	}

	// Update user data
	if err := h.userService.UpdateUser(vars["userHandle"], patch); err != nil {
		return RenderError(w, err)
	}

	return RenderResponse(w, response.NoContent())
}

/*
Delete a user. Requester needs admin role.
*/
// DeleteUser godoc
// @Summary      Delete a user.
// @Description  Requester needs admin role.
// @Tags         users
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      json
// @Param		 userHandle	path string true "User Handle"
// @Success      204
// @Failure		 400  {object}  response.Error
// @Failure      403  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v1/user/{userHandle} [delete]
// @Router       /v2/users/{userHandle} [delete]
func (h UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	if err := middlewares.ValidateRoles(r, []string{entity.ADMIN}); err != nil {
		return RenderError(w, err)
	}
	vars := mux.Vars(r)
	if vars["userHandle"] == "root" {
		return RenderError(w, errors.ErrForbidden)
	} else {
		if err := h.tokenService.LogOffUser(vars["userHandle"]); err != nil {
			return RenderError(w, err)
		}
		if _, err := h.userService.DeleteUser(vars["userHandle"]); err != nil {
			return RenderError(w, err)
		}

		return RenderResponse(w, response.NoContent())
	}
}

func isRoot(user string) bool {
	return user == "root"
}

func purgeRoot(users []entity.UserResponse) []entity.UserResponse {
	res := make([]entity.UserResponse, 0)
	for _, user := range users {
		if !isRoot(user.Handle) {
			res = append(res, user)
		}
	}
	return res
}
