package neutron

type Network struct {
	ID           string   `json:"id,omitempty"`
	Name         string   `json:"name"`
	Description  string   `json:"description,omitempty"`
	Status       string   `json:"status,omitempty"`
	AdminStateUp bool     `json:"admin_state_up"`
	Subnets      []string `json:"subnets,omitempty"`
	TenantID     string   `json:"tenant_id,omitempty"`
	MTU          int      `json:"mtu,omitempty"`
	ProjectID    string   `json:"project_id,omitempty"`
}

type GetNetworks struct {
	Networks []Network `json:"networks"`
}

type SingleNetwork struct {
	Network Network `json:"network"`
}
