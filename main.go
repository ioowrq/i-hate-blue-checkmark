package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/dghubble/oauth1"
	"github.com/mrjones/oauth"
)

const (
	callbackURL = "http://localhost:4673/callback"
	tokenFile   = "twitter_token.json"
	credsFile   = "credentials.json"
)

type Credentials struct {
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}

type TokenData struct {
	Token  string `json:"token"`
	Secret string `json:"secret"`
}

type UserResponse struct {
	ProfileImageURL string `json:"profile_image_url"`
}

func getCredentials() Credentials {
	if data, err := ioutil.ReadFile(credsFile); err == nil {
		var creds Credentials
		if err := json.Unmarshal(data, &creds); err == nil {
			return creds
		}
	}

	credentials := os.Getenv("CREDENTIALS")
	if credentials != "" {
		var creds Credentials
		if err := json.Unmarshal([]byte(credentials), &creds); err == nil {
			return creds
		}
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter API Key: ")
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)

	fmt.Print("Enter API Secret: ")
	apiSecret, _ := reader.ReadString('\n')
	apiSecret = strings.TrimSpace(apiSecret)

	creds := Credentials{
		APIKey:    apiKey,
		APISecret: apiSecret,
	}

	if data, err := json.Marshal(creds); err == nil {
		ioutil.WriteFile(credsFile, data, 0600)
	}

	return creds
}

func main() {
	if len(os.Args) > 1 {
		log.Fatal("This program doesn't accept arguments")
	}

	creds := getCredentials()
	token := os.Getenv("TOKEN")

	if _, err := os.Stat(tokenFile); os.IsNotExist(err) && token == "" {
		handleAuth(creds)
	} else {
		updateProfileImage(creds)
	}
}

func handleAuth(creds Credentials) {
	consumer := oauth.NewConsumer(
		creds.APIKey,
		creds.APISecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		},
	)

	requestToken, url, err := consumer.GetRequestTokenAndUrl(callbackURL)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan TokenData)
	server := &http.Server{Addr: ":4673"}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		verifier := r.URL.Query().Get("oauth_verifier")
		accessToken, err := consumer.AuthorizeToken(requestToken, verifier)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tokenData := TokenData{
			Token:  accessToken.Token,
			Secret: accessToken.Secret,
		}

		saveToken(tokenData)

		fmt.Fprintf(w, "Authentication successful! You can close this window.")
		done <- tokenData
	})

	fmt.Printf("Please visit: %s\n", url)
	openURL(url)
	go server.ListenAndServe()

	<-done
	server.Close()

}

func updateProfileImage(creds Credentials) {
	tokenData := loadToken()

	config := oauth1.NewConfig(creds.APIKey, creds.APISecret)
	token := oauth1.NewToken(tokenData.Token, tokenData.Secret)
	httpClient := config.Client(oauth1.NoContext, token)

	resp, err := httpClient.Get("https://api.twitter.com/1.1/account/verify_credentials.json")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var user UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		log.Fatal(err)
	}

	imageURL := strings.Replace(user.ProfileImageURL, "_normal", "", 1)
	imgResp, err := http.Get(imageURL)
	if err != nil {
		log.Fatal(err)
	}
	defer imgResp.Body.Close()

	imgData, err := ioutil.ReadAll(imgResp.Body)
	if err != nil {
		log.Fatal(err)
	}

	base64Img := base64.StdEncoding.EncodeToString(imgData)

	form := url.Values{}
	form.Add("image", base64Img)

	req, err := http.NewRequest("POST",
		"https://api.twitter.com/1.1/account/update_profile_image.json",
		strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err = httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Update failed: %s", resp.Status)
	}

	fmt.Println("Profile image updated successfully!")
}

func saveToken(token TokenData) error {
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(tokenFile, data, 0600)
}

func loadToken() TokenData {
	var data []byte
	var err error
	var token TokenData
	tokenValue := os.Getenv("TOKEN")
	if tokenValue != "" {
		data = []byte(tokenValue)
		if err := json.Unmarshal(data, &token); err != nil {
			log.Fatal(err)
		}
	} else {
		data, err = ioutil.ReadFile(tokenFile)
		if err != nil {
			log.Fatal("Please authenticate first")
		}

		if err := json.Unmarshal(data, &token); err != nil {
			log.Fatal(err)
		}
	}
	return token
}

func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // "linux", "freebsd", "openbsd", "netbsd"
		// Check if running under WSL
		if isWSL() {
			// Use 'cmd.exe /c start' to open the URL in the default Windows browser
			cmd = "cmd.exe"
			args = []string{"/c", "start", url}
		} else {
			// Use xdg-open on native Linux environments
			cmd = "xdg-open"
			args = []string{url}
		}
	}
	if len(args) > 1 {
		// args[0] is used for 'start' command argument, to prevent issues with URLs starting with a quote
		args = append(args[:1], append([]string{""}, args[1:]...)...)
	}
	return exec.Command(cmd, args...).Start()
}

// isWSL checks if the Go program is running inside Windows Subsystem for Linux
func isWSL() bool {
	releaseData, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(releaseData)), "microsoft")
}
