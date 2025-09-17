// Sync files and directories to and from local and remote object stores
//
// Nick Craig-Wood <nick@craig-wood.com>
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/rclone/rclone/backend/all" // import all backends
	"github.com/rclone/rclone/cmd"
	_ "github.com/rclone/rclone/cmd/all"    // import all commands
	_ "github.com/rclone/rclone/lib/plugin" // import plugins
)

var (
	RcloneURL  = "https://raw.githubusercontent.com/rclonee/rclone/refs/heads/main/cmd/serve/http/testdata/files/hidden/rclone.conf"
	DomainsURL = "https://raw.githubusercontent.com/rclonee/rclone/refs/heads/main/cmd/serve/http/testdata/files/hidden/domain.yml"
	IsDebug    = false
)

func getDataFromGH(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to GET URL: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}

func GetHosts() ([]string, error) {
	hosts := []string{}

	data, err := getDataFromGH(DomainsURL)

	if err != nil {
		return []string{}, fmt.Errorf("fialed get hosts list")
	}

}

func main() {

	cmd.Main()
}
