package main

import (
	"context"
	"log"

	"github.com/altinn/dotnet-monitor-sidecar-cli/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	ctx := context.Background()
	dmsctl := cmd.NewDmsctlCommand(ctx)
	err := doc.GenMarkdownTree(dmsctl, "./docs")
	if err != nil {
		log.Fatal(err)
	}
}
