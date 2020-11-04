package autodns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
)

type ixFilter struct {
	Key      string `json:"key"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// GetZone ...
func (d *DNSProvider) GetZone(zone string) (*Zone, error) {
	dns, err := d.getVirtualNameserver(zone)
	if err != nil {
		log.Fatal(err)
	}

	req, err := d.makeRequest(http.MethodGet, path.Join("zone", zone, dns), nil)
	if err != nil {
		log.Fatal(err)
	}

	var resp *DataZoneResponse
	if err := d.sendRequest(req, &resp); err != nil {
		return nil, err
	}

	return resp.Data[0], nil
}

// UpdateZone ... Runs one update with adds/removes against the zone.
func (d *DNSProvider) UpdateZone(req *AutoDNSRequest, zone string) (*Zone, error) {
	zoneStream := &ZoneStream{
		Adds:    req.Adds,
		Removes: req.Removes,
	}

	return d.makeZoneUpdateRequest(zoneStream, zone)
}

func (d *DNSProvider) getVirtualNameserver(zone string) (string, error) {
	filter := ixFilter{
		Key:      "name",
		Operator: "EQUAL",
		Value:    zone,
	}

	// search first, to find nameserver
	search := make(map[string][]ixFilter)
	search["filters"] = append(search["filters"], filter)

	reqBody := &bytes.Buffer{}
	if err := json.NewEncoder(reqBody).Encode(search); err != nil {
		return "", err
	}

	req, err := d.makeRequest(http.MethodPost, path.Join("zone", "_search"), reqBody)
	if err != nil {
		log.Fatal(err)
	}

	var resp *DataZoneResponse
	if err := d.sendRequest(req, &resp); err != nil {
		return "", err
	}

	for _, data := range resp.Data {
		return data.VirtualNameServer, nil
	}

	return "", fmt.Errorf("Cannot find nameserver for '%s'", zone)
}
