package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_initEngineAndDataSources_local_uri(t *testing.T) {
	configFileURI := "examples/demo/chinook-demo-local.yaml"
	var gitOptions *gitOptions = nil

	err := initEngineAndDataSources(configFileURI, gitOptions)
	assert.Equal(t, err, nil)
}

func Test_initEngineAndDataSources_local_uri_web(t *testing.T) {
	configFileURI := "examples/demo/chinook-demo-local-web.yaml"
	var gitOptions *gitOptions = nil

	err := initEngineAndDataSources(configFileURI, gitOptions)
	assert.Equal(t, err, nil)
}

func Test_initEngineAndDataSources_local_uri_git(t *testing.T) {
	configFileURI := "examples/demo/chinook-demo-local-git.yaml"
	var gitOptions *gitOptions = nil

	err := initEngineAndDataSources(configFileURI, gitOptions)
	assert.Equal(t, err, nil)
}
