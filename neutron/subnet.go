package neutron

type Subnet struct {
	ID              string           `json:"id,omitempty"`
	Name            string           `json:"name,omitempty"`
	SubnetPoolID    string           `json:"subnetpool_id,omitempty"`
	EnableDHCP      bool             `json:"enable_dhcp,omitempty"`
	NetworkID       string           `json:"network_id"`
	SegmentID       string           `json:"segment_id,omitempty"`
	ProjectID       string           `json:"project_id,omitempty"`
	TenantID        string           `json:"tenant_id,omitempty"`
	DNSNameservers  []string         `json:"dns_nameservers,omitempty"`
	AllocationPools []AllocationPool `json:"allocation_pools,omitempty"`
	HostRoutes      []string         `json:"host_routes,omitempty"`
	IPVersion       int              `json:"ip_version"`
	GatewayIP       string           `json:"gateway_ip,omitempty"`
	CIDR            string           `json:"cidr"`
}

type AllocationPool struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type GetSubnets struct {
	Subnets []Subnet `json:"subnets"`
}

type SingleSubnet struct {
	Subnet Subnet `json:"subnet"`
}
