package lgtm

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/brunograsselli/lgtm/github"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/ssh/terminal"
)

func Login(secrets *Secrets, config *Config) error {
	token, _ := secrets.Token()

	if token != nil {
		fmt.Println("You are already logged in.")
		return nil
	}

	user, password, err := askForCredentials()

	if err != nil {
		return err
	}

	fingerprint := ksuid.New().String()

	resp, err := github.Authorize(user, password, fingerprint, "")

	if err != nil {
		return err
	}

	if resp.StatusCode == 401 && resp.Header["X-Github-Otp"] != nil {
		code := askFor2FACode()

		resp, err = github.Authorize(user, password, fingerprint, code)

		if err != nil {
			return err
		}
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 201 {
		return errors.New(string(body))
	}

	var auth Authorization
	json.Unmarshal(body, &auth)

	err = secrets.SaveToken(auth.Token)

	if err != nil {
		return err
	}

	err = config.SaveUserName(user)

	if err == nil {
		fmt.Println("You are logged in!")
	}

	return err
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
