package akaunting

import "fmt"

type GetCustomersResponse struct {
	Data []Customer `json:"data"`
}

type GetCustomerResponse struct {
	Data Customer `json:"data"`
}

type Customer struct {
	ID           int    `json:"id" schema:"id"`
	Name         string `json:"name" schema:"name"`
	Address      string `json:"address" schema:"address"`
	Phone        string `json:"phone" schema:"phone"`
	TaxNumber    string `json:"tax_number" schema:"tax_number"`
	Email        string `json:"email" schema:"email"`
	Website      string `json:"website" schema:"website"`
	CurrencyCode string `json:"currency_code" schema:"currency_code"`
	Enabled      int    `json:"enabled" schema:"enabled"`
}

func (c *Client) GetCustomers() ([]Customer, error) {
	customersResponse := GetCustomersResponse{}
	customers := make([]Customer, 0)
	req, err := c.Get("customers")
	if err != nil {
		return customers, err
	}
	req.Limit(10000)
	err = req.Json(&customersResponse)
	if err != nil {
		return customers, err
	}
	return customersResponse.Data, err
}

func (c *Client) GetCustomer(id int) (Customer, error) {
	customerResponse := GetCustomerResponse{}
	customer := Customer{}
	req, err := c.Get(fmt.Sprintf("customers/%d", id))
	if err != nil {
		return customer, err
	}
	err = req.Json(&customerResponse)
	if err != nil {
		return customer, err
	}
	return customerResponse.Data, err
}
