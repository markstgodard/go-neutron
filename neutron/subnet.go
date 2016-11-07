package neutron

type Subnet struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	SubnetPoolID    string           `json:"subnetpool_id"`
	EnableDHCP      bool             `json:"enable_dhcp"`
	NetworkID       string           `json:"network_id"`
	SegmentID       string           `json:"segment_id"`
	ProjectID       string           `json:"project_id"`
	TenantID        string           `json:"tenant_id"`
	DNSNameservers  []string         `json:"dns_nameservers"`
	AllocationPools []AllocationPool `json:"allocation_pools"`
	HostRoutes      []string         `json:"host_routes"`
	IPVersion       int              `json:"ip_version"`
	GatewayIP       string           `json:"gateway_ip"`
	CIDR            string           `json:"cidr"`
}

type AllocationPool struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type GetSubnets struct {
	Subnets []Subnet `json:"subnets"`
}
