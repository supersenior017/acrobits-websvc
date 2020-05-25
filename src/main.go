package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	acrobitswebsvc "./lib"
)

const (
	FlagConfig  = "c"
	ExitSuccess = 0
)

type Opts struct {
	config acrobitswebsvc.Config
}

func parseArgs() *Opts {
	opts := &Opts{}
	isHelp := flag.Bool("h", false, "Print help and exit")
	isVersion := flag.Bool("V", false, "Print version and exit")
	flag.Var(&opts.config, FlagConfig, "Path to configuration file")
	flag.Parse()
	if *isHelp {
		flag.Usage()
		os.Exit(ExitSuccess)
	}
	if *isVersion {
		fmt.Println(acrobitswebsvc.Version)
		os.Exit(ExitSuccess)
	}
	if !opts.config.IsSet() {
		opts.config.SetDefaults()
	}
	opts.config.Balance.Func = getBalance
	return opts
}

func getBalance(
	ctx context.Context,
	params acrobitswebsvc.Params,
) (float64, error) {
	return 0, fmt.Errorf("NotImplemented")
}

func main() {
	opts := parseArgs()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		defer signal.Stop(sigint)
		<-sigint
		cancel()
	}()
	log.Println("Starting on:", opts.config.Addr)
	if err := acrobitswebsvc.ListenAndServe(
		ctx,
		opts.config,
	); err != nil && err != http.ErrServerClosed {
		log.Fatalln(err)
	}
}
