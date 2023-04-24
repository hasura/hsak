package commands

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

const HSAK_CURRENT_VERSION string = "v1.0.1"
const HSAK_LATEST_VERSION_API_PATH string = "https://api.github.com/repos/hasura/hsak/releases/latest"

func init() {
	rootCommand.AddCommand(runVersionCommand)
}

var runVersionCommand = &cobra.Command{
	Use:   "version",
	Short: "HSAK version and check for updates",
	RunE: func(command *cobra.Command, args []string) error {
		return runVersion()
	},
}

func runVersion() error {
	httpResponse, err := http.Get(HSAK_LATEST_VERSION_API_PATH)
	if err != nil {
		return fmt.Errorf("error retrieve lastest version info from GitHub: %w", err)
	}

	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return fmt.Errorf("error reading HTTP response: %w", err)
	}
	latestVersion := gjson.GetBytes(responseBody, "tag_name").String()
	if strings.Compare(latestVersion, HSAK_CURRENT_VERSION) == 0 {
		fmt.Printf("The HSAK version is '%s' and it is the latest released version\n", HSAK_CURRENT_VERSION)
	} else {
		fmt.Printf("The HSAK version is '%s' and there is newer released version '%s'\n", HSAK_CURRENT_VERSION, latestVersion)
	}

	return nil
}
