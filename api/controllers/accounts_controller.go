package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/alejandrehl/simple-bank-api/api/auth"
	"github.com/alejandrehl/simple-bank-api/api/models"
	"github.com/alejandrehl/simple-bank-api/api/responses"
	"github.com/alejandrehl/simple-bank-api/api/utils/formaterror"
	"github.com/gorilla/mux"
)

func (server *Server) CreateAccount(w http.ResponseWriter, r *http.Request) {
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	account := models.Account{}
	err = json.Unmarshal(body, &account)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	account.Prepare()
	err = account.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	account.OwnerID = uid;
	accountCreated, err := account.Save(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, accountCreated.ID))
	responses.JSON(w, http.StatusCreated, accountCreated)
}

func (server *Server) GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	account := models.Account{}
	accounts, err := account.FindAll(server.DB, uid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, accounts)
}

func (server *Server) GetAccountById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	aid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	account := models.Account{}
	accountReceived, err := account.FindByID(server.DB, aid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if (accountReceived.OwnerID != uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	responses.JSON(w, http.StatusOK, accountReceived)
}

func (server *Server) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	aid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	account := models.Account{}
	err = server.DB.Debug().Model(models.Account{}).Where("id = ?", aid).Take(&account).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("account not found"))
		return
	}

	if uid != account.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	accountUpdate := models.Account{}
	err = json.Unmarshal(body, &accountUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if uid != accountUpdate.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	err = accountUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	accountUpdate.ID = account.ID //this is important to tell the model the item id to update, the other update field are set above
	accountUpdated, err := accountUpdate.Update(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	responses.JSON(w, http.StatusOK, accountUpdated)
}

func (server *Server) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	aid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the account exist
	account := models.Account{}
	err = server.DB.Debug().Model(models.Account{}).Where("id = ?", aid).Take(&account).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this account?
	if uid != account.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the account balance is more than zero
	if account.Balance > 0 {
		err = errors.New("you cannot delete an account with a balance greater than $0")
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	_, err = account.Delete(server.DB, aid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", aid))
	responses.JSON(w, http.StatusNoContent, "")
}