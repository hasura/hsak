package commands

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_runSQL_local_uri(t *testing.T) {
	fileURI := "../test/chinook-music.sql"
	url := "http://localhost:8050"
	paths := &urlPaths{
		queryPath: HASURA_QUERY_API_PATH,
	}
	hasuraAdminSecret := "myadminsecretkey"
	dataSourceName := "test"
	var gitOptions *gitOptions = nil

	err := runSQL(fileURI, url, paths, hasuraAdminSecret, dataSourceName, gitOptions)
	assert.Equal(t, err, nil)
}

func Test_runSQL_web_uri(t *testing.T) {
	fileURI := "https://raw.githubusercontent.com/hasura/chinook-demo/main/postgres/data-init/music.sql"
	url := "http://localhost:8050"
	paths := &urlPaths{
		queryPath: HASURA_QUERY_API_PATH,
	}
	hasuraAdminSecret := "myadminsecretkey"
	dataSourceName := "test"
	var gitOptions *gitOptions = nil

	err := runSQL(fileURI, url, paths, hasuraAdminSecret, dataSourceName, gitOptions)
	assert.Equal(t, err, nil)
}

func Test_runSQL_git_uri(t *testing.T) {
	fileURI := "postgres/data-init/music.sql"
	url := "http://localhost:8050"
	paths := &urlPaths{
		queryPath: HASURA_QUERY_API_PATH,
	}
	hasuraAdminSecret := "myadminsecretkey"
	dataSourceName := "test"
	gitOptions := &gitOptions{
		Uri:           "https://github.com/hasura/chinook-demo.git",
		Branch:        "main",
		Username:      os.Getenv("HASURA_GIT_USERNAME"),
		PasswordOrPAT: os.Getenv("HASURA_GIT_PWD_OR_PAT"),
	}

	err := runSQL(fileURI, url, paths, hasuraAdminSecret, dataSourceName, gitOptions)
	assert.Equal(t, err, nil)
}
