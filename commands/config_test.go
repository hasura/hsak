package commands

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_importMetdata_local_uri(t *testing.T) {
	fileURI := "test/hasura-metadata.json"
	url := "http://localhost:8050/v1/metadata"
	hasuraAdminSecret := "myadminsecretkey"
	var gitOptions *gitOptions = nil

	err := importMetdata(fileURI, gitOptions, url, hasuraAdminSecret)
	assert.Equal(t, err, nil)
}

func Test_importMetdata_web_uri(t *testing.T) {
	fileURI := "https://raw.githubusercontent.com/chris-hasura/test-metadata/main/import/metadata.json"
	url := "http://localhost:8050/v1/metadata"
	hasuraAdminSecret := "myadminsecretkey"
	var gitOptions *gitOptions = nil

	err := importMetdata(fileURI, gitOptions, url, hasuraAdminSecret)
	assert.Equal(t, err, nil)
}

func Test_importMetdata_git_uri(t *testing.T) {
	fileURI := "import/metadata.json"
	url := "http://localhost:8050/v1/metadata"
	hasuraAdminSecret := "myadminsecretkey"
	gitOptions := &gitOptions{
		Uri:           "https://github.com/chris-hasura/test-metadata.git",
		Branch:        "main",
		Username:      os.Getenv("HASURA_GIT_USERNAME"),
		PasswordOrPAT: os.Getenv("HASURA_GIT_PWD_OR_PAT"),
	}

	err := importMetdata(fileURI, gitOptions, url, hasuraAdminSecret)
	assert.Equal(t, err, nil)
}

func Test_exportMetadata_local_uri(t *testing.T) {
	fileURI := "temp/metadata-export.json"
	url := "http://localhost:8050/v1/metadata"
	hasuraAdminSecret := "myadminsecretkey"
	var gitOptions *gitOptions = nil
	gitCommitMessage := ""

	err := exportMetadata(url, hasuraAdminSecret, fileURI, gitOptions, gitCommitMessage)
	assert.Equal(t, err, nil)
}

func Test_exportMetadata_git_uri(t *testing.T) {
	fileURI := "export/metadata.json"
	url := "http://localhost:8050/v1/metadata"
	hasuraAdminSecret := "myadminsecretkey"
	gitOptions := &gitOptions{
		Uri:           "https://github.com/chris-hasura/test-metadata.git",
		Branch:        "main",
		Username:      os.Getenv("HASURA_GIT_USERNAME"),
		PasswordOrPAT: os.Getenv("HASURA_GIT_PWD_OR_PAT"),
	}
	gitCommitMessage := "test commit"

	err := exportMetadata(url, hasuraAdminSecret, fileURI, gitOptions, gitCommitMessage)
	assert.Equal(t, err, nil)
}

func Test_reloadMetdata(t *testing.T) {
	url := "http://localhost:8050/v1/metadata"
	hasuraAdminSecret := "myadminsecretkey"

	err := reloadMetdata(url, hasuraAdminSecret, true, true, true)
	assert.Equal(t, err, nil)
}
