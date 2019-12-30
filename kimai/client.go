package kimai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	BaseUrl    string
	Username   string
	ApiToken   string
	httpClient *http.Client
}

type ClientRequest struct {
	Client  *Client
	Request *http.Request
}

func NewClient(baseURL, userName, apiToken string) *Client {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	if !strings.HasSuffix(baseURL, "api/") {
		baseURL += "api/"
	}
	return &Client{
		BaseUrl:    baseURL,
		Username:   userName,
		ApiToken:   apiToken,
		httpClient: &http.Client{},
	}
}

func (c *Client) AppendAuthorizationAndContentType(r *http.Request) {
	r.Header.Add("X-AUTH-USER", c.Username)
	r.Header.Add("X-AUTH-TOKEN", c.ApiToken)
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

func (c *Client) Post(path string, payload interface{}) (*ClientRequest, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	r, err := http.NewRequest("POST", c.getUrl(path), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")
	c.AppendAuthorizationAndContentType(r)
	return &ClientRequest{Client: c, Request: r}, nil
}

func (cr *ClientRequest) Limit(limit int) {
	urlValues := cr.Request.URL.Query()
	urlValues.Set("size", fmt.Sprintf("%d", limit))
	cr.Request.URL.RawQuery = urlValues.Encode()
}

func (c *Client) getUrl(path string) string {
	u, _ := url.ParseRequestURI(c.BaseUrl)
	u.Path = fmt.Sprintf("%s/%s", u.Path, path)
	return u.String()
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
