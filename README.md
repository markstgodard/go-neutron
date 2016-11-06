# go-neutron

Go library for Neutron v2.0 API

### Installation

Install using `go get github.com/markstgodard/go-neutron`.


### Usage

```go
// create new client
client, err := neutron.NewClient("http://192.168.56.101:9696", "some-keystone-token")

// get networks

networks, err := client.Networks()
```
