package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/seed"
)

func main() {
	var op int
	var n int

	flag.IntVar(&op, "op", 0, "operation to perform (1: add users)")
	flag.IntVar(&n, "n", 0, "number of records")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	seed, db, err := seed.New(logger)
	if err != nil {
		slog.Error("failed to create seed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	switch op {
	case 0:
		slog.Error("no operation specified")
	case 1:
		if n == 0 {
			slog.Error("no number of users specified")
		} else {
			cnt, err := seed.AddRandomUsers(n)
			if err != nil {
				slog.Error("failed to add random users", slog.String("error", err.Error()))
			} else {
				slog.Info("added random users successfully", slog.Int("success", cnt))
			}
		}
	default:
		slog.Error("invalid operation")
	}
}
