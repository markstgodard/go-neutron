package neutron

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const X_AUTH_TOKEN_HEADER = "X-Auth-Token"

type Network struct {
	id          string   `json:"id"`
	name        string   `json:"name"`
	description string   `json:"description"`
	status      string   `json:"status"`
	subnets     []string `json:"subnets"`
	tenant_id   string   `json:"tenant_id"`
	mtu         uint16   `json:"mtu"`
	project_id  string   `json:"project_id"`
}

type getNetworks struct {
	Networks []Network `json:"networks"`
}

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
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2.0/networks", c.URL), nil)
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

	var n getNetworks
	err = json.Unmarshal(body, &n)
	if err != nil {
		return nil, err
	}

	return n.Networks, nil
}
