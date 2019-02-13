package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"

	"github.com/octo/retry"
)

var (
	attempts = flag.Int("attempts", 0, "Number of attempts or zero to retry indefinitely.")
	timeout  = flag.Duration("timeout", 0, "Timeout of each individual attempt or zero to block indefinitely.")
)

func run(ctx context.Context) error {
	log.Println("Attempt", retry.Attempt(ctx), "...")

	args := flag.Args()
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	select {
	case <-ctx.Done():
		log.Println("CANCELLED:", ctx.Err())
	default:
		if err != nil {
			log.Println("FAILURE:", err)
		}
	}
	return err
}

func main() {
	ctx := context.Background()
	flag.Parse()

	opts := []retry.Option{
		retry.Attempts(*attempts),
		retry.Timeout(*timeout),
	}

	if err := retry.Do(ctx, run, opts...); err != nil {
		log.Fatal(err)
	}
}
