package lgtm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/ssh/terminal"
)

type Authorization struct {
	Token string `json:"token"`
}

func Login() error {
	credentialsPath := fmt.Sprintf("%s/.lgtm.secret", os.Getenv("HOME"))

	if _, err := os.Stat(credentialsPath); err == nil {
		fmt.Println("You are already logged in.")
		return nil
	}

	user, password, err := credentials()

	if err != nil {
		return err
	}

	client := &http.Client{}

	fingerprint := ksuid.New().String()
	reqBody := []byte(fmt.Sprintf(`{"note":"lgtm","scopes":["repo"],"fingerprint":"%s"}`, fingerprint))

	req, err := http.NewRequest("POST", "https://api.github.com/authorizations", bytes.NewBuffer(reqBody))
	req.SetBasicAuth(user, password)

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 401 && resp.Header["X-Github-Otp"] != nil {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter 2FA authentication code: ")
		code, _ := reader.ReadString('\n')
		code = strings.TrimSpace(code)

		req, err := http.NewRequest("POST", "https://api.github.com/authorizations", bytes.NewBuffer(reqBody))
		req.SetBasicAuth(user, password)
		req.Header.Add("X-GitHub-OTP", code)

		resp, err = client.Do(req)

		if err != nil {
			return err
		}

		defer resp.Body.Close()

		body, err = ioutil.ReadAll(resp.Body)
	}

	if resp.StatusCode != 201 {
		fmt.Println(string(body))
		return nil
	}

	var auth Authorization
	json.Unmarshal(body, &auth)

	err = ioutil.WriteFile(credentialsPath, []byte(auth.Token), 0644)
	if err != nil {
		return err
	}
	fmt.Println("Success!")
	return nil
}

func credentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your GitHub user: ")
	user, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(0)
	fmt.Println("")
	if err != nil {
		return "", "", err
	}
	password := string(bytePassword)

	return strings.TrimSpace(user), strings.TrimSpace(password), nil
}
