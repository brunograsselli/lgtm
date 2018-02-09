package lgtm

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/ssh/terminal"
)

func Login() error {
	credentialsPath := fmt.Sprintf("%s/.lgtm.secret", os.Getenv("HOME"))

	if _, err := os.Stat(credentialsPath); err == nil {
		fmt.Println("You are already logged in.")
		return nil
	}

	user, password, err := askForCredentials()

	if err != nil {
		return err
	}

	fingerprint := ksuid.New().String()

	resp, err := GitHubAuthoriza(user, password, fingerprint, "")

	if err != nil {
		return err
	}

	if resp.StatusCode == 401 && resp.Header["X-Github-Otp"] != nil {
		code := askFor2FACode()

		resp, err = GitHubAuthoriza(user, password, fingerprint, code)

		if err != nil {
			return err
		}
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

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

func askForCredentials() (string, string, error) {
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

func askFor2FACode() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter 2FA authentication code: ")
	code, _ := reader.ReadString('\n')
	return strings.TrimSpace(code)
}
