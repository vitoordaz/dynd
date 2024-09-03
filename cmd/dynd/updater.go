package main

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/vitoordaz/dynd/internal/dns"
	"github.com/vitoordaz/dynd/internal/myip"
)

const defaultTTL = 5 * time.Minute

func dnsUpdater(
	ctx context.Context,
	domain string,
	recordNames []string,
	dnsClient dns.Client,
	ipClient myip.Client,
	pollInterval time.Duration,
) error {
	logVerbose.Printf("starting updater for domain %s with poll interval %s\n", domain, pollInterval)

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		ip, err := ipClient.GetIPAddress(ctx)
		if err != nil {
			return err
		}
		logVerbose.Printf("current IP address is %s\n", ip)

		records, err := dnsClient.GetRecords(ctx, domain)
		if err != nil {
			return err
		}

		for _, recordName := range recordNames {
			aRecord := getRecord(records, "A", recordName)
			if aRecord == nil {
				logVerbose.Printf("there are no A record with name %s for domain %s\n", domain, recordName)
				logVerbose.Printf("creating A record %s [%s] for domain %s\n", recordName, ip, domain)
				if err := dnsClient.CreateRecord(ctx, domain, recordName, "A", []string{ip}, defaultTTL); err != nil {
					return err
				}
				logVerbose.Printf("created A record %s [%s] for domain %s\n", recordName, ip, domain)
			} else {
				logVerbose.Printf("domain %s has A record %s %s\n", domain, aRecord.Name, strings.Join(aRecord.Values, ", "))
				if !reflect.DeepEqual(aRecord.Values, []string{ip}) {
					logVerbose.Printf("updating domain %s A record %s value to [%s]\n", domain, aRecord.Name, ip)
					if err := dnsClient.ReplaceRecord(ctx, domain, recordName, "A", []string{ip}, defaultTTL); err != nil {
						return err
					}
					logVerbose.Printf("updated domain %s A record %s value to [%s]\n", domain, aRecord.Name, ip)
				}
			}
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return nil
		}
	}
}

func getRecord(records []*dns.Record, recordType, recordName string) *dns.Record {
	for _, record := range records {
		if record.Type == recordType && record.Name == recordName {
			return record
		}
	}
	return nil
}
