package kimai

import "fmt"

type Project struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Rate        float64 `json:"hourlyRate"`
	CustomerId  int     `json:"customer"`
	OrderNumber string  `json:"orderNumber"`
}

func (c *Client) GetProjects() ([]Project, error) {
	projects := make([]Project, 0)
	req, err := c.Get("projects")
	if err != nil {
		return projects, err
	}
	req.Limit(10000)
	err = req.Json(&projects)
	if err != nil {
		return projects, err
	}
	return projects, err
}

func (c *Client) GetProject(id int) (Project, error) {
	project := Project{}
	req, err := c.Get(fmt.Sprintf("projects/%d", id))
	if err != nil {
		return project, err
	}
	err = req.Json(&project)
	if err != nil {
		return project, err
	}
	return project, err
}
