package lgtm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
)

func Open(number int32) error {
	credentialsPath := fmt.Sprintf("%s/.lgtm.secret", os.Getenv("HOME"))

	if _, err := os.Stat(credentialsPath); os.IsNotExist(err) {
		fmt.Println("Please log in first (lgtm login)")
		return nil
	}

	if _, err := os.Stat("/tmp/lgtm.json"); os.IsNotExist(err) {
		fmt.Printf("Don't know how to open PR %d\n", number)
		return nil
	}

	c, err := ioutil.ReadFile("/tmp/lgtm.json")
	if err != nil {
		return err
	}

	var repos map[string][]PullRequest

	json.Unmarshal(c, &repos)

	for _, prs := range repos {
		for _, pr := range prs {
			if pr.Number == number {
				err := openBrowser(pr.HTMLURL)

				if err != nil {
					return err
				}

				return nil
			}
		}
	}

	fmt.Printf("Don't know how to open PR %d\n", number)
	return nil
}

func openBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}
