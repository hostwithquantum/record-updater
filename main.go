package main

import (
	"fmt"
	"log"
	"net"

	"github.com/hostwithquantum/record-updater/autodns"
	"gopkg.in/ini.v1"
)

func main() {
	provider, err := autodns.NewDNSProvider()
	if err != nil {
		log.Fatal(err)
	}

	customer := "luzilla"

	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatal(err)
	}

	// .customer.planetary-quantum.com > $customer.customer.planetary-quantum.net
	dnsRecords := cfg.Section("").Key("records").Strings(",")

	zoneName := cfg.Section("").Key("zone").String()
	log.Println(zoneName)

	zone, err := provider.GetZone(zoneName)
	if err != nil {
		log.Fatal(err)
	}

	existingRecords := make(map[string]*autodns.ResourceRecord)
	for _, zone := range zone.ResourceRecords {
		existingRecords[zone.Name] = zone
	}

	recordValue := cfg.Section("").Key("target").String()
	isCNAME := true
	if net.ParseIP(recordValue) != nil {
		isCNAME = false
		log.Println("It's an IP!")
	} else {
		recordValue = fmt.Sprintf(recordValue, customer)
		log.Println("It's a CNAME!")
	}

	log.Println(recordValue)

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
}
