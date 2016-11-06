package neutron

type Network struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	Subnets     []string `json:"subnets"`
	TenantID    string   `json:"tenant_id"`
	MTU         int      `json:"mtu"`
	ProjectID   string   `json:"project_id"`
}

type GetNetworks struct {
	Networks []Network `json:"networks"`
}
