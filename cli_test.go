package main

import (
	"os"
	"test/commands"
	"testing"

	"github.com/stretchr/testify/assert"
)

const QUERY_URL = "http://localhost:8050/v2/query"
const METADATA_URL = "http://localhost:8050/v1/metadata"
const ADMIN_SECRET = "myadminsecretkey"
const DATA_SOURCE = "test"
const GIT_REPO_URL = "https://github.com/chris-hasura/chinook-demo.git"
const GIT_REPO_BRANCH = "main"

func testCLI(args []string) error {
	testArgs := append(os.Args[0:1], args...)
	os.Args = testArgs

	commands.ResetCLIFlags()
	commands.ResetCLISqlFlags()
	commands.ResetCLIConfigFlags()
	commands.ResetCLIDemoFlags()

	return commands.Execute()
}

func Test_sql_local_uri(t *testing.T) {
	err := testCLI(
		[]string{
			"sql",
			"-f",
			"./test/chinook-music.sql",
			"-u",
			QUERY_URL,
			"-S",
			ADMIN_SECRET,
			"-s",
			DATA_SOURCE,
		},
	)
	assert.Equal(t, err, nil)
}

func Test_sql_web_uri(t *testing.T) {
	err := testCLI(
		[]string{
			"sql",
			"-f",
			"https://raw.githubusercontent.com/hasura/chinook-demo/main/data-init/music.sql",
			"-u",
			QUERY_URL,
			"-S",
			ADMIN_SECRET,
			"-s",
			DATA_SOURCE,
		},
	)
	assert.Equal(t, err, nil)
}

func Test_sql_git_uri(t *testing.T) {
	err := testCLI(
		[]string{
			"sql",
			"-f",
			"./data-init/music.sql",
			"-u",
			QUERY_URL,
			"-S",
			ADMIN_SECRET,
			"-s",
			DATA_SOURCE,
			"--gitRepoURI",
			GIT_REPO_URL,
			"--gitRepoBranch",
			GIT_REPO_BRANCH,
			//			"--gitUsername",
			//			GIT_USERNAME,
			//			"--gitPasswordOrPAT",
			//			GIT_PWD_PAT,
		},
	)
	assert.Equal(t, err, nil)
}

func Test_config_import_local_uri(t *testing.T) {
	err := testCLI(
		[]string{
			"config",
			"import",
			"-f",
			"./test/hasura-metadata.json",
			"-u",
			METADATA_URL,
			"-S",
			ADMIN_SECRET,
		},
	)
	assert.Equal(t, err, nil)
}

func Test_config_import_web_uri(t *testing.T) {
	err := testCLI(
		[]string{
			"config",
			"import",
			"-f",
			"https://raw.githubusercontent.com/chris-hasura/chinook-demo/main/metadata/hasura-metadata.json",
			"-u",
			METADATA_URL,
			"-S",
			ADMIN_SECRET,
		},
	)
	assert.Equal(t, err, nil)
}

func Test_config_import_git_uri(t *testing.T) {
	err := testCLI(
		[]string{
			"config",
			"import",
			"-f",
			"metadata/hasura-metadata.json",
			"-u",
			METADATA_URL,
			"-S",
			ADMIN_SECRET,
			"--gitRepoURI",
			GIT_REPO_URL,
			"--gitRepoBranch",
			GIT_REPO_BRANCH,
			//			"--gitUsername",
			//			GIT_USERNAME,
			//			"--gitPasswordOrPAT",
			//			GIT_PWD_PAT,
		},
	)
	assert.Equal(t, err, nil)
}

func Test_config_export_local_uri(t *testing.T) {
	err := testCLI(
		[]string{
			"config",
			"export",
			"-f",
			"./temp/hasura-metadata-export.json",
			"-u",
			METADATA_URL,
			"-S",
			ADMIN_SECRET,
		},
	)
	assert.Equal(t, err, nil)
}

func Test_config_export_git_uri(t *testing.T) {
	err := testCLI(
		[]string{
			"config",
			"export",
			"-f",
			"test/metadata-export.json",
			"-u",
			METADATA_URL,
			"-S",
			ADMIN_SECRET,
			"--gitRepoURI",
			GIT_REPO_URL,
			"--gitRepoBranch",
			GIT_REPO_BRANCH,
			//			"--gitUsername",
			//			GIT_USERNAME,
			//			"--gitPasswordOrPAT",
			//			GIT_PWD_PAT,
			"--gitCommitMessage",
			"test commit",
		},
	)
	assert.Equal(t, err, nil)
}

func Test_config_reload(t *testing.T) {
	err := testCLI(
		[]string{
			"config",
			"reload",
			"-u",
			METADATA_URL,
			"-S",
			ADMIN_SECRET,
		},
	)
	assert.Equal(t, err, nil)
}

func Test_demo_init_local_uri(t *testing.T) {
	err := testCLI(
		[]string{
			"demo",
			"init",
			"-f",
			"../test/chinook-ct-cloud-demo.yaml",
		},
	)
	assert.Equal(t, err, nil)
}

func Test_demo_init_local_uri_git(t *testing.T) {
	err := testCLI(
		[]string{
			"demo",
			"init",
			"-f",
			"../test/chinook-ct-cloud-demo-git.yaml",
		},
	)
	assert.Equal(t, err, nil)
}
