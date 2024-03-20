package structs

type Device struct {
	ID        string `json:"id"`
	Is_active bool   `json:"is_active"`
	Name      string `json:"name"`
}

type Devices struct {
	Devices []Device `json:"devices"`
}
