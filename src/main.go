package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	acrobitsbalance "./lib"
)

const (
	FlagConfig = "c"

	ExOk     = 0
	ExArgErr = 2
)

type Opts struct {
	config acrobitsbalance.Config
}

func parseArgs() *Opts {
	opts := &Opts{}
	isHelp := flag.Bool("h", false, "Print help and exit")
	isVersion := flag.Bool("V", false, "Print version and exit")
	flag.Var(&opts.config, FlagConfig, "Path to configuration file")
	flag.Parse()
	if *isHelp {
		flag.Usage()
		os.Exit(ExOk)
	}
	if *isVersion {
		fmt.Println(acrobitsbalance.Version)
		os.Exit(ExOk)
	}
	if opts.config.FilePath() == "" {
		fmt.Fprintf(os.Stderr, "-%s required\n", FlagConfig)
		os.Exit(ExArgErr)
	}
	return opts
}

func getBalance(username, password string) (float64, error) {
	return 0, nil
}

func main() {
	opts := parseArgs()
	opts.config.Func = getBalance
	log.Println("Starting on:", opts.config.Addr)
	log.Fatalln(acrobitsbalance.ListenAndServe(opts.config))
}
