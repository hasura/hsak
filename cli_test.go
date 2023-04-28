package main

import (
	"os"
	"test/commands"
	"testing"

	"github.com/stretchr/testify/assert"
)

const URL = "http://localhost:8050"
const ADMIN_SECRET = "myadminsecretkey"
const DATA_SOURCE = "test"
const GIT_REPO_URL = "https://github.com/chris-hasura/test-metadata.git"
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
			URL,
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
			"https://raw.githubusercontent.com/hasura/chinook-demo/main/postgres/data-init/music.sql",
			"-u",
			URL,
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
			"postgres/data-init/music.sql",
			"-u",
			URL,
			"-S",
			ADMIN_SECRET,
			"-s",
			DATA_SOURCE,
			"--gitRepoURI",
			"https://github.com/hasura/chinook-demo.git",
			"--gitRepoBranch",
			"main",
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
			URL,
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
			"https://raw.githubusercontent.com/chris-hasura/test-metadata/main/import/metadata.json",
			"-u",
			URL,
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
			"import/metadata.json",
			"-u",
			URL,
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
			URL,
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
			"export/metadata.json",
			"-u",
			URL,
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
			URL,
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
			"examples/demo/chinook-demo-local.yaml",
		},
	)
	assert.Equal(t, err, nil)
}

func Test_demo_init_local_uri_web(t *testing.T) {
	err := testCLI(
		[]string{
			"demo",
			"init",
			"-f",
			"examples/demo/chinook-demo-local-web.yaml",
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
			"examples/demo/chinook-demo-local-git.yaml",
		},
	)
	assert.Equal(t, err, nil)
}
