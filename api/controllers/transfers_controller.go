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
	// Verifiy authentication
	_, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Get data from body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Create transference with data from body
	transfer := models.Transfer{}
	err = json.Unmarshal(body, &transfer)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Prepare and validate transfer data
	transfer.Prepare()
	err = transfer.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Check if from_account exists
	from_account := models.Account{}
	fmt.Println("FromAccount")
	fmt.Println(from_account)
	_, err = from_account.CheckAccountExist(server.DB, uint64(transfer.FromAccountID))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// Check from_account balance
	if (from_account.Balance < transfer.Amount) {
		var err = errors.New("insufficient balance to make this transfer")
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// Check if to_account exists
	to_account := models.Account{} 
	fmt.Println("ToAccount")
	fmt.Println(to_account)
	_, err = to_account.CheckAccountExist(server.DB, uint64(transfer.ToAccountID))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// Update from_account balance Account.Balance - Transfer.Amount
	from_account.Balance = from_account.Balance - transfer.Amount
	_, err = from_account.Update(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// Update from_account balance Account.Balance - Transfer.Amount
	to_account.Balance = to_account.Balance + transfer.Amount
	_, err = to_account.Update(server.DB)
	if err != nil {
		// TODO: Devolver el dinero a from_account

		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// Create new transfer record
	transferCreated, err := transfer.Save(server.DB)
	if err != nil {
		// TODO: Devolver el dinero a from_account y descontarla de to_account
	
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