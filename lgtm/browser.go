package lgtm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"

	"github.com/brunograsselli/lgtm/github"
)

type Browser struct {
	LastResultsPath string
}

func (b *Browser) Open(number int32) error {
	if _, err := os.Stat(b.LastResultsPath); os.IsNotExist(err) {
		fmt.Printf("Don't know how to open PR %d\n", number)
		return nil
	}

	c, err := ioutil.ReadFile(b.LastResultsPath)
	if err != nil {
		return err
	}

	var repos map[string][]github.PullRequest

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
