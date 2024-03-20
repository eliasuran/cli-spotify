package main

import (
	"fmt"
	"os"

	"github.com/eliasuran/cli-spotify/functions/auth"
	"github.com/eliasuran/cli-spotify/functions/playback"
	"github.com/joho/godotenv"
)

func check(message string, err error) {
	if err != nil {
		fmt.Printf(message, err)
		os.Exit(1)
	}
}

func main() {
	err := godotenv.Load()
	check("Error loading env variables: %v", err)

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
		check("Error getting auth code: %v\n", err)

		// get refresh token
		refresh_token, err = auth.GetRefreshToken(client_id, client_secret, redirect_uri, auth_code)
		check("Error getting refresh token: %v\n", err)

		// write the refresh token to env
		err = auth.WriteToEnv(refresh_token)
		check("Error writing to env file: %v\n", err)
	}

	access_token, err := auth.GetAccessToken(client_id, client_secret, refresh_token)
	check("Could not get access token: %v\n", err)

	// devices, err := devices.GetAllDevices(access_token.Access_token)

	err = playback.ResumePlayback(access_token.Access_token)
	check("Could not resume playback: %v\n", err)
}
