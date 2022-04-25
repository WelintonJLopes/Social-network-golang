package controllers

import (
	"api/src/answers"
	"api/src/authentication"
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"api/src/security"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
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

	db, err := database.Connect()
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewRepositoryUsers(db)
	userSaveDatabase, err := repository.SearchEmail(user.Email)
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}

	if err = security.PasswordCheck(userSaveDatabase.Password, user.Password); err != nil {
		answers.Err(w, http.StatusUnauthorized, err)
		return
	}

	token, err := authentication.CreateToken(userSaveDatabase.ID)
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}

	w.Write([]byte(token))
}
