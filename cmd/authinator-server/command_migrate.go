package main

import (
	"fmt"
	"os"

	r "github.com/dancannon/gorethink"
	"github.com/spf13/cobra"
	"github.com/wolfeidau/authinator/store/users"
)

var (
	cmdMigrate = &cobra.Command{
		Use:   "migrate",
		Short: "Perform a schema migration on the RethinkDB database",
		Long:  ``,
		Run:   runCmdMigrate,
	}

	migrateOpts struct {
		ConnectionAddr string
	}
)

func init() {
	cmdMigrate.PersistentFlags().StringVar(&migrateOpts.ConnectionAddr, "connection-addr", "localhost:28015", "Configure a connection address")
	cmdRoot.AddCommand(cmdMigrate)

}

func runCmdMigrate(cmd *cobra.Command, args []string) {

	session, err := r.Connect(r.ConnectOpts{
		Address: migrateOpts.ConnectionAddr,
	})

	if err != nil {
		fmt.Printf("Opening RethinkDB session failed: %s", err)
		os.Exit(1)
	}

	resp, err := r.DBCreate(users.DBName).RunWrite(session)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%d DB created\n", resp.DBsCreated)

	r.DB(users.DBName).TableCreate(users.TableName).Exec(session)

	fmt.Printf("Table created\n")
}
