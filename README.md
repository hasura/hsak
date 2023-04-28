# hsak
Hasura Swiss Army Knife CLI for SEs

hsak Commands:
- config: Supports JSON metadata from either the Hasura CLI or the Hasura Console
  - import: Import HGE metadata from a local file, web server or Git repository
  - export: Export HGE metadata to a local file or Git repository
  - reload: Reload HGE metadata
- sql: Execute a SQL script on a HGE datasource, from a local file, web server or Git repository
- demo
  - init: Initialize demo data & HGE server(s)

## Usage
### Environment File
hsak looks for the `config.env` file in the current working directory or the directory where the executable file is located.

This behaviour facilitates double-clicking on the executable and hsak can be run with preset arguments using the `HASURA_DEMO_ARGS` environment variable.

While the Git username & password/PAT can be specified on the command line. It is a best practice to set these variables using environment variables (`HASURA_GIT_USERNAME` & `HASURA_GIT_PWD_OR_PAT` respectively) or the config.env file.

```config.env
HASURA_DEMO_ARGS="demo init -f ./test/chinook-ct-cloud-demo.yaml -v"
HASURA_GIT_USERNAME=<please replace with username>
HASURA_GIT_PWD_OR_PAT=<please replace with password or person access token>
```

### Example Usage
Execute SQL a local script file on HGE data source
```
hsak sql -f ./test/chinook-music.sql -u http://localhost:8050 -S myadminsecretkey -s test
```

Execute SQL a web server script on HGE data source
```
hsak sql -f https://raw.githubusercontent.com/hasura/chinook-demo/main/postgres/data-init/music.sql -u http://localhost:8050 -S myadminsecretkey -s test
```

Execute SQL a script from a Git repository on HGE data source
```
hsak sql -f ./data-init/music.sql -u http://localhost:8050 -S myadminsecretkey -s test --gitRepoURI GIT_REPO_URL --gitRepoBranch GIT_REPO_BRANCH
```

Import HGE metadata from a local file
```
hsak config import -f ./test/hasura-metadata.json -u http://localhost:8050 -S myadminsecretkey
```

Import HGE metadata from a web server file
```
hsak config import -f https://raw.githubusercontent.com/hasura/chinook-demo/main/postgres/metadata/hasura-metadata.json -u http://localhost:8050 -S myadminsecretkey
```

Import HGE metadata from a file in a Git repository
```
hsak config import -f metadata/hasura-metadata.json -u http://localhost:8050 -S myadminsecretkey --gitRepoURI GIT_REPO_URL --gitRepoBranch GIT_REPO_BRANCH
```

Export HGE metadata to a local file
```
hsak config export -f ./temp/hasura-metadata-export.json -u http://localhost:8050 -S myadminsecretkey
```

Export HGE metadata to a Git repository
```
hsak config export -f test/metadata-export.json -u http://localhost:8050 -S myadminsecretkey --gitRepoURI GIT_REPO_URL --gitRepoBranch GIT_REPO_BRANCH --gitCommitMessage "test commit"
```

Reload HGE metadata
```
hsak config reload -u http://localhost:8050 -S myadminsecretkey
```

Initialize a demo using a local config file that references a demo from a web server
```
hsak demo init -f ../test/chinook-ct-cloud-demo.yaml
```

Initialize a demo using a local config file that references a demo in a Git repository
```
hsak demo init -f ../test/chinook-ct-cloud-demo-git.yaml
```
