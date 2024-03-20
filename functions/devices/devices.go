package devices

import (
	"github.com/eliasuran/cli-spotify/functions/requests"
	"github.com/eliasuran/cli-spotify/functions/structs"
)

func GetAllDevices(token string) ([]structs.Device, error) {
	var deviceData structs.Devices
	err := requests.Get(token, "https://api.spotify.com/v1/me/player/devices", &deviceData)
	if err != nil {
		return []structs.Device{}, err
	}
	return deviceData.Devices, nil
}
