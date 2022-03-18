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

func (server *Server) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	_, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	transfer := models.Transfer{}
	err = json.Unmarshal(body, &transfer)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	transfer.Prepare()
	err = transfer.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var account models.Account

	account = models.Account{}
	_, err = account.CheckAccountExist(server.DB, uint64(transfer.FromAccountID))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	account = models.Account{}
	_, err = account.CheckAccountExist(server.DB, uint64(transfer.ToAccountID))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	transferCreated, err := transfer.Save(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, transferCreated.ID))
	responses.JSON(w, http.StatusCreated, transferCreated)
}

func (server *Server) GetAllTransfers(w http.ResponseWriter, r *http.Request) {
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	transfer := models.Transfer{}

	transfers, err := transfer.FindByOwnerId(server.DB, uid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, transfers)
}

func (server *Server) GetTransferById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	transfer := models.Transfer{}

	transferReceived, err := transfer.FindByID(server.DB, tid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if (transferReceived.FromAccountID != uid && transferReceived.ToAccountID != uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	responses.JSON(w, http.StatusOK, transferReceived)
}