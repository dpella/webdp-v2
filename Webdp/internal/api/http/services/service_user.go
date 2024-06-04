package services

import (
	"fmt"
	errors "webdp/internal/api/http"
	"webdp/internal/api/http/entity"
	"webdp/internal/api/http/repo/postgres"
	"webdp/internal/api/http/utils"
)

type UserService struct {
	postg postgres.UserPostgres
}

func NewUserService(postUser postgres.UserPostgres) UserService {
	return UserService{postg: postUser}
}

func (u UserService) GetAllUsers() ([]entity.UserResponse, error) {
	users, err := u.postg.GetUsers()
	if err != nil {
		return []entity.UserResponse{}, errors.WrapDBError(err, "get users", "all")
	}
	return users, nil
}

func (u UserService) GetUser(userHandle string) (entity.UserResponse, error) {
	user, err := u.postg.GetUser(userHandle)
	if err != nil {
		return entity.UserResponse{}, errors.WrapDBError(err, "get", userHandle)
	}
	return user, nil
}

func (u UserService) GetUserRoles(userHandle string) ([]string, error) {
	user, err := u.GetUser(userHandle)
	if err != nil {
		return nil, err
	}
	return user.Roles, nil
}

func (u UserService) CreateUser(user entity.UserPost) (string, error) {
	pwd, err := utils.HashAndSalt(user.PWD)
	if err != nil {
		return "", fmt.Errorf("%w: could not set password for %s", err, user.Handle)
	}

	user.PWD = pwd

	res, err := u.postg.CreateUser(user)

	if err != nil {
		return "", errors.WrapDBError(err, "create user", user.Handle)
	}

	return res, nil
}

func (u UserService) UpdateUser(handle string, patch entity.UserPatch) error {
	if len(patch.PWD) > 0 {
		pwd, err := utils.HashAndSalt(patch.PWD)
		if err != nil {
			return fmt.Errorf("%w: could not hash password for %s", err, handle)
		}
		patch.PWD = pwd
	}
	if _, err := u.postg.UpdateUser(handle, patch); err != nil {
		return errors.WrapDBError(err, "update", handle)
	}
	return nil
}

func (u UserService) DeleteUser(handle string) (string, error) {
	if err := u.postg.DeleteUser(handle); err != nil {
		return "", errors.WrapDBError(err, "delete", handle)
	}
	return fmt.Sprintf("user with handle %s has been deleted", handle), nil
}

func (u UserService) CompareHashedPlain(loginReq entity.LoginRequest) bool {
	res, err := u.postg.GetByHandle(loginReq.Username)
	if err != nil {
		return false
	}
	return utils.ComparePasswords(res.PWD, loginReq.PWD)
}
