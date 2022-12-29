package commands

import (
	"errors"
	"fmt"
	"io"

	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	gitHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/joho/godotenv"
	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var verbose bool = false

var rootCommand = &cobra.Command{
	Use:   "hsak",
	Short: "Hasura SE Swiss Army Knife",
	Long:  `A CLI to manage demo build/maintenance tasks`,
}

func init() {
	rootCommand.PersistentFlags().BoolP("verbose", "v", false, "Output more processing info during execution")
	rootCommand.PersistentFlags().String("gitRepoURI", "", "Git repository URI, which typically ends with .git (Env: HASURA_GIT_REPO_URI)")
	rootCommand.PersistentFlags().String("gitRepoBranch", "", "Git repository branch (Env: HASURA_GIT_REPO_BRANCH)")
	rootCommand.PersistentFlags().String("gitUsername", "", "Git username for authn & commits (Env: HASURA_GIT_USERNAME)")
	rootCommand.PersistentFlags().String("gitPasswordOrPAT", "", "Git password or personal access token (PAT) for authn & commits (Env: HASURA_GIT_PWD_OR_PAT)")

	verbose = slices.Contains(os.Args, "-v") || slices.Contains(os.Args, "--verbose")

	execPath, _ := os.Executable()
	//fmt.Println("execPath:", execPath)
	loadedEnvFile := false
	err := godotenv.Load("config.env")
	if err != nil {
		err = godotenv.Load(path.Join(path.Dir(execPath), "config.env"))
		if err == nil {
			loadedEnvFile = true
		}
	} else {
		loadedEnvFile = true
	}

	if len(os.Args[1:]) == 0 { // if there are no CLI args then try to set from HASURA_DEMO_ARGS
		argsEnv := os.Getenv("HASURA_DEMO_ARGS")
		if len(argsEnv) > 0 {
			args := strings.Split(argsEnv, " ")
			verbose = slices.Contains(args, "-v")
			if verbose {
				if loadedEnvFile {
					fmt.Println("Loaded environment variable file")
				}
				fmt.Println("Setting args from env")
			}
			rootCommand.SetArgs(args)
		}
	}

	setPersistentFlagFromEnv(rootCommand, "gitRepoURI", "HASURA_GIT_REPO_URI")
	setPersistentFlagFromEnv(rootCommand, "gitRepoBranch", "HASURA_GIT_REPO_BRANCH")
	setPersistentFlagFromEnv(rootCommand, "gitUsername", "HASURA_GIT_USERNAME")
	setPersistentFlagFromEnv(rootCommand, "gitPasswordOrPAT", "HASURA_GIT_PWD_OR_PAT")
}

func ResetCLIFlags() {
	rootCommand.PersistentFlags().Lookup("verbose").Value.Set("")
	rootCommand.PersistentFlags().Lookup("verbose").Changed = false
	rootCommand.PersistentFlags().Lookup("gitRepoURI").Value.Set("")
	rootCommand.PersistentFlags().Lookup("gitRepoURI").Changed = false
	rootCommand.PersistentFlags().Lookup("gitRepoBranch").Value.Set("")
	rootCommand.PersistentFlags().Lookup("gitRepoBranch").Changed = false
	rootCommand.PersistentFlags().Lookup("gitUsername").Value.Set("")
	rootCommand.PersistentFlags().Lookup("gitUsername").Changed = false
	rootCommand.PersistentFlags().Lookup("gitPasswordOrPAT").Value.Set("")
	rootCommand.PersistentFlags().Lookup("gitPasswordOrPAT").Changed = false

	setPersistentFlagFromEnv(rootCommand, "gitRepoURI", "HASURA_GIT_REPO_URI")
	setPersistentFlagFromEnv(rootCommand, "gitRepoBranch", "HASURA_GIT_REPO_BRANCH")
	setPersistentFlagFromEnv(rootCommand, "gitUsername", "HASURA_GIT_USERNAME")
	setPersistentFlagFromEnv(rootCommand, "gitPasswordOrPAT", "HASURA_GIT_PWD_OR_PAT")
}

func Execute() error {
	return rootCommand.Execute()
}

func setFlagFromEnv(command *cobra.Command, flagName, envName string) {
	value := os.Getenv(envName)
	if len(value) > 0 {
		command.Flags().Set(flagName, value)
	}
}

func setPersistentFlagFromEnv(command *cobra.Command, flagName, envName string) {
	value := os.Getenv(envName)
	if len(value) > 0 {
		command.PersistentFlags().Set(flagName, value)
	}
}

func initGitOptionsFromCommand(command *cobra.Command) *gitOptions {
	gitOptions := &gitOptions{}

	gitOptions.Uri, _ = command.Flags().GetString("gitRepoURI")
	gitOptions.Branch, _ = command.Flags().GetString("gitRepoBranch")
	gitOptions.Username, _ = command.Flags().GetString("gitUsername")
	gitOptions.PasswordOrPAT, _ = command.Flags().GetString("gitPasswordOrPAT")

	return gitOptions
}

func fileExists(fileURI string) (bool, error) {
	if strings.HasPrefix(fileURI, "http") {
		response, err := http.Head(fileURI)
		if err != nil {
			return false, fmt.Errorf("error accessing fileURI \"%s\": %w", fileURI, err)
		}
		if response.StatusCode >= 200 && response.StatusCode <= 299 {
			return true, nil
		} else if response.StatusCode == 404 {
			return false, nil
		} else {
			return false, fmt.Errorf("error accessing fileURI \"%s\": HTTP status %d, %s", fileURI, response.StatusCode, response.Status)
		}
	} else {
		if _, err := os.Stat(fileURI); err == nil {
			return true, nil
		} else if errors.Is(err, os.ErrNotExist) {
			return false, nil
		} else {
			return false, fmt.Errorf("error accessing file \"%s\": %w", fileURI, err)
		}
	}
}

func joinPath(path1, path2 string, isPath1Dir bool) string {
	if len(path1) == 0 {
		return path2
	} else if len(path2) == 0 {
		return path1
	}
	if strings.HasPrefix(path1, "http") {
		if isPath1Dir {
			relativeURI, _ := url.JoinPath(path1, path2)
			return relativeURI
		} else {
			relativeURI, _ := url.JoinPath(path1, "..", path2)
			return relativeURI
		}
	} else {
		if isPath1Dir {
			return path.Join(path1, path2)
		} else {
			return path.Join(path.Dir(path1), path2)
		}
	}
}

// repoUrl, branch, relativeFileURI, username, passwordOrPAT string, writer io.Writer
func readFile(fileURI, referenceURI, rootURI string, gitOptions *gitOptions, writer io.Writer) (string, error) {
	if gitOptions != nil && len(gitOptions.Uri) > 0 {
		return fileURI, gitRead(fileURI, rootURI, gitOptions, writer)
	} else {
		exists, err := fileExists(fileURI)
		if !exists {
			if len(rootURI) > 0 {
				relativeURI := joinPath(rootURI, fileURI, true)
				exists, err = fileExists(relativeURI)
				if err != nil {
					return fileURI, err
				} else if exists {
					fileURI = relativeURI
				} else {
					relativeURI = joinPath(referenceURI, rootURI, false)
					relativeURI = joinPath(relativeURI, fileURI, true)
					exists, err = fileExists(relativeURI)
					if err != nil {
						return fileURI, err
					} else if exists {
						fileURI = relativeURI
					} else {
						return fileURI, fmt.Errorf("error file does not exist: \"%s\"", fileURI)
					}
				}
			} else {
				relativeURI := joinPath(referenceURI, fileURI, false)
				exists, err := fileExists(relativeURI)
				if err != nil {
					return fileURI, err
				} else if exists {
					fileURI = relativeURI
				} else {
					return fileURI, fmt.Errorf("error file does not exist: \"%s\"", fileURI)
				}
			}
		}

		if !exists {
			return fileURI, err
		}
	}

	if strings.HasPrefix(fileURI, "http") {
		httpRequest, err := http.NewRequest(http.MethodGet, fileURI, nil)
		if err != nil {
			return fileURI, fmt.Errorf("error creating HTTP request: %w", err)
		}

		httpResponse, err := http.DefaultClient.Do(httpRequest)
		if err != nil {
			return fileURI, fmt.Errorf("error reading file: %w", err)
		}
		_, err = io.Copy(writer, httpResponse.Body)
		if err != nil {
			return fileURI, fmt.Errorf("error reading HTTP response: %w", err)
		}

		return fileURI, nil
	} else {
		file, err := os.Open(fileURI)
		if err != nil {
			return fileURI, fmt.Errorf("error opening file: %w", err)
		}
		_, err = io.Copy(writer, file)
		return fileURI, err
	}
}

func writeFile(fileData io.ReadCloser, fileURI string, gitOptions *gitOptions, gitCommitMessage string) error {
	if strings.HasPrefix(fileURI, "http") {
		return fmt.Errorf("HTTP write only supports Git repositories")
	} else if gitOptions != nil && len(gitOptions.Uri) > 0 {
		return gitWrite(fileData, fileURI, gitOptions, gitCommitMessage)
	} else {
		outputFile, err := os.Create(fileURI)
		if err != nil {
			fmt.Printf("error creating output file: %s\n", err)
			return err
		}
		defer outputFile.Close()

		_, err = io.Copy(outputFile, fileData)

		return err
	}
}

type gitOptions struct {
	Uri           string `yaml:"uri"`
	Branch        string `yaml:"branch"`
	Username      string `yaml:"username"`
	PasswordOrPAT string `yaml:"passwordOrPAT"`
}

func gitClone(repoUrl, branch string, auth *gitHttp.BasicAuth) (*git.Repository, *git.Worktree, error) {
	// Clones the given repository in memory, creating the remote, the local
	repo, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL:           repoUrl,
		Auth:          auth,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		SingleBranch:  true,
		Depth:         1,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error cloning repo: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, nil, fmt.Errorf("error getting worktree: %w", err)
	}

	return repo, worktree, nil
}

