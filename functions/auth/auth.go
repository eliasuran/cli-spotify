package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// getting auth token
var authCodeChan = make(chan string)

func GetAuthCode(client_id string, redirect_uri string, scopes []string) string {
	auth_code_url := "https://accounts.spotify.com/authorize?response_type=code&client_id=" + client_id + "&scope=" + strings.Join(scopes, " ") + "&redirect_uri=" + redirect_uri

	err := openBrowser(auth_code_url)
	if err != nil {
		fmt.Printf("Could not open browser: %v\n", err)
		return ""
	}

	http.HandleFunc("/callback", handleCallback)
	go http.ListenAndServe(":8888", nil)

	code := getCode()

	return code
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

type Response struct {
	Access_token  string `json:"access_token"`
	Token_type    string `json:"token_type"`
	Scope         string `json:"scope"`
	Expires_in    int    `json:"expires_in"`
	Refresh_token string `json:"refresh_token"`
}

// getting refresh token
func GetRefreshToken(client_id string, client_secret string, redirect_uri string, auth_code string) string {
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
		fmt.Printf("Error initializing request to refresh token: %v\n", err)
		return ""
	}
	// headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+encoded_auth)

	// execute request
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error executing request to get refresh token: %v\n", err)
		return ""
	}

	// read request
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return ""
	}

	// parse json
	var data Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Printf("Error parsing json: %v\n", err)
		return ""
	}

	return data.Refresh_token
}
