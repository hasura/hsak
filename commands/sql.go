package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

func init() {
	runSQLCommand.Flags().StringP("fileURI", "f", "", "Path to the SQL file to run.")
	runSQLCommand.MarkFlagRequired("fileURI")
	runSQLCommand.Flags().StringP("url", "u", "", "Hasura GraphQL Engine base URL")
	runSQLCommand.MarkFlagRequired("url")
	runSQLCommand.Flags().StringP("hasuraAdminSecret", "S", "", "Hasura Admin Secret")
	runSQLCommand.MarkFlagRequired("hasuraAdminSecret")
	runSQLCommand.Flags().StringP("dataSourceName", "s", "", "Hasura GraphQL Engine data source name")
	runSQLCommand.MarkFlagRequired("dataSourceName")

	rootCommand.AddCommand(runSQLCommand)
}

var runSQLCommand = &cobra.Command{
	Use:   "sql",
	Short: "Run SQL statements from a file",
	RunE: func(command *cobra.Command, args []string) error {
		fileURI, _ := command.Flags().GetString("fileURI")
		url, _ := command.Flags().GetString("url")
		hasuraAdminSecret, _ := command.Flags().GetString("hasuraAdminSecret")
		dataSourceName, _ := command.Flags().GetString("dataSourceName")
		gitOptions := initGitOptionsFromCommand(command)

		return runSQL(fileURI, url, hasuraAdminSecret, dataSourceName, gitOptions)
	},
}

func ResetCLISqlFlags() {
	runSQLCommand.Flags().Lookup("fileURI").Value.Set("")
	runSQLCommand.Flags().Lookup("fileURI").Changed = false
	runSQLCommand.Flags().Lookup("url").Value.Set("")
	runSQLCommand.Flags().Lookup("url").Changed = false
	runSQLCommand.Flags().Lookup("hasuraAdminSecret").Value.Set("")
	runSQLCommand.Flags().Lookup("hasuraAdminSecret").Changed = false
	runSQLCommand.Flags().Lookup("dataSourceName").Value.Set("")
	runSQLCommand.Flags().Lookup("dataSourceName").Changed = false
}

func runSQL(fileURI, url, hasuraAdminSecret, dataSourceName string, gitOptions *gitOptions) error {
	fmt.Printf("Executing script on data source \"%s\" for HGE at %s\n", dataSourceName, url)
	sqlBytes := &bytes.Buffer{}
	fileURI, err := readFile(fileURI, "", "", gitOptions, sqlBytes)
	if err != nil {
		return fmt.Errorf("error reading file: %s %w", fileURI, err)
	}
	if gitOptions != nil && len(gitOptions.Uri) > 0 {
		fmt.Printf("  1:2 Reading script file %s from %s \u2714\n", fileURI, gitOptions.Uri)
	} else {
		fmt.Printf("  1:2 Reading script file %s \u2714\n", fileURI)
	}

	payload := createQueryAPIPayload(dataSourceName, sqlBytes.Bytes())

	httpRequest, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}
	httpRequest.Header.Add("content-type", "application/json")
	httpRequest.Header.Add("x-hasura-admin-secret", hasuraAdminSecret)

	fmt.Printf("  2:2 Executing query ")
	_, httpResponseBody, err := callHasuraAPI(url, hasuraAdminSecret, payload)
	if err != nil {
		return fmt.Errorf("\nerror importing metadata, %w", err)
	}

	responseBody, err := io.ReadAll(httpResponseBody)
	if err != nil {
		return fmt.Errorf("\nerror reading HTTP response: %w", err)
	}
	result := gjson.GetBytes(responseBody, "result").String()
	if len(result) > 0 {
		fmt.Printf("\u2714\n  result: %s\n", result)
	} else {
		fmt.Println("\u2714")
	}

	return nil
}

func createQueryAPIPayload(dataSourceName string, sql []byte) []byte {
	escapedSql, err := json.Marshal(string(sql))
	if err != nil {
		fmt.Printf("error escaping SQL: %s", err)
		os.Exit(1)
	}

	payload := bytes.NewBuffer([]byte(
		`{ "type": "run_sql", "args": { "source": "`,
	))
	payload.WriteString(dataSourceName)
	payload.WriteString(`", "sql": `)
	payload.Write(escapedSql)
	payload.WriteString(` }}`)

	//fmt.Println(payload.String())
	return payload.Bytes()
}
