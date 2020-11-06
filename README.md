# record-updater

A cli tool to create records on AutoDNS (by InternetX GmbH).

For usage: `./record-updater help`

Basic format/configuration is explained in: [config.ini-dist](./config.ini-dist).

We use code from go-acme to talk to the AutoDNS service:
https://go-acme.github.io/lego/dns/autodns/ 

## Example:

_For records, `%s` replacement happens with whatever `QUANTUM_CUSTOMER`.

Assume the following:

```
$ export QUANTUM_CUSTOMER=acme
$ cat config.ini
target=127.0.0.1
zone=example.org
records=%s.customer,*.%s.customer
```

Running the tool creates the following records:

```
; zone: .your-zone.example.org
acme.customer IN A 127.0.0.1
*acme.customer IN A 127.0.0.1
```

Or the following hosts:

- `acme.customer.example.org`
- `*.acme.customer.example.org`
