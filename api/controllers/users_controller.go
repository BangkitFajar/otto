package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"tes/api/auth"
	"tes/api/models"
	"tes/api/responses"
	"tes/api/utils/formaterror"
)

func (server *Server) Register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)

	user.Prepare()
	err = user.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	userCreated, err := user.SaveUser(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	balance := models.Balance{
		UserId: userCreated.ID,
		Saldo:  "0",
		Status: "aktif",
	}
	err = json.Unmarshal(body, &balance)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	balanceCreated, err := balance.SaveBalance(server.DB)
	resultResponse := models.Result{
		Nama:    userCreated.Nama,
		Email:   userCreated.Email,
		Balance: *balanceCreated,
	}
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	responses.JSON(w, http.StatusCreated, resultResponse)

}

func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {

	user := models.User{}
	balance := models.Balance{}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	getUser, err := user.FindUserByID(server.DB, uint32(tokenID))
	if tokenID != uint32(getUser.ID) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	balanceResult, err := balance.GetLatestBalance(server.DB, uint32(tokenID))
	if tokenID != uint32(getUser.ID) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	responseResult := models.Result{
		ID:      getUser.ID,
		Nama:    getUser.Nama,
		Email:   getUser.Email,
		Balance: *balanceResult,
	}

	if tokenID != uint32(getUser.ID) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, responseResult)
}

func (server *Server) GetHistoryTransaction(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	balance := models.Balance{}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	result, err := user.FindUserByID(server.DB, uint32(tokenID))
	if tokenID != uint32(result.ID) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	balanceResult, err := balance.FindAllBalance(server.DB, uint32(tokenID))
	if tokenID != uint32(result.ID) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	responseResult := models.Results{
		ID:      result.ID,
		Nama:    result.Nama,
		Email:   result.Email,
		Balance: *balanceResult,
	}

	if tokenID != uint32(result.ID) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, responseResult)
}
