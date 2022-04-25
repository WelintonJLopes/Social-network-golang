package controllers

import (
	"api/src/answers"
	"api/src/authentication"
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"api/src/security"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		answers.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err = json.Unmarshal(bodyRequest, &user); err != nil {
		answers.Err(w, http.StatusBadRequest, err)
		return
	}

	if err = user.Prepare("register"); err != nil {
		answers.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewRepositoryUsers(db)
	user.ID, err = repository.Create(user)
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}

	answers.JSON(w, http.StatusCreated, user)
}

func SearchUsers(w http.ResponseWriter, r *http.Request) {
	nameOuNick := strings.ToLower(r.URL.Query().Get("user"))

	db, err := database.Connect()
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewRepositoryUsers(db)
	users, err := repository.Search(nameOuNick)

	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}

	answers.JSON(w, http.StatusOK, users)
}

func SearchUser(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)

	userID, err := strconv.ParseUint(parameters["userId"], 10, 64)
	if err != nil {
		answers.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewRepositoryUsers(db)
	user, err := repository.SearchID(userID)
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}

	answers.JSON(w, http.StatusOK, user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	userID, err := strconv.ParseUint(parameters["userId"], 10, 64)
	if err != nil {
		answers.Err(w, http.StatusBadRequest, err)
		return
	}

	userIDToken, err := authentication.ExtractUserID(r)
	if err != nil {
		answers.Err(w, http.StatusUnauthorized, err)
		return
	}

	if userID != userIDToken {
		answers.Err(w, http.StatusForbidden, errors.New("you cannot change a user other than yours"))
		return
	}

	bodyRequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		answers.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err = json.Unmarshal(bodyRequest, &user); err != nil {
		answers.Err(w, http.StatusBadRequest, err)
		return
	}

	if err = user.Prepare("edit"); err != nil {
		answers.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewRepositoryUsers(db)
	if err = repository.Update(userID, user); err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}

	answers.JSON(w, http.StatusNoContent, nil)

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	userID, err := strconv.ParseUint(parameters["userId"], 10, 64)
	if err != nil {
		answers.Err(w, http.StatusBadRequest, err)
		return
	}

	userIDToken, err := authentication.ExtractUserID(r)
	if err != nil {
		answers.Err(w, http.StatusUnauthorized, err)
		return
	}

	if userID != userIDToken {
		answers.Err(w, http.StatusForbidden, errors.New("you cannot delete a user other than yours"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewRepositoryUsers(db)
	if err = repository.Delete(userID); err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}

	answers.JSON(w, http.StatusNoContent, nil)
}

func FollowUser(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	userID, err := strconv.ParseUint(parameters["userId"], 10, 64)
	if err != nil {
		answers.Err(w, http.StatusBadRequest, err)
		return
	}

	followerID, err := authentication.ExtractUserID(r)
	if err != nil {
		answers.Err(w, http.StatusUnauthorized, err)
		return
	}

	if userID == followerID {
		answers.Err(w, http.StatusForbidden, errors.New("you cannot follow your own username"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewRepositoryUsers(db)
	if err = repository.Follow(userID, followerID); err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}

	answers.JSON(w, http.StatusNoContent, nil)
}

func UnfollowollowUser(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	userID, err := strconv.ParseUint(parameters["userId"], 10, 64)
	if err != nil {
		answers.Err(w, http.StatusBadRequest, err)
		return
	}

	followerID, err := authentication.ExtractUserID(r)
	if err != nil {
		answers.Err(w, http.StatusUnauthorized, err)
		return
	}

	if userID == followerID {
		answers.Err(w, http.StatusForbidden, errors.New("you can't stop following your own username"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewRepositoryUsers(db)
	if err = repository.Unfollowollow(userID, followerID); err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}

	answers.JSON(w, http.StatusNoContent, nil)
}

func SearchFollowers(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	userID, err := strconv.ParseUint(parameters["userId"], 10, 64)
	if err != nil {
		answers.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewRepositoryUsers(db)
	followers, err := repository.SearchFollowers(userID)
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}

	answers.JSON(w, http.StatusOK, followers)
}

func SearchFollowing(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	userID, err := strconv.ParseUint(parameters["userId"], 10, 64)
	if err != nil {
		answers.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewRepositoryUsers(db)
	following, err := repository.SearchFollowing(userID)
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}

	answers.JSON(w, http.StatusOK, following)
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	userID, err := strconv.ParseUint(parameters["userId"], 10, 64)
	if err != nil {
		answers.Err(w, http.StatusBadRequest, err)
		return
	}

	userIDToken, err := authentication.ExtractUserID(r)
	if err != nil {
		answers.Err(w, http.StatusUnauthorized, err)
		return
	}

	if userID != userIDToken {
		answers.Err(w, http.StatusForbidden, errors.New("you are not allowed to update the password of a user other than yours"))
		return
	}

	bodyRequest, err := ioutil.ReadAll(r.Body)
	var password models.Password
	if err := json.Unmarshal(bodyRequest, &password); err != nil {
		answers.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewRepositoryUsers(db)
	passwordSavedDatabase, err := repository.SearchPassword(userID)
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}
	if err = security.PasswordCheck(passwordSavedDatabase, password.Current); err != nil {
		answers.Err(w, http.StatusUnauthorized, errors.New("the new password does not match the current password"))
		return
	}
	passwordWithHash, err := security.HashPassword(password.New)
	if err != nil {
		answers.Err(w, http.StatusBadRequest, err)
		return
	}

	if err = repository.UpdatePassword(userID, string(passwordWithHash)); err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}

	answers.JSON(w, http.StatusNoContent, nil)
}
