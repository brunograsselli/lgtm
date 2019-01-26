package lgtm

import (
	"bufio"
	"fmt"
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

	client, err := github.Authorize(user, password, fingerprint)
	if err != nil {
		switch err {
		case github.Need2FAErr:
			code := askFor2FACode()

			client, err = github.AuthorizeWith2FA(user, password, fingerprint, code)

			if err != nil {
				return err
			}
		default:
			return err
		}
	}

	if err := secrets.SaveToken(client.GetToken()); err != nil {
		return err
	}

	if err := config.SaveUserName(user); err != nil {
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
