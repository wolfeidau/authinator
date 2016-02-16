package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	r "github.com/dancannon/gorethink"
	"github.com/emicklei/go-restful"
	"github.com/spf13/cobra"
	"github.com/wolfeidau/authinator/api"
	"github.com/wolfeidau/authinator/auth"
	"github.com/wolfeidau/authinator/store/users"
)

var (
	cmdServe = &cobra.Command{
		Use:   "serve",
		Short: "Start the authinator server",
		Long:  ``,
		Run:   runCmdServe,
	}

	serveOpts struct {
		ConnectionAddr string
	}
)

func init() {
	cmdServe.PersistentFlags().StringVar(&serveOpts.ConnectionAddr, "connection-addr", "localhost:28015", "Configure a connection address")
	cmdRoot.AddCommand(cmdServe)
}

func runCmdServe(cmd *cobra.Command, args []string) {

	session, err := r.Connect(r.ConnectOpts{
		Address: serveOpts.ConnectionAddr,
	})

	if err != nil {
		fmt.Printf("Opening RethinkDB session failed: %s", err)
		os.Exit(1)
	}

	userStore := users.NewUserStoreRethinkDB(session)

	wsContainer := restful.NewContainer()

	certs, err := auth.GenerateTestCerts()
	if err != nil {
		fmt.Printf("Opening RethinkDB session failed: %s", err)
		os.Exit(1)
	}

	jwtAuth := api.BuildJWTAuthFunc(userStore, certs)

	ar := api.NewAuthResource(userStore, jwtAuth, certs)

	ar.Register(wsContainer)

	ur := api.NewUserResource(userStore, jwtAuth)

	ur.Register(wsContainer)

	log.Printf("start listening on localhost:9090")
	server := &http.Server{Addr: ":9090", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
