package commands

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/tidwall/sjson"
)

const HASURA_METADATA_API_EXPORT_PAYLOAD string = `{ "type": "export_metadata", "version": 2, "args": {} }`
const HASURA_METADATA_API_RELOAD_PAYLOAD string = `{ "type": "reload_metadata", "args": { "reload_remote_schemas": true, "reload_sources": true, "recreate_event_triggers": true }}`
const HASURA_METADATA_API_CLEAR_PAYLOAD string = `{ "type": "clear_metadata", "version": 2, "args": {}}`

func init() {
	importConfigCommand.Flags().StringP("fileURI", "f", "", "Path to the metadata file to be imported")
	importConfigCommand.MarkFlagRequired("fileURI")
	importConfigCommand.Flags().StringP("url", "u", "", "Hasura GraphQL Engine base URL")
	importConfigCommand.MarkFlagRequired("url")
	importConfigCommand.Flags().StringP("hasuraAdminSecret", "S", "", "Hasura Admin Secret")
	importConfigCommand.MarkFlagRequired("hasuraAdminSecret")

	configCommand.AddCommand(importConfigCommand)

	exportConfigCommand.Flags().StringP("url", "u", "", "Hasura GraphQL Engine base URL")
	exportConfigCommand.MarkFlagRequired("url")
	exportConfigCommand.Flags().StringP("hasuraAdminSecret", "S", "", "Hasura Admin Secret")
	exportConfigCommand.MarkFlagRequired("hasuraAdminSecret")
	exportConfigCommand.Flags().StringP("fileURI", "f", "", "Path to the exported metadata file that is exported")
	exportConfigCommand.MarkFlagRequired("fileURI")
	exportConfigCommand.Flags().StringP("gitCommitMessage", "m", "", "Description for Git commit")
	configCommand.AddCommand(exportConfigCommand)

	reloadMetadataCommand.Flags().StringP("url", "u", "", "Hasura GraphQL Engine base URL")
	reloadMetadataCommand.MarkFlagRequired("url")
	reloadMetadataCommand.Flags().StringP("hasuraAdminSecret", "S", "", "Hasura Admin Secret")
	reloadMetadataCommand.MarkFlagRequired("hasuraAdminSecret")

	configCommand.AddCommand(reloadMetadataCommand)
	rootCommand.AddCommand(configCommand)
}

func ResetCLIConfigFlags() {
	importConfigCommand.Flags().Lookup("fileURI").Value.Set("")
	importConfigCommand.Flags().Lookup("fileURI").Changed = false
	importConfigCommand.Flags().Lookup("url").Value.Set("")
	importConfigCommand.Flags().Lookup("url").Changed = false
	importConfigCommand.Flags().Lookup("hasuraAdminSecret").Value.Set("")
	importConfigCommand.Flags().Lookup("hasuraAdminSecret").Changed = false

	exportConfigCommand.Flags().Lookup("fileURI").Value.Set("")
	exportConfigCommand.Flags().Lookup("fileURI").Changed = false
	exportConfigCommand.Flags().Lookup("url").Value.Set("")
	exportConfigCommand.Flags().Lookup("url").Changed = false
	exportConfigCommand.Flags().Lookup("hasuraAdminSecret").Value.Set("")
	exportConfigCommand.Flags().Lookup("hasuraAdminSecret").Changed = false
	exportConfigCommand.Flags().Lookup("gitCommitMessage").Value.Set("")
	exportConfigCommand.Flags().Lookup("gitCommitMessage").Changed = false

	reloadMetadataCommand.Flags().Lookup("url").Value.Set("")
	reloadMetadataCommand.Flags().Lookup("url").Changed = false
	reloadMetadataCommand.Flags().Lookup("hasuraAdminSecret").Value.Set("")
	reloadMetadataCommand.Flags().Lookup("hasuraAdminSecret").Changed = false
}

var configCommand = &cobra.Command{
	Use:   "config",
	Short: "Metadata configuration related commands",
}

var importConfigCommand = &cobra.Command{
	Use:   "import",
	Short: "Import metadata config from a file into a Hasura GraphQL Engine",
	RunE: func(command *cobra.Command, args []string) error {
		fileURI, _ := command.Flags().GetString("fileURI")
		url, _ := command.Flags().GetString("url")
		hasuraAdminSecret, _ := command.Flags().GetString("hasuraAdminSecret")
		gitOptions := initGitOptionsFromCommand(command)

		return importMetdata(fileURI, gitOptions, url, hasuraAdminSecret)
	},
}

var exportConfigCommand = &cobra.Command{
	Use:   "export",
	Short: "Export metadata config from a Hasura GraphQL Engine into a file",
	RunE: func(command *cobra.Command, args []string) error {
		fileURI, _ := command.Flags().GetString("fileURI")
		url, _ := command.Flags().GetString("url")
		hasuraAdminSecret, _ := command.Flags().GetString("hasuraAdminSecret")
		gitCommitMessage, _ := command.Flags().GetString("gitCommitMessage")
		gitOptions := initGitOptionsFromCommand(command)

		return exportMetadata(url, hasuraAdminSecret, fileURI, gitOptions, gitCommitMessage)
	},
}

var reloadMetadataCommand = &cobra.Command{
	Use:   "reload",
	Short: "Reload existing metadata in aHasura GraphQL Engine",
	RunE: func(command *cobra.Command, args []string) error {
		url, _ := command.Flags().GetString("url")
		hasuraAdminSecret, _ := command.Flags().GetString("hasuraAdminSecret")

		return reloadMetdata(url, hasuraAdminSecret, true, true, true)
	},
}

