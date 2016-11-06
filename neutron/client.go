package neutron

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

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

type Config struct {
	URL string
}

type Client struct {
	config Config
}

func NewClient(config Config) *Client {
	return &Client{config: config}
}

func (c *Client) Networks() (Networks, error) {
	resp, err := http.Get(c.config.URL)
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
