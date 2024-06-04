package entity

import (
	"fmt"
	"time"
	errors "webdp/internal/api/http"
	"webdp/internal/api/http/utils"
)

const (
	ADMIN   = "Admin"
	CURATOR = "Curator"
	ANALYST = "Analyst"
)

type User struct {
	Handle    string    `json:"handle"`
	Name      string    `json:"name"`
	Roles     []string  `json:"roles"`
	PWD       string    `json:"password"`
	CreatedOn time.Time `json:"created_time"`
	UpdatedOn time.Time `json:"updated_time"`
}

type UserResponse struct {
	Handle    string    `json:"handle"`
	Name      string    `json:"name"`
	Roles     []string  `json:"roles"`
	CreatedOn time.Time `json:"created_time"`
	UpdatedOn time.Time `json:"updated_time"`
}

type UserPost struct {
	Handle string   `json:"handle" dpvalidation:"non-empty-string"`
	Name   string   `json:"name" dpvalidation:"non-empty-string"`
	Roles  []string `json:"roles"`
	PWD    string   `json:"password" dpvalidation:"non-empty-string"`
}

type UserPatch struct {
	Name  string   `json:"name" dpvalidation:"non-empty-string"`
	Roles []string `json:"roles"`
	PWD   string   `json:"password" dpvalidation:"non-empty-string"`
}

type LoginRequest struct {
	Username string `json:"username" dpvalidation:"non-empty-string"`
	PWD      string `json:"password" dpvalidation:"non-empty-string"`
}

func (u UserPost) Valid() error {
	err := utils.ValidateNonEmptyString(u)
	if err != nil {
		return err
	}
	return validateRoles(u.Roles)
}

func (u UserPatch) Valid() error {
	err := utils.ValidateNonEmptyString(u)
	if err != nil {
		return err
	}
	return validateRoles(u.Roles)
}

func (l LoginRequest) Valid() error {
	return utils.ValidateNonEmptyString(l)
}

func validateRoles(rs []string) error {
	if len(rs) == 0 || len(rs) > 3 {
		return fmt.Errorf("%w: unexpected amount of roles: %d", errors.ErrBadInput, len(rs))
	}

	for _, role := range rs {
		if role != ADMIN && role != CURATOR && role != ANALYST {
			return fmt.Errorf("%w: unrecognized role: %s", errors.ErrBadInput, role)
		}
	}
	return nil
}
