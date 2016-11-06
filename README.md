# go-neutron

Go library for Neutron v2.0 API

### Installation

Install using `go get github.com/markstgodard/go-neutron`.


### Usage

```go
func main() {
	// create new client
	client, err := neutron.NewClient("http://192.168.56.101:9696", "some-keystone-token")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// get networks
	networks, err := client.Networks()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, n := range networks {
		fmt.Println(n)
	}
}
```
