package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/alek101/GoMikroservisChiNinja/application"
)

func main() {
	app, err := application.New()
	if err != nil {
		fmt.Println("failed to start app:", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err == nil {
		err = app.Start(ctx)
		if err != nil {
			fmt.Println("failed to start app:", err)
		}
	}

}
