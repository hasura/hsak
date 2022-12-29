package commands

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"gopkg.in/yaml.v3"
)

const HASURA_METADATA_API_PATH string = "/v1/metadata"
const HASURA_QUERY_API_PATH string = "/v2/query"
const HASURA_DATA_SOURCE string = `{ "name": "", "kind": "", "configuration": {}, "customization": {}, "tables": [] }`
const HASURA_METADATA_API_REPLACE_METADATA_PAYLOAD string = `{ "type": "replace_metadata", "version": 2, "args": { "allow_inconsistent_metadata": true, "metadata": { "version": 3, "sources": [] } } }`

func init() {
	initDemoCommand.Flags().StringP("configFileURI", "f", "", "URI of demo config file to be imported")
	initDemoCommand.MarkFlagRequired("configFileURI")

	setFlagFromEnv(initDemoCommand, "configFileURI", "HASURA_DEMO_CONFIG_FILE_URI")

	demoCommand.AddCommand(initDemoCommand)
	rootCommand.AddCommand(demoCommand)
}

func ResetCLIDemoFlags() {
	initDemoCommand.Flags().Lookup("configFileURI").Value.Set("")
	initDemoCommand.Flags().Lookup("configFileURI").Changed = false
}

var demoCommand = &cobra.Command{
	Use:   "demo",
	Short: "Demo related commands",
}

var initDemoCommand = &cobra.Command{
	Use:   "init",
	Short: "Initialize the demo HGE & database",
	RunE: func(command *cobra.Command, args []string) error {
		configFileURI, _ := command.Flags().GetString("configFileURI")
		gitOptions := initGitOptionsFromCommand(command)

		return initEngineAndDataSources(configFileURI, gitOptions)
	},
}

type demoConfig struct {
	Name          string
	Desc          string
	ConfigRootURI string      `yaml:"configRootURI"`
	ConfigGitRepo *gitOptions `yaml:"configGitRepo"`
	Engines       []hgeConfig
}

type hgeConfig struct {
	Name            string
	Desc            string
	Metadata        metadata
	HgeURL          string     `yaml:"hgeURL"`
	HgeAdminSecret  string     `yaml:"hgeAdminSecret"`
	DataSourceInits []dataInit `yaml:"dataSourceInits"`
	DependsOn       []string   `yaml:"dependsOn"`
}

type metadata struct {
	Uri       string
	Overrides []override
}

type dataInit struct {
	MetadataDataSourceName string   `yaml:"metadataDataSourceName"`
	FileURIs               []string `yaml:"fileURIs"`
}

