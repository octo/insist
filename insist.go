package main

import (
	"bytes"
	"context"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/octo/retry"
)

var (
	attempts = flag.Int("attempts", 0, "Number of attempts or zero to retry indefinitely.")
	timeout  = flag.Duration("timeout", 0, "Timeout of each individual attempt or zero to block indefinitely.")
)

func run(ctx context.Context, stdin io.Reader) error {
	log.Println("Attempt", retry.Attempt(ctx), "...")

	args := flag.Args()
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdin = stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if os.IsNotExist(err) {
		return retry.Abort(err)
	}

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

	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("reading STDIN:", err)
	}

	opts := []retry.Option{
		retry.Attempts(*attempts),
		retry.Timeout(*timeout),
	}

	err = retry.Do(ctx, func(ctx context.Context) error {
		return run(ctx, bytes.NewReader(in))
	}, opts...)
	if err != nil {
		log.Fatal(err)
	}
}
