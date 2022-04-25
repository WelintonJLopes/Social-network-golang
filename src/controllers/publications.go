package controllers

import (
	"api/src/answers"
	"api/src/authentication"
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreatePublication(w http.ResponseWriter, r *http.Request) {
	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		answers.Err(w, http.StatusUnauthorized, err)
		return
	}

	bodyRequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		answers.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var publication models.Publication
	if err := json.Unmarshal(bodyRequest, &publication); err != nil {
		answers.Err(w, http.StatusBadRequest, err)
		return
	}

	publication.AuthorID = userID

	if err = publication.Prepare(); err != nil {
		answers.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewRepositoryPublications(db)
	publication.ID, err = repository.Create(publication)
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}

	answers.JSON(w, http.StatusCreated, publication)
}

func SearchPublications(w http.ResponseWriter, r *http.Request) {
	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		answers.Err(w, http.StatusUnauthorized, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewRepositoryPublications(db)
	publications, err := repository.Search(userID)
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}

	answers.JSON(w, http.StatusOK, publications)
}

func SearchPublication(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	publicationID, err := strconv.ParseUint(parameters["publicationId"], 10, 64)
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

	repository := repositories.NewRepositoryPublications(db)
	publication, err := repository.SearchID(publicationID)
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}

	answers.JSON(w, http.StatusOK, publication)
}

func UpdatePublication(w http.ResponseWriter, r *http.Request) {
	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		answers.Err(w, http.StatusUnauthorized, err)
		return
	}

	parameters := mux.Vars(r)
	publicationID, err := strconv.ParseUint(parameters["publicationId"], 10, 64)
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

	repository := repositories.NewRepositoryPublications(db)
	publicationSavedDatabase, err := repository.SearchID(publicationID)
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}

	if userID != publicationSavedDatabase.AuthorID {
		answers.Err(w, http.StatusForbidden, errors.New("you cannot change a post that is not yours"))
		return
	}

	bodyRequest, err := ioutil.ReadAll(r.Body)
	if err != nil {
		answers.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var publication models.Publication
	if err := json.Unmarshal(bodyRequest, &publication); err != nil {
		answers.Err(w, http.StatusBadRequest, err)
		return
	}

	if err = publication.Prepare(); err != nil {
		answers.Err(w, http.StatusBadRequest, err)
		return
	}

	if err = repository.Update(publicationID, publication); err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}

	answers.JSON(w, http.StatusNoContent, nil)
}

func DeletePublication(w http.ResponseWriter, r *http.Request) {
	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		answers.Err(w, http.StatusUnauthorized, err)
		return
	}

	parameters := mux.Vars(r)
	publicationID, err := strconv.ParseUint(parameters["publicationId"], 10, 64)
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

	repository := repositories.NewRepositoryPublications(db)
	publicationSavedDatabase, err := repository.SearchID(publicationID)
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}

	if userID != publicationSavedDatabase.AuthorID {
		answers.Err(w, http.StatusForbidden, errors.New("you cannot delete a post that is not yours"))
		return
	}

	if err = repository.Delete(publicationID); err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}

	answers.JSON(w, http.StatusNoContent, nil)
}

func SearchPublicationsUser(w http.ResponseWriter, r *http.Request) {
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

	repository := repositories.NewRepositoryPublications(db)
	publications, err := repository.SearchUser(userID)
	if err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}

	answers.JSON(w, http.StatusOK, publications)
}

func LikePublication(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	publicationID, err := strconv.ParseUint(parameters["publicationId"], 10, 64)
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

	repository := repositories.NewRepositoryPublications(db)
	if err = repository.LikePublication(publicationID); err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}

	answers.JSON(w, http.StatusNoContent, nil)
}

func DislikePublication(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	publicationID, err := strconv.ParseUint(parameters["publicationId"], 10, 64)
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

	repository := repositories.NewRepositoryPublications(db)
	if err = repository.DislikePublication(publicationID); err != nil {
		answers.Err(w, http.StatusInternalServerError, err)
		return
	}

	answers.JSON(w, http.StatusNoContent, nil)
}
