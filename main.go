package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/hostwithquantum/record-updater/autodns"
	"github.com/urfave/cli/v2"
)

var version string

func main() {
	app := &cli.App{
		Name:    "record-updater",
		Usage:   "A cli tool to (bulk) update DNS records on InternetX' AutoDNS",
		Version: version,
		Authors: []*cli.Author{
			{
				Name:  "Till Klampaeckel",
				Email: "till@planetary-quantum.com",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "customer",
				Usage:   "The customer to pull nodes for",
				EnvVars: []string{"QUANTUM_CUSTOMER"},
			},
			&cli.StringFlag{
				Name:    "config",
				Usage:   "Settings for groups, etc. for --list",
				Value:   "config.ini",
				EnvVars: []string{"QUANTUM_DNS_CONFIG"},
			},
			&cli.StringFlag{
				Name:    "target",
				Usage:   "Set the target via command-line (overrules .ini)",
				Value:   "",
				EnvVars: []string{"QUANTUM_DNS_TARGET"},
			},
		},
		Action: func(c *cli.Context) error {
			provider, err := autodns.NewDNSProvider()
			if err != nil {
				return err
			}

			customer := c.String("customer")
			if customer == "" {
				return errors.New("Please add a customer")
			}
			log.Printf("Customer: %s", customer)

			ru := NewRecordUpdater(c.String("config"))
			dnsRecords, err := ru.GetStrings("records", ",", "")
			if err != nil {
				return err
			}

			zoneName, err := ru.GetString("zone", "")
			if err != nil {
				return err
			}

			log.Printf("Zone: %s\n", zoneName)

			zone, err := provider.GetZone(zoneName)
			if err != nil {
				return err
			}

			existingRecords := make(map[string]*autodns.ResourceRecord)
			for _, zone := range zone.ResourceRecords {
				existingRecords[zone.Name] = zone
			}

			var recordValue string
			if c.String("target") != "" {
				recordValue = c.String("target")
			} else {
				recordValue, err = ru.GetString("target", "")
				if err != nil {
					return err
				}
			}

			if recordValue == "" {
				return fmt.Errorf("Missing 'target' via --target or '%s'", c.String("config"))
			}

			isCNAME := true
			if net.ParseIP(recordValue) != nil {
				isCNAME = false
				log.Println("It's an IP!")
			} else {
				recordValue = fmt.Sprintf(recordValue, customer)
				log.Println("It's a CNAME!")
			}

			log.Printf("Record Value: %s\n", recordValue)

			request := autodns.NewAutoDNSRequest(recordValue)

			for _, record := range dnsRecords {
				finalRecord := fmt.Sprintf(record, customer)
				existingRecord, ok := existingRecords[finalRecord]
				if !ok {
					if finalRecord == customer {
						request.AddA(finalRecord)
						continue
					}
					if isCNAME {
						request.AddCNAME(finalRecord)
					} else {
						request.AddA(finalRecord)
					}
					continue
				}

				if finalRecord == customer {
					request.AddA(finalRecord)
					continue
				}

				if isCNAME {
					if existingRecord.Type != "CNAME" {
						request.AddCNAME(finalRecord)
						request.RemoveRecord(existingRecord)
						continue
					}
				} else {
					if existingRecord.Type != "A" {
						request.AddA(finalRecord)
						request.RemoveRecord(existingRecord)
						continue
					}
				}

				if existingRecord.Value == recordValue {
					continue
				}

				if isCNAME {
					request.AddCNAME(finalRecord)
				} else {
					request.AddA(finalRecord)
				}
				request.RemoveRecord(existingRecord)
			}

			log.Println("We are commiting the following:")
			log.Printf("   Adds: %d\n", len(request.Adds))
			log.Printf("Removes: %d\n", len(request.Removes))

			for _, r := range append(request.Adds, request.Removes...) {
				fmt.Printf(" - %s IN %s %s\n", r.Name, r.Type, r.Value)
			}

			provider.UpdateZone(request, zoneName)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
