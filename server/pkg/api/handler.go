package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"passwordless-mail-client/pkg/account"
	"strconv"

	mail "passwordless-mail-server/pkg/mail"
	"passwordless-mail-server/pkg/model"
)

type Handler struct {
	service mail.MailService
}

type MailHandler interface {
	HealthCheck(w http.ResponseWriter, r *http.Request)
	GetInbox(w http.ResponseWriter, r *http.Request)
	GetMail(w http.ResponseWriter, r *http.Request)
	SendMail(w http.ResponseWriter, r *http.Request)
}

func NewHandler(service mail.MailService) MailHandler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Printf("Request received from %s\n", r.RemoteAddr)
	fmt.Fprintln(w, "OK")
}

func (h *Handler) GetInbox(w http.ResponseWriter, r *http.Request) {
	// validate public key in header
	publicKey := r.Header.Get("x-public-key")
	if publicKey == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// validate request body
	body := model.RequestBody{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate public key
	parsePublicKey, err := account.HexToPublicKey(publicKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// read query params
	params := r.URL.Query()
	pageString := params.Get("page")
	limitString := params.Get("limit")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	serviceQuery := mail.ServiceGetInboxQuery{
		Recipient: publicKey,
		Page:      page,
		Limit:     limit,
	}

	inbox, err := h.service.GetInbox(body, parsePublicKey, serviceQuery)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(inbox)
		return
	}
	if err.Error() == "validation failed" ||
		err.Error() == "uuid is already used" ||
		err.Error() == "message timeout" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err.Error() == "bad request" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
}

func (h *Handler) GetMail(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Println("get mail called")
	fmt.Fprintln(w, "get mail in-progress")
}

func (h *Handler) SendMail(w http.ResponseWriter, r *http.Request) {
	// validate public key in header
	hexPublicKey := r.Header.Get("x-public-key")
	if hexPublicKey == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// validate public key
	publicKey, err := account.HexToPublicKey(hexPublicKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// validate request body
	body := model.RequestBody{}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := h.service.SendMail(body, publicKey)
	if err == nil {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(result)
		return
	}
	if err.Error() == "validation failed" ||
		err.Error() == "uuid is already used" ||
		err.Error() == "message timeout" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err.Error() == "invalid recipient public key" ||
		err.Error() == "bad request" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// fmt.Println(result)

	w.WriteHeader(http.StatusInternalServerError)
}