func importMetdata(fileURI string, gitOptions *gitOptions, url, hasuraAdminSecret string) error {
	fmt.Printf("Importing metadata to HGE at %s\n", url)

	if gitOptions != nil && len(gitOptions.Uri) > 0 {
		fmt.Printf("  1:2 Reading file \"%s\" from %s ", fileURI, gitOptions.Uri)
	} else {
		fmt.Printf("  1:2 Reading file \"%s\" ", fileURI)
	}
	configMetadata := &bytes.Buffer{}
	fileURI, err := readFile(fileURI, "", "", gitOptions, configMetadata)
	if err != nil {
		return fmt.Errorf("error reading file: %s %w", fileURI, err)
	}
	fmt.Println("\u2714")

	payload := createMetadataPayload(configMetadata.Bytes())

	fmt.Printf("  2:2 Calling metadata API ")
	_, _, err = callHasuraAPI(url, hasuraAdminSecret, payload)
	if err != nil {
		return fmt.Errorf("\nerror importing metadata, %w", err)
	}
	fmt.Println("\u2714")

	return nil
}

func exportMetadata(url, hasuraAdminSecret, fileURI string, gitOptions *gitOptions, gitCommitMessage string) error {
	fmt.Printf("Exporting metadata from HGE at %s\n", url)
	payload := []byte(HASURA_METADATA_API_EXPORT_PAYLOAD)

	fmt.Printf("  1:2 Calling metadata API ")
	_, httpResponseBody, err := callHasuraAPI(url, hasuraAdminSecret, payload)
	if err != nil {
		return fmt.Errorf("\nerror exporting metadata, %w", err)
	}
	fmt.Println("\u2714")

	if gitOptions != nil && len(gitOptions.Uri) > 0 {
		fmt.Printf("  2:2 Writing file \"%s\" to %s ", fileURI, gitOptions.Uri)
	} else {
		fmt.Printf("  2:2 Writing file \"%s\" ", fileURI)
	}
	err = writeFile(httpResponseBody, fileURI, gitOptions, gitCommitMessage)
	if err != nil {
		return fmt.Errorf("error writing output file: %w", err)
	}
	fmt.Println("\u2714")

	return nil
}

func reloadMetdata(url, hasuraAdminSecret string, reloadSources, reloadRemoteSchemas, recreateEventTriggers bool) error {
	fmt.Printf("Reloading metadata for HGE at %s\n", url)
	payload := []byte(HASURA_METADATA_API_RELOAD_PAYLOAD)

	payload, _ = sjson.SetBytes(payload, "args.reload_sources", reloadSources)
	payload, _ = sjson.SetBytes(payload, "args.reload_remote_schemas", reloadRemoteSchemas)
	payload, _ = sjson.SetBytes(payload, "args.recreate_event_triggers", recreateEventTriggers)

	fmt.Printf("  1:1 Calling metadata API ")
	_, _, err := callHasuraAPI(url, hasuraAdminSecret, payload)
	if err != nil {
		return fmt.Errorf("\nerror reloading metadata, %w", err)
	}
	fmt.Println("\u2714")

	return nil
}

func callHasuraAPI(url, hasuraAdminSecret string, payload []byte) (int, io.ReadCloser, error) {
	httpRequest, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return -1, nil, fmt.Errorf("error creating HTTP request: %w", err)
	}
	httpRequest.Header.Add("content-type", "application/json")
	httpRequest.Header.Add("x-hasura-admin-secret", hasuraAdminSecret)

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return -1, nil, fmt.Errorf("error calling %s\n%w", url, err)
	}
	if httpResponse.StatusCode < 200 || httpResponse.StatusCode > 299 {
		responseBody, _ := io.ReadAll(httpResponse.Body)
		return httpResponse.StatusCode, nil, fmt.Errorf("error calling %s, status code: %d\n%s", url, httpResponse.StatusCode, responseBody)
	}

	return httpResponse.StatusCode, httpResponse.Body, nil
}

func createMetadataPayload(configMetadata []byte) []byte {
	// determine the position of the metadata object with the exported JSON
	metadataPosition := bytes.Index(configMetadata, []byte("metadata"))
	var payload *bytes.Buffer
	if metadataPosition == -1 {
		payload = bytes.NewBuffer([]byte(
			`{ "type": "replace_metadata", "version": 2, "args": { "allow_inconsistent_metadata": true, "metadata": `,
		))
		payload.Write(configMetadata)
		payload.WriteString("}}")
	} else {
		metadataObjectOpenCurlyBracePosition := bytes.Index(configMetadata[metadataPosition:], []byte("{"))
		metadataObjectOpenCurlyBracePosition = metadataPosition + metadataObjectOpenCurlyBracePosition
		lastCloseCurlyBracePosition := bytes.LastIndex(configMetadata[metadataObjectOpenCurlyBracePosition:], []byte("}"))
		lastCloseCurlyBracePosition = metadataObjectOpenCurlyBracePosition + lastCloseCurlyBracePosition

		payload = bytes.NewBuffer([]byte(
			`{ "type": "replace_metadata", "version": 2, "args": { "allow_inconsistent_metadata": true, "metadata": `,
		))
		payload.Write(configMetadata[metadataObjectOpenCurlyBracePosition:lastCloseCurlyBracePosition])
		payload.WriteString("}}")
	}

	return payload.Bytes()
}
