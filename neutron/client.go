package neutron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const X_AUTH_TOKEN_HEADER = "X-Auth-Token"

type request struct {
	URL          string
	Method       string
	Body         []byte
	OkStatusCode int
}

type response struct {
	Body       []byte
	StatusCode int
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

func (c *Client) doRequest(r request) (response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(r.Method, r.URL, bytes.NewBuffer(r.Body))
	if err != nil {
		return response{}, err
	}

	req.Header.Add(X_AUTH_TOKEN_HEADER, c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return response{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response{}, err
	}

	if resp.StatusCode != r.OkStatusCode {
		return response{}, fmt.Errorf("Error: %s details: %s\n", resp.Status, body)
	}
	return response{Body: body, StatusCode: resp.StatusCode}, nil
}

func (c *Client) CreateNetwork(net Network) (Network, error) {
	jsonStr, err := json.Marshal(SingleNetwork{Network: net})
	if err != nil {
		return Network{}, fmt.Errorf("invalid network: ", err)
	}

	resp, err := c.doRequest(request{
		URL:          fmt.Sprintf("%s/v2.0/networks", c.URL),
		Method:       http.MethodPost,
		Body:         jsonStr,
		OkStatusCode: http.StatusCreated,
	})
	if err != nil {
		return Network{}, err
	}

	var r SingleNetwork
	err = json.Unmarshal(resp.Body, &r)
	if err != nil {
		return Network{}, err
	}

	return r.Network, nil
}

func (c *Client) DeleteNetwork(id string) error {
	if id == "" {
		return fmt.Errorf("empty 'id' parameter")
	}
	_, err := c.doRequest(request{
		URL:          fmt.Sprintf("%s/v2.0/networks/%s", c.URL, id),
		Method:       http.MethodDelete,
		OkStatusCode: http.StatusNoContent,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Networks() ([]Network, error) {
	resp, err := c.doRequest(request{
		URL:          fmt.Sprintf("%s/v2.0/networks", c.URL),
		Method:       http.MethodGet,
		OkStatusCode: http.StatusOK,
	})
	if err != nil {
		return nil, err
	}

	var r GetNetworks
	err = json.Unmarshal(resp.Body, &r)
	if err != nil {
		return nil, err
	}
	return r.Networks, nil
}

func (c *Client) NetworksByName(name string) ([]Network, error) {
	if name == "" {
		return nil, fmt.Errorf("empty 'name' parameter")
	}

	resp, err := c.doRequest(request{
		URL:          fmt.Sprintf("%s/v2.0/networks?name=%s", c.URL, name),
		Method:       http.MethodGet,
		OkStatusCode: http.StatusOK,
	})
	if err != nil {
		return nil, err
	}

	var r GetNetworks
	err = json.Unmarshal(resp.Body, &r)
	if err != nil {
		return nil, err
	}
	return r.Networks, nil
}

func (c *Client) Subnets() ([]Subnet, error) {
	resp, err := c.doRequest(request{
		URL:          fmt.Sprintf("%s/v2.0/subnets", c.URL),
		Method:       http.MethodGet,
		OkStatusCode: http.StatusOK,
	})
	if err != nil {
		return nil, err
	}

	var r GetSubnets
	err = json.Unmarshal(resp.Body, &r)
	if err != nil {
		return nil, err
	}
	return r.Subnets, nil
}

func (c *Client) SubnetsByName(name string) ([]Subnet, error) {
	if name == "" {
		return nil, fmt.Errorf("empty 'name' parameter")
	}

	resp, err := c.doRequest(request{
		URL:          fmt.Sprintf("%s/v2.0/subnets?name=%s", c.URL, name),
		Method:       http.MethodGet,
		OkStatusCode: http.StatusOK,
	})
	if err != nil {
		return nil, err
	}

	var r GetSubnets
	err = json.Unmarshal(resp.Body, &r)
	if err != nil {
		return nil, err
	}
	return r.Subnets, nil
}

func (c *Client) CreateSubnet(s Subnet) (Subnet, error) {
	jsonStr, err := json.Marshal(SingleSubnet{Subnet: s})
	if err != nil {
		return Subnet{}, fmt.Errorf("invalid subnet: ", err)
	}

	resp, err := c.doRequest(request{
		URL:          fmt.Sprintf("%s/v2.0/subnets", c.URL),
		Method:       http.MethodPost,
		Body:         jsonStr,
		OkStatusCode: http.StatusCreated,
	})

	if err != nil {
		return Subnet{}, err
	}

	var r SingleSubnet
	err = json.Unmarshal(resp.Body, &r)
	if err != nil {
		return Subnet{}, err
	}

	return r.Subnet, nil
}

func (c *Client) CreatePort(p Port) (Port, error) {
	jsonStr, err := json.Marshal(SinglePort{Port: p})
	if err != nil {
		return Port{}, fmt.Errorf("invalid port: ", err)
	}

	resp, err := c.doRequest(request{
		URL:          fmt.Sprintf("%s/v2.0/ports", c.URL),
		Method:       http.MethodPost,
		Body:         jsonStr,
		OkStatusCode: http.StatusCreated,
	})

	if err != nil {
		return Port{}, err
	}

	var r SinglePort
	err = json.Unmarshal(resp.Body, &r)
	if err != nil {
		return Port{}, err
	}

	return r.Port, nil
}

func (c *Client) DeletePort(id string) error {
	if id == "" {
		return fmt.Errorf("empty 'id' parameter")
	}
	_, err := c.doRequest(request{
		URL:          fmt.Sprintf("%s/v2.0/ports/%s", c.URL, id),
		Method:       http.MethodDelete,
		OkStatusCode: http.StatusNoContent,
	})
	if err != nil {
		return err
	}
	return nil
}
