package neutron

type Subnet struct {
	ID             string
	Name           string
	EnableDHCP     bool
	NetworkID      string
	SegmentID      string
	ProjectID      string
	TenantID       string
	DNSNameservers []string
	// AllocationPools []AllocationPool
	HostRoutes []string
	IPVersion  int
	GatewayIP  string
	CIDR       string
}
