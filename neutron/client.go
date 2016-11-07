package neutron

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const X_AUTH_TOKEN_HEADER = "X-Auth-Token"

type Client struct {
	URL   string
	token string
}

func NewClient(url, token string) (*Client, error) {
	if url == "" {
		return nil, fmt.Errorf("missing URL")
	}
	if token == "" {
		return nil, fmt.Errorf("missing token")
	}
	return &Client{URL: url, token: token}, nil
}

func (c *Client) Networks() ([]Network, error) {
	return c.getNetworks(fmt.Sprintf("%s/v2.0/networks", c.URL))
}

func (c *Client) NetworksByName(name string) ([]Network, error) {
	if name == "" {
		return nil, fmt.Errorf("empty 'name' parameter")
	}
	return c.getNetworks(fmt.Sprintf("%s/v2.0/networks?name=%s", c.URL, name))
}

func (c *Client) getNetworks(url string) ([]Network, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add(X_AUTH_TOKEN_HEADER, c.token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: %s\n", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)

	var gn GetNetworks
	err = json.Unmarshal(body, &gn)
	if err != nil {
		return nil, err
	}

	return gn.Networks, nil
}
