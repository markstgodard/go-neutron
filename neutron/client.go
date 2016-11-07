package neutron

import (
	"bytes"
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

func (c *Client) CreateNetwork(net Network) (Network, error) {
	return c.postNetwork(fmt.Sprintf("%s/v2.0/networks", c.URL), net)
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

	var r GetNetworks
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	return r.Networks, nil
}

func (c *Client) postNetwork(url string, net Network) (Network, error) {
	client := &http.Client{}
	jsonStr, err := json.Marshal(net)
	if err != nil {
		return Network{}, fmt.Errorf("invalid network: ", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add(X_AUTH_TOKEN_HEADER, c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return Network{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Network{}, fmt.Errorf("Error: %s\n", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)

	var r GetNetwork
	err = json.Unmarshal(body, &r)
	if err != nil {
		return Network{}, err
	}

	return r.Network, nil
}

func (c *Client) Subnets() ([]Subnet, error) {
	return c.getSubnets(fmt.Sprintf("%s/v2.0/subnets", c.URL))
}

func (c *Client) SubnetsByName(name string) ([]Subnet, error) {
	if name == "" {
		return nil, fmt.Errorf("empty 'name' parameter")
	}
	return c.getSubnets(fmt.Sprintf("%s/v2.0/subnets?name=%s", c.URL, name))
}

func (c *Client) getSubnets(url string) ([]Subnet, error) {
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

	var r GetSubnets
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	return r.Subnets, nil
}
