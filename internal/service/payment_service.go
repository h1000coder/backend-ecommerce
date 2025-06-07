package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"soulstreet/internal/model"
	"soulstreet/internal/repository"
)

type PaymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (p *PaymentService) CreatePayment(paymentData model.Payment) (string, error) {
	url := "https://api.mercadopago.com/checkout/preferences"
	
	urlSucess := os.Getenv("URL_SUCESS")
	urlFailure := os.Getenv("URL_FAILURE")
	urlPending := os.Getenv("URL_PENDING")
	mpToken := os.Getenv("MP_TOKEN")
	webhookUrl := os.Getenv("URL_WEBHOOK")
	
	if urlSucess == "" || urlFailure == "" || urlPending == "" || webhookUrl == "" {
		log.Fatal("Some dotenv var is empty")
	}
	
	productInfo, err := p.repo.GetProductInfo(paymentData.ProductID)
	if err != nil{
		return "", errors.New("Erro ao buscar informações do produto")
	}

	payload := map[string]interface{} {
		"items": []map[string]interface{}{
			{
				"title": productInfo.Name,
				"quantity": 1,
				"currency_id": "BRL",
				"unit_price": productInfo.Price,
			},
		},
		"payer": map[string]interface{}{
            "email": paymentData.Email,
        },
        "back_urls": map[string]string{
            "success": urlSucess,
            "failure": urlFailure,
            "pending": urlPending,
        },
        "auto_return": "approved",
		"notification_url": webhookUrl,
	}
	
	body, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+mpToken)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil{
		panic(err)
	}
	
	defer resp.Body.Close()
	
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	
	if resp.StatusCode != 201 {
		fmt.Println("Erro ao criar preferencia: ", string(respBody))
		return "", err
	}
	var responseData map[string]interface{}
    err = json.Unmarshal(respBody, &responseData)
    if err != nil {
        panic(err)
    }
    fmt.Println(responseData)
    if initPoint, ok := responseData["init_point"].(string); ok {
        fmt.Println("Link para Checkout Pro:", initPoint)
        return initPoint, nil
    } else {
        fmt.Println("Resposta inesperada:", string(respBody))
        return "", errors.New("Unexpected response")
    }
}

func GetStatus(id int64) error {
	statusURL := os.Getenv("STATUS_URL")
	mpToken := os.Getenv("MP_TOKEN")
	if statusURL == "" || mpToken == "" {
		log.Fatal("Some dotenv var is empty")
	}
	

	statusURL = fmt.Sprintf("%s%d", statusURL, id)

	req, err := http.NewRequest("GET", statusURL, nil)
	if err != nil {
		return errors.New("Erro ao criar requisição de status")
	}


	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+mpToken)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return errors.New("Erro ao fazer requisição de status")
	}

	defer resp.Body.Close()

	var responseData map[string]interface{}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Erro ao ler resposta da requisição")
	}

	if err = json.Unmarshal(respBody, &responseData); err != nil {
		return errors.New("Erro ao dar unmarshal no respBody")
	}

	if responseData["status"] == "approved" {
		payer, ok := responseData["payer"].(map[string]interface{})
		if !ok {
			return errors.New("Campo payer não encontrado")
		}

		email, ok := payer["email"].(string)
		if !ok {
			return errors.New("Campo email não encontrado")
		}
		fmt.Println(email)
	} else {
		fmt.Printf("O status do pagamento foi atualizado para: %s\n", responseData["status"])
		return nil
	}
	return nil
}
