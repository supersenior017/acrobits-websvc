package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"bitbucket.org/x31a/acrobits-websvc/src/websvc"
)

const (
	FlagConfig  = "c"
	ExitSuccess = 0
)

type Opts struct {
	config websvc.Config
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
		fmt.Println(websvc.Version)
		os.Exit(ExitSuccess)
	}
	if !opts.config.IsSet() {
		opts.config.SetDefaults()
	}
	opts.config.Balance.Func = getBalance
	opts.config.Rate.Func = getRate
	return opts
}

func getBalance(
	ctx context.Context,
	account websvc.Account,
) (websvc.Balance, error) {
	return websvc.Balance{}, fmt.Errorf("NotImplemented")
}

func getRate(
	ctx context.Context,
	rate websvc.RateParams,
) (websvc.Rate, error) {
	return websvc.Rate{}, fmt.Errorf("NotImplemented")
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
	log.Printf("Starting on: %q\n", opts.config.Addr)
	if err := websvc.ListenAndServe(
		ctx,
		opts.config,
	); err != nil && err != http.ErrServerClosed {
		log.Fatalln(err)
	}
}
