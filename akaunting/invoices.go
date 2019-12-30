package akaunting

import (
	"fmt"
	"net/url"
)

var _ UrlValuesConvertable = Invoice{}

type Invoice struct {
	InvoiceNumber   string          `json:"invoice_number" schema:"invoice_number"`
	OrderNumber     string          `json:"order_number" schema:"order_number"`
	InvoicedAt      string          `json:"invoiced_at" schema:"invoiced_at"`
	DueAt           string          `json:"due_at" schema:"due_at"`
	CurrencyCode    string          `json:"currency_code" schema:"currency_code"`
	CurrencyRate    float64         `json:"currency_rate" schema:"currency_rate"`
	Amount          float64         `json:"amount" schema:"amount"`
	CustomerID      int             `json:"customer_id" schema:"customer_id"`
	CustomerName    string          `json:"customer_name" schema:"customer_name"`
	CustomerAddress string          `json:"customer_address" schema:"customer_address"`
	CategoryID      int             `json:"category_id" schema:"category_id"`
	StatusCode      string          `json:"invoice_status_code" schema:"invoice_status_code"`
	Items           InvoiceItemData `json:"items"`
}

type InvoiceItemData struct {
	Data []InvoiceItem `json:"data"`
}

type InvoiceItem struct {
	ItemID          int     `json:"item_id" schema:"item_id"`
	Quantity        float64 `json:"quantity" schema:"quantity"`
	Price           float64 `json:"price" schema:"price"`
	Total           float64 `json:"total" schema:"total"`
	TaxID           int     `json:"tax_id" schema:"tax_id"`
	CurrencyCode    string  `json:"currency" schema:"currency"`
	Name            string  `json:"name" schema:"name"`
	LedgerAccountID int     `json:"de_account_id" schema:"de_account_id"`
}

func (i Invoice) ToUrlValues() *url.Values {
	val := &url.Values{}
	val.Set("invoice_number", i.InvoiceNumber)
	val.Set("invoiced_at", i.InvoicedAt)
	val.Set("due_at", i.DueAt)
	val.Set("currency_code", i.CurrencyCode)
	val.Set("order_number", i.OrderNumber)
	val.Set("currency_rate", fmt.Sprintf("%f", i.CurrencyRate))
	val.Set("amount", fmt.Sprintf("%f", i.Amount))
	val.Set("customer_id", fmt.Sprintf("%d", i.CustomerID))
	val.Set("customer_name", i.CustomerName)
	val.Set("customer_address", i.CustomerAddress)
	val.Set("category_id", fmt.Sprintf("%d", i.CategoryID))
	val.Set("invoice_status_code", i.StatusCode)

	for i, it := range i.Items.Data {
		val.Set(fmt.Sprintf("item[%03d][de_account_id]", i), fmt.Sprintf("%d", it.LedgerAccountID))
		val.Set(fmt.Sprintf("item[%03d][name]", i), it.Name)
		val.Set(fmt.Sprintf("item[%03d][quantity]", i), fmt.Sprintf("%f", it.Quantity))
		val.Set(fmt.Sprintf("item[%03d][price]", i), fmt.Sprintf("%f", it.Price))
		val.Set(fmt.Sprintf("item[%03d][currency]", i), it.CurrencyCode)
		val.Set(fmt.Sprintf("item[%03d][tax_id]", i), fmt.Sprintf("%d", it.TaxID))
		val.Set(fmt.Sprintf("item[%03d][total]", i), fmt.Sprintf("%f", it.Total))
	}

	return val
}

func (c *Client) CreateInvoice(invoice Invoice) error {
	req, err := c.Post("invoices", invoice)
	if err != nil {
		return err
	}

	return req.ExpectStatus(201)
}
