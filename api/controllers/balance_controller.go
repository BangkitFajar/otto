package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"tes/api/auth"
	"tes/api/models"
	"tes/api/responses"
)

func (server *Server) Checking(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	user := models.User{}
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

	balance := models.Balance{}
	err = json.Unmarshal(body, &balance)

	balanceResult, err := balance.GetLatestBalance(server.DB, uint32(tokenID))
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	paymentMethodVA := models.Paymentmethod{}
	err = json.Unmarshal(body, &paymentMethodVA)
	checkingVA, err := paymentMethodVA.FindBalanceByVA(server.DB, paymentMethodVA.VA)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	saldoExisting, err := strconv.Atoi(balanceResult.Saldo)
	topUpSaldo, err := strconv.Atoi(checkingVA.Nominal)
	calculate := saldoExisting + topUpSaldo
	tes := models.Balance{
		UserId: getUser.ID,
		Saldo:  strconv.Itoa(calculate),
		Status: "topup",
	}
	topUp, err := tes.SaveBalance(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	update := models.Paymentmethod{
		Status:    "success",
		UpdatedAt: time.Now(),
	}
	updatePayment, err := update.UpdatePayment(server.DB, checkingVA.ID)
	log.Println(updatePayment)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, topUp)
}

func (server *Server) GenerateVA(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	user := models.User{}
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

	reqParam := models.Paymentmethod{}
	err = json.Unmarshal(body, &reqParam)

	uniqCode := fmt.Sprint(time.Now().Nanosecond())[:3]
	numberVA := fmt.Sprint(time.Now().Year(), "", int(time.Now().Month()), uniqCode)

	paymentMethod := models.Paymentmethod{
		Nama:      reqParam.Nama,
		UserId:    tokenID,
		Status:    "pending",
		VA:        numberVA,
		Nominal:   reqParam.Nominal,
		CreatedAt: time.Now(),
	}
	err = json.Unmarshal(body, &paymentMethod)

	paymentMethodSave, err := paymentMethod.SavePaymentMethod(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, paymentMethodSave)
}
