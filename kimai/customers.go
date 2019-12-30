package kimai

import "fmt"

type Customer struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Company  string  `json:"company"`
	Contact  string  `json:"contact"`
	Address  string  `json:"address"`
	Country  string  `json:"country"`
	Email    string  `json:"email"`
	Rate     float64 `json:"hourlyRate"`
	Currency string  `json:"currency"`
	Number   string  `json:"number"`
}

func (c *Client) GetCustomers() ([]Customer, error) {
	customers := make([]Customer, 0)
	req, err := c.Get("customers")
	if err != nil {
		return customers, err
	}
	req.Limit(10000)
	err = req.Json(&customers)
	if err != nil {
		return customers, err
	}
	return customers, err
}

func (c *Client) GetCustomer(id int) (Customer, error) {
	customer := Customer{}
	req, err := c.Get(fmt.Sprintf("customers/%d", id))
	if err != nil {
		return customer, err
	}
	err = req.Json(&customer)
	if err != nil {
		return customer, err
	}
	return customer, err
}
