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
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	item := models.Account{}

	err = json.Unmarshal(body, &item)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	item.Prepare()

	err = item.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	item.OwnerID = uid;
	itemCreated, err := item.Save(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, itemCreated.ID))
	responses.JSON(w, http.StatusCreated, itemCreated)
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

	itemReceived, err := account.FindByID(server.DB, aid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, itemReceived)
}

func (server *Server) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Check if the item id is valid
	rid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the item exist
	item := models.Account{}
	err = server.DB.Debug().Model(models.Account{}).Where("id = ?", rid).Take(&item).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("account not found"))
		return
	}

	// If a user attempt to update a item not belonging to him
	if uid != item.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	itemUpdate := models.Account{}
	err = json.Unmarshal(body, &itemUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != itemUpdate.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	itemUpdate.Prepare()
	err = itemUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	itemUpdate.ID = item.ID //this is important to tell the model the item id to update, the other update field are set above

	itemUpdated, err := itemUpdate.Update(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, itemUpdated)
}

func (server *Server) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Is a valid item id given to us?
	rid, err := strconv.ParseUint(vars["id"], 10, 64)
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

	// Check if the post exist
	item := models.Account{}
	err = server.DB.Debug().Model(models.Account{}).Where("id = ?", rid).Take(&item).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this post?
	if uid != item.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = item.Delete(server.DB, rid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", rid))
	responses.JSON(w, http.StatusNoContent, "")
}