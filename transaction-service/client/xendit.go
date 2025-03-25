package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type XenditClient struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

type CreateInvoiceRequest struct {
	ExternalID         string  `json:"external_id"`
	Amount             float64 `json:"amount"`
	PayerEmail         string  `json:"payer_email,omitempty"`
	Description        string  `json:"description"`
	InvoiceDuration    int     `json:"invoice_duration,omitempty"`
	CustomerName       string  `json:"customer_name,omitempty"`
	CallbackURL        string  `json:"callback_url,omitempty"`
	SuccessRedirectURL string  `json:"success_redirect_url,omitempty"`
	FailureRedirectURL string  `json:"failure_redirect_url,omitempty"`
}

type InvoiceResponse struct {
	ID                      string         `json:"id"`
	ExternalID              string         `json:"external_id"`
	UserID                  string         `json:"user_id"`
	Status                  string         `json:"status"`
	MerchantName            string         `json:"merchant_name"`
	Amount                  float64        `json:"amount"`
	PayerEmail              string         `json:"payer_email"`
	Description             string         `json:"description"`
	InvoiceURL              string         `json:"invoice_url"`
	ExpiryDate              time.Time      `json:"expiry_date"`
	AvailableBanks          []Bank         `json:"available_banks"`
	AvailableRetailOutlets  []RetailOutlet `json:"available_retail_outlets"`
	AvailableEWallets       []EWallet      `json:"available_ewallets"`
	ShouldExcludeCreditCard bool           `json:"should_exclude_credit_card"`
}

type Bank struct {
	BankCode          string  `json:"bank_code"`
	CollectionType    string  `json:"collection_type"`
	BankAccountNumber string  `json:"bank_account_number"`
	TransferAmount    float64 `json:"transfer_amount"`
}

type RetailOutlet struct {
	RetailOutletName string  `json:"retail_outlet_name"`
	PaymentCode      string  `json:"payment_code"`
	TransferAmount   float64 `json:"transfer_amount"`
}

type EWallet struct {
	EWalletType string `json:"ewallet_type"`
}

func NewXenditClient() *XenditClient {
	apiKey := os.Getenv("XENDIT_SECRET_KEY")
	if apiKey == "" {
		apiKey = "xnd_development_your_api_key" // Replace with your actual API key or use env var
	}

	return &XenditClient{
		APIKey:  apiKey,
		BaseURL: "https://api.xendit.co",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *XenditClient) CreateInvoice(req CreateInvoiceRequest) (*InvoiceResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	request, err := http.NewRequest("POST", c.BaseURL+"/invoices", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	request.SetBasicAuth(c.APIKey, "")
	request.Header.Set("Content-Type", "application/json")

	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&errorResponse); err != nil {
			return nil, fmt.Errorf("failed to decode error response: %v", err)
		}
		return nil, fmt.Errorf("xendit API error: %v", errorResponse)
	}

	var invoiceResponse InvoiceResponse
	if err := json.NewDecoder(response.Body).Decode(&invoiceResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &invoiceResponse, nil
}

func (c *XenditClient) GetInvoice(invoiceID string) (*InvoiceResponse, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/invoices/%s", c.BaseURL, invoiceID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	request.SetBasicAuth(c.APIKey, "")
	request.Header.Set("Content-Type", "application/json")

	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&errorResponse); err != nil {
			return nil, fmt.Errorf("failed to decode error response: %v", err)
		}
		return nil, fmt.Errorf("xendit API error: %v", errorResponse)
	}

	var invoiceResponse InvoiceResponse
	if err := json.NewDecoder(response.Body).Decode(&invoiceResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &invoiceResponse, nil
}
