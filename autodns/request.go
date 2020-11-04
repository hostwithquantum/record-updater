package autodns

// AutoDNSRequest ... A simple wrapper that builds properties for ZoneStream.
type AutoDNSRequest struct {
	target  string
	Adds    []*ResourceRecord
	Removes []*ResourceRecord
}

// NewAutoDNSRequest ... factory/CTOR
func NewAutoDNSRequest(target string) *AutoDNSRequest {
	return &AutoDNSRequest{
		target: target,
	}
}

// AddA ...
func (r *AutoDNSRequest) AddA(name string) {
	r.add(&ResourceRecord{
		Name:  name,
		Type:  "A",
		TTL:   9600,
		Value: r.target,
	})
}

// AddCNAME ...
func (r *AutoDNSRequest) AddCNAME(name string) {
	r.add(&ResourceRecord{
		Name:  name,
		Type:  "CNAME",
		TTL:   9600,
		Value: r.target,
	})
}

// RemoveRecord ...
func (r *AutoDNSRequest) RemoveRecord(record *ResourceRecord) {
	r.Removes = append(r.Removes, record)
}

func (r *AutoDNSRequest) add(record *ResourceRecord) {
	r.Adds = append(r.Adds, record)
}
