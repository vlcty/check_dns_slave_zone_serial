# check_dns_slave_zone_serial

Icinga plugin to check your DNS slaves serial numbers against your master.

## Usage

The plugin takes three arguments:

- zone: The zone/domain to query. For example `vlcty.de`
- master: IP address of your master server
- slaves: IP addresses of your slave servers (comma seperated)

Example call:

> :-$ go run check_dns_slave_zone_serial.go --zone vlcty.de --master "[2001:db9:150:b00b::abcd]" --slaves "[2a01:4f8:1c1c:4213::1af],[2a0c:25
00:571:b00:14c9:1f5b:e73a:c45b]"   
> OK - Zone: vlcty.de, Master serial: 2021030201, Slaves: [2a01:4f8:1c1c:4213::1af] is in sync, [2a0c:2500:571:b00:14c9:1f5b:e73a:c45b] is in sync

## Installation

Compile the go source file for your arch. For linux on amd64: `GOOS=linux GOARCH=amd64 go build -o check_dns_slave_zone_serial check_dns_slave_zone_serial.go`

## IP formats

IP addresses must be in brackets. Legacy IP address don't need to be.