func gitRead(fileURI, rootURI string, gitOptions *gitOptions, writer io.Writer) error {
	if gitOptions == nil {
		return fmt.Errorf("gitRead error Git options are missing")
	}

	auth := &gitHttp.BasicAuth{
		Username: gitOptions.Username,
		Password: gitOptions.PasswordOrPAT,
	}

	_, worktree, err := gitClone(gitOptions.Uri, gitOptions.Branch, auth)
	if err != nil {
		return fmt.Errorf("gitRead: %w", err)
	}

	fileHandle, err := worktree.Filesystem.Open(fileURI)
	if err != nil {
		if len(rootURI) > 0 {
			relativeURI := joinPath(rootURI, fileURI, true)
			fileHandle, err = worktree.Filesystem.Open(relativeURI)
			if err != nil {
				return fmt.Errorf("error opening file: %w", err)
			}
		} else {
			return fmt.Errorf("error opening file: %w", err)
		}
	}

	_, err = io.Copy(writer, fileHandle)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

func gitWrite(fileData io.ReadCloser, fileURI string, gitOptions *gitOptions, gitCommitMessage string) error {
	if gitOptions == nil {
		return fmt.Errorf("gitWrite error Git options are missing")
	}

	auth := &gitHttp.BasicAuth{
		Username: gitOptions.Username,
		Password: gitOptions.PasswordOrPAT,
	}

	repo, worktree, err := gitClone(gitOptions.Uri, gitOptions.Branch, auth)
	if err != nil {
		return fmt.Errorf("gitRead: %w", err)
	}

	fileHandle, err := worktree.Filesystem.Create(fileURI)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}

	io.Copy(fileHandle, fileData)

	_, err = worktree.Add(".")
	if err != nil {
		return fmt.Errorf("add error: %w", err)
	}

	status, err := worktree.Status()
	if err != nil {
		return fmt.Errorf("status error: %w", err)
	}
	if verbose {
		fmt.Println(status)
	}

	commit, err := worktree.Commit(gitCommitMessage, &git.CommitOptions{})
	if err != nil {
		return fmt.Errorf("commit error: %w", err)
	}

	commitObj, err := repo.CommitObject(commit)
	if err != nil {
		return fmt.Errorf("error retriving commit info: %w", err)
	}
	if verbose {
		fmt.Println(commitObj)
	}

	err = repo.Push(&git.PushOptions{Auth: auth})
	if err != nil {
		return fmt.Errorf("push error: %w", err)
	}

	return nil
}

type override struct {
	Find        string
	ReplaceWith string `yaml:"replaceWith"`
}

/*
	type jsonParser struct {
		parser oj.Parser
	}
*/
func overrideJSON(json []byte, overrides []override) (string, error) {
	overriden, err := overrideJSONBytes([]byte(json), overrides)
	if err != nil {
		return "", err
	}
	return string(overriden), err
}
func overrideJSONBytes(json []byte, overrides []override) ([]byte, error) {
	if len(overrides) == 0 {
		return json, nil
	}

	parsedJson, err := oj.Parse(json)
	if err != nil {
		return nil, err
	}
	for _, override := range overrides {
		query, err := jp.ParseString(override.Find)
		if err != nil {
			return nil, err
		}
		//value := query.Get(parsedJson)
		query.Set(parsedJson, override.ReplaceWith)
	}

	return oj.Marshal(parsedJson, ojg.Options{Sort: true})
}
