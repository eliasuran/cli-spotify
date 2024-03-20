package structs

type AccessToken struct {
	Access_token string `json:"access_token"`
	Expires_in   int    `json:"expires_in"`
}

type Device struct {
	ID        string `json:"id"`
	Is_active bool   `json:"is_active"`
	Name      string `json:"name"`
}

type Devices struct {
	Devices []Device `json:"devices"`
}
