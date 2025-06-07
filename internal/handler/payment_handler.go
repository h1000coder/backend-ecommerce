package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"soulstreet/internal/model"
	"soulstreet/internal/service"
	sendjson "soulstreet/pkg/json"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var paymentData model.Payment

	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendjson.SendJsonError(w, http.StatusBadRequest, errors.New("Error to read req body"))
		return
	}

	err = json.Unmarshal(body, &paymentData)
	if err != nil {
		sendjson.SendJsonError(w, http.StatusBadRequest, errors.New("Invalid JSON"))
		return
	}

	paymentUrl, err := h.paymentService.CreatePayment(paymentData)
	if err != nil {
		sendjson.SendJsonError(w, http.StatusInternalServerError, err)
		return
	}

	sendjson.SendJson(w, http.StatusCreated, paymentUrl)

}

func (h *PaymentHandler) WebHookPayment(w http.ResponseWriter, r *http.Request) {
	var webhookData model.WebHookData

	jsonBody, err := io.ReadAll(r.Body)
	if err != nil {
		sendjson.SendJsonError(w, http.StatusUnprocessableEntity, errors.New("Não foi possivel ler o corpo da requisição"))
		return
	}

	err = json.Unmarshal(jsonBody, &webhookData)
	if err != nil {
		sendjson.SendJsonError(w, http.StatusBadRequest, errors.New("Json invalido"))
		return
	}

	if webhookData.Action == "payment.updated" {
		service.GetStatus(webhookData.Data.ID)
	} else {
		return
	}
}
