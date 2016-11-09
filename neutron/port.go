package neutron

type Port struct {
	ID           string    `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	NetworkID    string    `json:"network_id"`
	TenantID     string    `json:"tenant_id,omitempty"`
	Status       string    `json:"status,omitempty"`
	AdminStateUp bool      `json:"admin_state_up,omitempty"`
	MacAddress   string    `json:"mac_address,omitempty"`
	DeviceOwner  string    `json:"device_owner,omitempty"`
	DeviceID     string    `json:"device_id,omitempty"`
	FixedIPs     []FixedIP `json:"fixed_ips,omitempty"`
}

type FixedIP struct {
	IPAddress string `json:"ip_address"`
	SubnetID  string `json:"subnet_id"`
}

type GetPorts struct {
	Ports []Port `json:"ports"`
}

type SinglePort struct {
	Port Port `json:"port"`
}
