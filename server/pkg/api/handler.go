package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"passwordless-mail-client/pkg/account"
	"strconv"
	"time"

	mail "passwordless-mail-server/pkg/mail"
	"passwordless-mail-server/pkg/model"

	"github.com/google/uuid"
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
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Printf("Request received from %s\n", r.RemoteAddr)
	fmt.Fprintln(w, "OK")
}

func (h *Handler) GetInbox(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
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
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userAddress := r.Header.Get("x-public-key")
	if userAddress == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// validate public key
	publicKey, err := account.HexToPublicKey(userAddress)
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

	// validate message
	var message map[string]interface{}
	err = json.Unmarshal([]byte(body.Data), &message)
	if err != nil {
		fmt.Println("err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = uuid.Parse(message["email_id"].(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = uuid.Parse(message["id"].(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = time.Parse(time.RFC3339, message["timestamp"].(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := h.service.GetMail(body, publicKey, userAddress)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
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
	if err.Error() == "mail not found" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
}

func (h *Handler) SendMail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
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
