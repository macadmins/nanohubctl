package main

import (
	"context"
	"log"

	"github.com/macadmins/nanohubctl/internal/cli"
)

func main() {
	ctx := context.Background()
	err := cli.ExecuteWithContext(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
