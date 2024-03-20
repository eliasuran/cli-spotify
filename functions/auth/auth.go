package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// getting auth token
var authCodeChan = make(chan string)

func GetAuthCode(client_id string, redirect_uri string, scopes []string) (string, error) {
	auth_code_url := "https://accounts.spotify.com/authorize?response_type=code&client_id=" + client_id + "&scope=" + strings.Join(scopes, " ") + "&redirect_uri=" + redirect_uri

	err := openBrowser(auth_code_url)
	if err != nil {
		return "", err
	}

	http.HandleFunc("/callback", handleCallback)
	go http.ListenAndServe(":8888", nil)

	code := getCode()

	return code, nil
}

// check os and open browser, returns err
func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Printf("Could not open browser window: %v\n", err)
	}
	return err
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Auth code not found", http.StatusBadRequest)
	}
	authCodeChan <- code
	fmt.Fprintln(w, "Auth successful! You can close this window")
}

func getCode() string {
	for {
		code := <-authCodeChan
		if code != "" {
			return code
		}
		time.Sleep(1 * time.Second)
	}
}

type RefreshToken struct {
	Access_token  string `json:"access_token"`
	Token_type    string `json:"token_type"`
	Scope         string `json:"scope"`
	Expires_in    int    `json:"expires_in"`
	Refresh_token string `json:"refresh_token"`
}

// getting refresh token
func GetRefreshToken(client_id string, client_secret string, redirect_uri string, auth_code string) (string, error) {
	// b64 encode
	auth := client_id + ":" + client_secret
	encoded_auth := base64.StdEncoding.EncodeToString([]byte(auth))

	// initialize request
	// body
	req_body_data := url.Values{}
	req_body_data.Set("grant_type", "authorization_code")
	req_body_data.Set("code", auth_code)
	req_body_data.Set("redirect_uri", redirect_uri)
	req_body := strings.NewReader(req_body_data.Encode())

	client := http.Client{}
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", req_body)
	if err != nil {
		return "", err
	}
	// headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+encoded_auth)

	// execute request
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	// read request
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	// parse json
	var data RefreshToken
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	return data.Refresh_token, nil
}

func WriteToEnv(refresh_token string) error {
	data := "SPOTIFY_REFRESH_TOKEN=" + refresh_token
	f, err := os.OpenFile("./.env", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(data); err != nil {
		return err
	}
	return nil
}

type AccessToken struct {
	Access_token string `json:"access_token"`
}

func GetAccessToken(client_id string, client_secret string, refresh_token string) (string, error) {
	auth := client_id + ":" + client_secret
	encoded_auth := base64.StdEncoding.EncodeToString([]byte(auth))

	// body
	req_body_data := url.Values{}
	req_body_data.Set("grant_type", "refresh_token")
	req_body_data.Set("refresh_token", refresh_token)
	req_body := strings.NewReader(req_body_data.Encode())

	client := http.Client{}
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", req_body)
	if err != nil {
		return "", err
	}
	// headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+encoded_auth)

	// execute request
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	// read request
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	// parse json
	var data AccessToken
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	return data.Access_token, nil
}
