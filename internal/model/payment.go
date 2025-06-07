package model

type Payment struct {
	ID        string
	ProductID int     `json:"product_id"`
	Email     string  `json:"email"`
	Telephone string  `json:"telephone"`
	Address   Address `json:"address"`
	CPF       string  `json:"cpf"`
}

type Address struct {
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Number       string `json:"number"`
	State        string `json:"state"`
	Country      string `json:"country"`
}


type Item struct {
	Title       string  `json:"title"`
	Quantity    int     `json:"quantity"`
	CurrencyID  string  `json:"currency_id"`
	UnitPrice   float64 `json:"unit_price"`
}

type BackURLs struct {
	Success string `json:"success"`
	Failure string `json:"failure"`
	Pending string `json:"pending"`
}


type PreferenceRequest struct {
	Items      []Item    `json:"items"`
	BackURLs   BackURLs  `json:"back_urls"`
	AutoReturn string    `json:"auto_return"`
}


type PreferenceResponse struct {
	InitPoint string `json:"init_point"`
	ID        string `json:"id"`
}

type WebHookData struct {
	Action string `json:"action"`
	Data struct {
		ID int64 `json:"id"`
	} `json:"data"`
}