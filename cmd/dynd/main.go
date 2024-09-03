package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/vitoordaz/dynd/internal/dns"
	"github.com/vitoordaz/dynd/internal/helpers"
	"github.com/vitoordaz/dynd/internal/myip"
)

var (
	logVerbose = log.New(os.Stdout, "D: ", 0)
	logError   = log.New(os.Stderr, "ERROR: ", 0)

	domain           = flag.String("domain", "", "domain to update")
	recordNames      = flag.String("record-names", "*", "a comma separated list of record names that will be updated")
	gandiAccessToken = flag.String("gandi-access-token", "", "gandi access token")
	pollInterval     = flag.Int("poll-interval", 60, "IP address polling interval in seconds")
)

const (
	exitCodeError = 2
	exitCodeOk    = 0
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	flag.Usage = printUsage
	flag.Parse()
	if *domain == "" {
		logError.Println("domain is required")
		os.Exit(exitCodeError)
	}
	if *gandiAccessToken == "" {
		logError.Println("gandi access token is required")
		os.Exit(exitCodeError)
	}

	dnsClient, err := dns.NewGandiClient(ctx, *gandiAccessToken)
	if err != nil {
		logError.Println(err)
		os.Exit(exitCodeError)
	}

	myIPClient := myip.NewIPIFYClient()

	if err := dnsUpdater(
		ctx,
		*domain,
		helpers.TrimStringSpaces(strings.Split(*recordNames, ",")),
		dnsClient,
		myIPClient,
		time.Duration(*pollInterval)*time.Second,
	); err != nil {
		logError.Println(err)
		os.Exit(exitCodeError)
	}
	os.Exit(exitCodeOk)
}
