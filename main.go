package main

import (
	"fmt"
	"os"

	"github.com/eliasuran/cli-spotify/functions/auth"
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

	refresh_token := os.Getenv("SPOTIFY_REFRESH_TOKEN")

	if refresh_token == "" {
		client_id := os.Getenv("SPOTIFY_CLIENT_ID")
		client_secret := os.Getenv("SPOTIFY_CLIENT_SECRET")

		if client_id == "" || client_secret == "" {
			fmt.Println("One or more env variables missing")
			os.Exit(1)
		}

		scopes := []string{"user-top-read", "user-read-currently-playing"}
		redirect_uri := "http://localhost:8888/callback"

		// get auth code
		auth_code, err := auth.GetAuthCode(client_id, redirect_uri, scopes)
		if err != nil {
			fmt.Printf("Error getting auth code: %v\n", err)
			return
		}

		// get refresh token
		refresh_token, err = auth.GetRefreshToken(client_id, client_secret, redirect_uri, auth_code)
		if err != nil {
			fmt.Printf("Error getting refresh token: %v\n", err)
			return
		}

		// write the refresh token to env
		err = auth.WriteToEnv(refresh_token)
		if err != nil {
			fmt.Printf("Error writing to env file: %v\n", err)
			return
		}
	}
	fmt.Println("You are authorized:", refresh_token)
}
