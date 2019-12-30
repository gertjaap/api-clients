package akaunting

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type UrlValuesConvertable interface {
	ToUrlValues() *url.Values
}

type Client struct {
	BaseUrl    string
	Username   string
	Password   string
	CompanyId  int
	httpClient *http.Client
}

type ClientRequest struct {
	Client  *Client
	Request *http.Request
}

func NewClient(baseURL, userName, password string, companyId int) *Client {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	if !strings.HasSuffix(baseURL, "api/") {
		baseURL += "api/"
	}
	return &Client{
		BaseUrl:    baseURL,
		Username:   userName,
		Password:   password,
		CompanyId:  companyId,
		httpClient: &http.Client{},
	}
}

func (c *Client) AppendAuthorizationAndContentType(r *http.Request) {
	r.SetBasicAuth(c.Username, c.Password)
	r.Header.Set("Accept", "application/json")
}

func (c *Client) Get(path string) (*ClientRequest, error) {
	r, err := http.NewRequest("GET", c.getUrl(path), nil)
	if err != nil {
		return nil, err
	}
	c.AppendAuthorizationAndContentType(r)
	return &ClientRequest{Client: c, Request: r}, nil
}

func (c *Client) Post(path string, payload UrlValuesConvertable) (*ClientRequest, error) {
	payloadString := payload.ToUrlValues().Encode()
	//log.Printf("Sending payload:\n%s\n\n", payloadString)
	r, err := http.NewRequest("POST", c.getUrl(path), strings.NewReader(payloadString))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.AppendAuthorizationAndContentType(r)
	return &ClientRequest{Client: c, Request: r}, nil
}

func (c *Client) getUrl(path string) string {
	u, _ := url.ParseRequestURI(c.BaseUrl)
	u.Path = fmt.Sprintf("%s/%s", u.Path, path)
	u.RawQuery = fmt.Sprintf("company_id=%d", c.CompanyId)
	return u.String()
}

func (cr *ClientRequest) Limit(limit int) {
	urlValues := cr.Request.URL.Query()
	urlValues.Set("limit", fmt.Sprintf("%d", limit))
	cr.Request.URL.RawQuery = urlValues.Encode()
}

func (cr *ClientRequest) Json(dst interface{}) error {
	resp, err := cr.Client.httpClient.Do(cr.Request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	/*b, _ := ioutil.ReadAll(resp.Body)
	log.Printf("Response: %s", string(b))
	return json.NewDecoder(bytes.NewBuffer(b)).Decode(dst)*/

	return json.NewDecoder(resp.Body).Decode(dst)
}

func (c *Client) AkauntingDate(d time.Time) string {
	return d.Format("2006-01-02 15:04:05")
}

func (cr *ClientRequest) ExpectStatus(statusCode int) error {
	resp, err := cr.Client.httpClient.Do(cr.Request)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != statusCode {
		return fmt.Errorf("Received unexpected status code %d - %s", resp.StatusCode, resp.Status)
	}
	return nil
}
