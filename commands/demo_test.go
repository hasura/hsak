package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_initEngineAndDataSources_local_uri(t *testing.T) {
	configFileURI := "examples/demo/chinook-demo-local.yaml"
	paths := &urlPaths{
		metadataPath: HASURA_METADATA_API_PATH,
		queryPath:    HASURA_QUERY_API_PATH,
	}
	var gitOptions *gitOptions = nil

	err := initEngineAndDataSources(configFileURI, paths, gitOptions)
	assert.Equal(t, err, nil)
}

func Test_initEngineAndDataSources_local_uri_web(t *testing.T) {
	configFileURI := "examples/demo/chinook-demo-local-web.yaml"
	paths := &urlPaths{
		metadataPath: HASURA_METADATA_API_PATH,
		queryPath:    HASURA_QUERY_API_PATH,
	}
	var gitOptions *gitOptions = nil

	err := initEngineAndDataSources(configFileURI, paths, gitOptions)
	assert.Equal(t, err, nil)
}

func Test_initEngineAndDataSources_local_uri_git(t *testing.T) {
	configFileURI := "examples/demo/chinook-demo-local-git.yaml"
	paths := &urlPaths{
		metadataPath: HASURA_METADATA_API_PATH,
		queryPath:    HASURA_QUERY_API_PATH,
	}
	var gitOptions *gitOptions = nil

	err := initEngineAndDataSources(configFileURI, paths, gitOptions)
	assert.Equal(t, err, nil)
}
