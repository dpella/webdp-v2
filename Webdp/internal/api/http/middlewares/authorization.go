package middlewares

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"webdp/internal/api/http/entity"
	"webdp/internal/api/http/services"
	"webdp/internal/api/http/utils"

	errors "webdp/internal/api/http"

	"github.com/gorilla/mux"
)

/*
Checks whether the requester has at least one of the provided roles.
*/
func ValidateRoles(r *http.Request, roles []string) error {
	var userToken services.JWTTokenClaims
	if _, err := ExtracAuthnHeader(r.Header, &userToken); err != nil {
		return err
	}
	for _, role := range roles {
		if slices.Contains(userToken.Roles, role) {
			return nil
		}
	}
	return fmt.Errorf("%w: missing role", errors.ErrForbidden)
}

/*
Checks whether the requester is the requestee.
*/
func ValidateSelfRequest(r *http.Request) error {
	user, err := extractUserHandle(&r.Header)
	if err != nil {
		return err
	}
	if user != mux.Vars(r)["userHandle"] {
		return fmt.Errorf("%w: requester is not authorized", errors.ErrForbidden)
	}
	return nil
}

/*
Checks whether the requester is the owner of the requested dataset.
*/
func ValidateOwnership(r *http.Request, ds *services.DatasetService) error {
	user, err := extractUserHandle(&r.Header)
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(mux.Vars(r)["datasetId"], 10, 64)
	if err != nil {
		return err
	}
	owner, err := ds.GetDatasetOwner(id)
	if err != nil || user != owner {
		return fmt.Errorf("%w: user is not owner", errors.ErrForbidden)
	}
	return nil
}

/*
Checks whether owner of a new dataset has a certain role.
*/
func ValidateNewOwnerRole(r *http.Request, us *services.UserService, role string) error {
	var d entity.DatasetCreate
	if err := utils.ParseJsonRequestBody(r, &d); err != nil {
		return err
	}
	roles, err := us.GetUserRoles(d.Owner)
	if err != nil || !slices.Contains(roles, role) {
		return fmt.Errorf("%w: missing role", errors.ErrForbidden)
	}
	return nil
}

/*
Checks whether the requester has granted access to the dataset.
A user has granted access to a dataset to which it has a budget (allocated or consumed).
*/
func ValidateGrantedAccess(r *http.Request, bs *services.BudgetService) error {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["datasetId"], 10, 64)
	if err != nil {
		return err
	}

	user, err := extractUserHandle(&r.Header)
	if err != nil {
		return err
	}

	return hasAccess(bs, id, user)
}

/*
Filters a list of datasets: Removes those where the requester does not have granted access.
A user has granted access to a dataset to which it has a budget (allocated or consumed).
*/
func FilterGrantedAccess(r *http.Request, bs *services.BudgetService, datasets *[]entity.DatasetInfo) error {
	if datasets == nil {
		return errors.ErrNotFound
	}

	user, err := extractUserHandle(&r.Header)
	if err != nil {
		return err
	}
	var accessed []entity.DatasetInfo
	for _, set := range *datasets {
		if err := hasAccess(bs, set.Id, user); err == nil {
			accessed = append(accessed, set)
		}
	}
	if len(accessed) == 0 {
		return fmt.Errorf("%w: no access to datasets", errors.ErrForbidden)
	}
	*datasets = accessed
	return nil
}

func IsRootRequestor(r *http.Request) bool {
	user, err := extractUserHandle(&r.Header)
	if err != nil {
		return false
	}

	if user == "root" {
		return true
	}

	return false

}

func extractUserHandle(rh *http.Header) (string, error) {
	var userToken services.JWTTokenClaims
	if _, err := ExtracAuthnHeader(*rh, &userToken); err != nil {
		return "", err
	}
	return userToken.Handle, nil
}

func hasAccess(bs *services.BudgetService, id int64, user string) error {
	hasAccess, err := bs.UserHasDatasetBudget(user, id)
	if err != nil || !hasAccess {
		return fmt.Errorf("%w: user has no access to dataset", errors.ErrForbidden)
	}
	return nil
}
