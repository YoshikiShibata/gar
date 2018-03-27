// Copyright © 2018 Yoshiki Shibata. All rights reserved.

package main

import (
	"fmt"
	"os"
)

// gar (Github Assets Retriever) retrieve assets from the Github.com
// gar uses a personal access token for authentication: The personal access
// token must be presented with an enironment variable named GITHUB_TOKEN.
//
// gar owner repository tag
func main() {
	if len(os.Args) != 4 {
		showUsageAndExit()
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "GITHUB_TOKEN is not set")
	}

	err := retrieveAssets(os.Args[1], os.Args[2], os.Args[3], token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}

func showUsageAndExit() {
	fmt.Fprintln(os.Stderr, "usage: gar owner repository tag\n")
	os.Exit(1)
}
