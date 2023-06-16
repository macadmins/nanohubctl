package main

import (
	"context"
	"log"

	"github.com/macadmins/ddmctl/internal/cli"
)

func main() {
	ctx := context.Background()
	err := cli.ExecuteWithContext(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