func initEngineAndDataSources(configFileURI string, gitOptions *gitOptions) error {
	execPath, _ := os.Executable()
	configBytes := &bytes.Buffer{}
	configFileURI, err := readFile(configFileURI, execPath, "", gitOptions, configBytes)
	if err != nil {
		return fmt.Errorf("error reading file: %s %w", configFileURI, err)
	}
	if verbose {
		fmt.Printf("Initializing HGE engine(s) using config file at: %s\n", configFileURI)
	} else {
		fmt.Println("Initializing HGE engine(s)")
	}

	config := demoConfig{}
	err = yaml.Unmarshal(configBytes.Bytes(), &config)
	if err != nil {
		return fmt.Errorf("error marshalling YAML: %w", err)
	}

	if config.ConfigGitRepo != nil && len(config.ConfigGitRepo.Username) == 0 && len(config.ConfigGitRepo.PasswordOrPAT) == 0 &&
		gitOptions != nil && len(gitOptions.Username) > 0 && len(gitOptions.PasswordOrPAT) > 0 {
		config.ConfigGitRepo.Username = gitOptions.Username
		config.ConfigGitRepo.PasswordOrPAT = gitOptions.PasswordOrPAT
	}

	for engineIndex, engine := range config.Engines {
		engineURL := engine.HgeURL
		hgeMetadataBuffer := &bytes.Buffer{}
		metadataURI, err := readFile(engine.Metadata.Uri, configFileURI, config.ConfigRootURI, config.ConfigGitRepo, hgeMetadataBuffer)
		if err != nil {
			return fmt.Errorf("error reading file: %s %s", configFileURI, err)
		}
		if verbose {
			fmt.Printf("  %d:%d Configuring HGE \"%s\" at %s with metadata from %s\n", engineIndex+1, len(config.Engines), engine.Name, engineURL, metadataURI)
		} else {
			fmt.Printf("  %d:%d Configuring HGE \"%s\" at %s\n", engineIndex+1, len(config.Engines), engine.Name, engineURL)
		}

		hgeMetadataBytes, err := overrideJSONBytes(hgeMetadataBuffer.Bytes(), engine.Metadata.Overrides)
		if err != nil {
			return fmt.Errorf("error applying metadata overrides: %w", err)
		}
		//fmt.Println("hgeMetadataBytes:", string(hgeMetadataBytes))

		if len(engine.DataSourceInits) > 0 {
			fmt.Printf("    1:3 Loading HGE data sources connection info ")
			sourcesPath := "sources"
			if bytes.Index(hgeMetadataBytes, []byte("metadata")) > 0 {
				sourcesPath = "metadata.sources"
			}
			result := gjson.GetBytes(hgeMetadataBytes, sourcesPath)
			payload := HASURA_METADATA_API_REPLACE_METADATA_PAYLOAD
			result.ForEach(func(key, source gjson.Result) bool {
				dataSource := HASURA_DATA_SOURCE
				dataSource, _ = sjson.Set(dataSource, "name", source.Get("name").Str)
				dataSource, _ = sjson.Set(dataSource, "kind", source.Get("kind").Str)
				dataSource, _ = sjson.SetRaw(dataSource, "configuration", source.Get("configuration").Raw)
				payload, _ = sjson.SetRaw(payload, "args.metadata.sources.-1", dataSource)
				return true
			})

			_, _, err = callHasuraAPI(engineURL+HASURA_METADATA_API_PATH, engine.HgeAdminSecret, []byte(payload))
			if err != nil {
				return fmt.Errorf("\nerror importing metadata, %s", err)
			}
			fmt.Println("\u2714")

			fmt.Println("    2:3 Initializing data source(s)")
			for dataSourceIndex, dataSourceInit := range engine.DataSourceInits {
				for _, dataInitFileURI := range dataSourceInit.FileURIs {
					initBytes := &bytes.Buffer{}
					dataSourceInitURI, err := readFile(dataInitFileURI, configFileURI, config.ConfigRootURI, config.ConfigGitRepo, initBytes)
					if err != nil {
						return fmt.Errorf("error reading file: %s %w", configFileURI, err)
					}
					if verbose {
						fmt.Printf("      %d:%d Initializing data source \"%s\" using %s ", dataSourceIndex+1, len(engine.DataSourceInits), dataSourceInit.MetadataDataSourceName, dataSourceInitURI)
					} else {
						fmt.Printf("      %d:%d Initializing data source \"%s\" ", dataSourceIndex+1, len(engine.DataSourceInits), dataSourceInit.MetadataDataSourceName)
					}

					_, _, err = callHasuraAPI(
						engineURL+HASURA_QUERY_API_PATH,
						engine.HgeAdminSecret,
						createQueryAPIPayload(dataSourceInit.MetadataDataSourceName, initBytes.Bytes()),
					)
					if err != nil {
						return fmt.Errorf("\nerror importing metadata, %w", err)
					}
					fmt.Println("\u2714")
				}
			}
		}

		fmt.Printf("    3:3 Loading remaining metadata ")
		_, _, err = callHasuraAPI(
			engineURL+HASURA_METADATA_API_PATH,
			engine.HgeAdminSecret,
			createMetadataPayload(hgeMetadataBytes),
		)
		if err != nil {
			return fmt.Errorf("\nerror importing metadata, %w", err)
		}
		fmt.Println("\u2714")
	}

	return nil
}
