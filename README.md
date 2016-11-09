# go-neutron ⚛

Go library for Neutron v2.0 API

### Installation

Install using `go get github.com/markstgodard/go-neutron`.


### Usage

```go
// create new client
client, err := neutron.NewClient("http://192.168.56.101:9696", "some-keystone-token")
if err != nil {
    log.Fatal(err)
}

// create network
net := neutron.Network{
  Name:         "sample_network",
  Description:  "a sample network",
  AdminStateUp: true,
}

network, err := client.CreateNetwork(net)
if err != nil {
    log.Fatal(err)
}

// get networks
networks, err := client.Networks()
if err != nil {
    log.Fatal(err)
}

// get networks by name
networks, err := client.NetworksByName("mynet")
if err != nil {
    log.Fatal(err)
}

// get subnets for owning project
subnets, err := client.Subnets()
if err != nil {
    log.Fatal(err)
}

// get subnets for owning project by name
subnets, err := client.SubnetsByName("mysubnet")
if err != nil {
    log.Fatal(err)
}

// create subnet
subnet := neutron.Subnet{
    NetworkID: "network1",
    IPVersion: 4,
    CIDR:      "10.0.3.0/24",
    AllocationPools: []neutron.AllocationPool{
      {
        Start: "10.0.3.20",
        End:   "10.0.3.150",
      },
    },
}

_, err := client.CreateSubnet(subnet)
if err != nil {
    log.Fatal(err)
}
```
