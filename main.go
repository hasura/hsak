package main

import (
	"fmt"
	"os"
	"test/commands"
)

/*
const filepath string = "./hasura-metadata/exported-metadata.json"
const sqlpath string = "./sql/corp.sql"
const HASURA_METADATA_API_URL string = "http://localhost:8010/v1/metadata"
const HASURA_SCHEMA_API_URL string = "http://localhost:8010/v2/query"
const hasuraAdminSecret string = "ZtuSwFP24OFpsSdfqTFelNWXWLfXKo6OW5wjibBLF2llSl2gsZHz2I2a4FKNIwdh"
*/
func main() {
	commands.Execute()
	if err := commands.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
