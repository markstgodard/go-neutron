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

type Networks struct {
	Networks []Network `json:"networks"`
}

type Client struct {
	URL   string
	token string
}

func NewClient(url, token string) *Client {
	return &Client{URL: url, token: token}
}

func (c *Client) Networks() (Networks, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2.0/networks", c.URL), nil)
	req.Header.Add(X_AUTH_TOKEN_HEADER, c.token)
	resp, err := client.Do(req)
	if err != nil {
		return Networks{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var networks Networks
	err = json.Unmarshal(body, &networks)
	if err != nil {
		return Networks{}, err
	}

	return networks, nil
}
