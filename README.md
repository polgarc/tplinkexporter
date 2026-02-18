# TPLink EasySmart Switch Exporter

Exports port stats from TPLink's EasySmart Switches. Very alpha, basic functionality so far:
- Port Stats

Tested on:

TPLink EasySmart Gigabit 8 Port Switch (TL-SG108E):
- v2 Firmware: 1.0.1 Build 20211021 Rel.56865
- v3 Firmware: 1.0.0 Build 20171214 Rel.70905 (based on web GUI)
- v4 Firmware: 1.0.0 Build 20181120 Rel.40749
- v5 Firmware: 1.0.0 Build 20230218 Rel.51358

Should work on the switches of the same family, but untested personally:
- TPLink EasySmart Gigabit 16 Port Switch TL-SG116E

## Grafana Dashboard

Basic Grafana dashboard using the exporter - https://grafana.com/grafana/dashboards/12517

## Usage

Command line:

`go run main.go --host 10.0.0.3 --username admin --password admin`

Docker:

`docker run -p 9717:9717 polgarc/tplinkexporter --host 10.0.0.3 --username admin --password admin`

Docker compose:

```yaml
services:
  tplink-exporter:
    image: polgarc/tplinkexporter:latest
    environment:
      HOST: tplink-1.lan
      USERNAME: admin
      PASSWORD: admin
    ports:
      - 9717:9717
```

## Metrics

Exported on `<host>:9717/metrics`:

```
tplinkexporter_portstats_state{portnum="1"-"8",host="host"}
tplinkexporter_portstats_linkstatus{portnum="1"-"8",host="host"}
tplinkexporter_portstats_rxgoodpkt{portnum="1"-"8",host="host"}
tplinkexporter_portstats_rxbadpkt{portnum="1"-"8",host="host"}
tplinkexporter_portstats_txgoodpkt{portnum="1"-"8",host="host"}
tplinkexporter_portstats_txbadpkt{portnum="1"-"8",host="host"}
```
