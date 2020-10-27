package luci

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// HostURL - Default Luci URL
const HostURL string = "http://localhost"

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

// RequestStruct -
type RequestStruct struct {
	ID int `json:"id"`
	Method string `json:"method"`
	Params []string `json:"params"`
}

// AuthResponse -
type AuthResponse struct {
	ID     int    `json:"id"`
	Result string `json:"result"`
	Error  string `json:"error"`
}

// NewClient -
func NewClient(host, username, password *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		// Default Luci URL
		HostURL: HostURL,
	}

	if host != nil {
		c.HostURL = *host
	}

	if (username != nil) && (password != nil) {
		// form request body
		rb, err := json.Marshal(RequestStruct{
			ID: 1,
			Method: "login",
			Params: []string{*username, *password},
		})

		if err != nil {
			return nil, err
		}

		// authenticate
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/cgi-bin/luci/rpc/auth", c.HostURL), strings.NewReader(string(rb)))
		if err != nil {
			return nil, err
		}

		body, err := c.doRequest(req)

		// parse response body
		ar := AuthResponse{}
		err = json.Unmarshal(body, &ar)
		if err != nil {
			return nil, err
		}

		c.Token = ar.Result
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	q := req.URL.Query()
	q.Add("auth", c.Token)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
