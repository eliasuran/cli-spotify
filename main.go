package main

import (
	"fmt"
	"os"

	"github.com/eliasuran/cli-spotify/functions/auth"
	"github.com/eliasuran/cli-spotify/functions/devices"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading env variables: %v", err)
		os.Exit(1)
	}

	// check if there is a refresh token env variable
	// if there is none, run the auth process

	client_id := os.Getenv("SPOTIFY_CLIENT_ID")
	client_secret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if client_id == "" || client_secret == "" {
		fmt.Println("One or more env variables missing")
		os.Exit(1)
	}

	refresh_token := os.Getenv("SPOTIFY_REFRESH_TOKEN")

	// get refresh token if none is found
	if refresh_token == "" {
		scopes := []string{"user-read-playback-state"}
		redirect_uri := "http://localhost:8888/callback"

		// get auth code
		auth_code, err := auth.GetAuthCode(client_id, redirect_uri, scopes)
		if err != nil {
			fmt.Printf("Error getting auth code: %v\n", err)
			os.Exit(1)
		}

		// get refresh token
		refresh_token, err = auth.GetRefreshToken(client_id, client_secret, redirect_uri, auth_code)
		if err != nil {
			fmt.Printf("Error getting refresh token: %v\n", err)
			os.Exit(1)
		}

		// write the refresh token to env
		err = auth.WriteToEnv(refresh_token)
		if err != nil {
			fmt.Printf("Error writing to env file: %v\n", err)
			os.Exit(1)
		}
	}

	access_token, err := auth.GetAccessToken(client_id, client_secret, refresh_token)
	if err != nil {
		fmt.Printf("Could not get access token: %v\n", err)
		return
	}

	devices, err := devices.GetAllDevices(access_token)
	if err != nil {
		fmt.Printf("Could not get devices: %v\n", err)
	}

	fmt.Println("devices: ", devices)
}
