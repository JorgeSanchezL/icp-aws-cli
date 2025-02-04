package main

import (
	"fmt"
	commands "icp-aws-cli/cmd/icp-aws-cli"
	"icp-aws-cli/pkg/awsclient"
	"os"
)

func main() {
	clients, err := awsclient.NewAWSClientCollection()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initializing AWS clients: %v\n", err)
		os.Exit(1)
	}

	commands.InitCommands(clients)
	commands.Execute()
}
