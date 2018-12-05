# Smsc balance exporter

The smsc balance exporter for [prometheus](https://prometheus.io) allows exporting balance for [smsc gateway](https://smsc.ru)

## How it works
Exporter querying balance every hour (by default) and store it value in memory.
When prometheus make request, exporter retrieve balance value from memory for make response.

## Configuration
You must set environment variables:

* `SMSC_LOGIN` - your login
* `SMSC_PASSWORD` - your password

## Command-line flags

* `listen-address` - The address to listen on for HTTP requests. (Default: `0.0.0.0:9601`)
* `interval` - Interval (in seconds) for querying balance. (Default: `0.0.0.0:9601`)

## Running with systemctl

```sh
# Set variables you need
SMSC_EXPORTER_VERSION=v0.1.0-alpha
SMSC_EXPORTER_PLATFORM=linux
SMSC_EXPORTER_ARCH=amd64

# Download release
wget https://github.com/xxxcoltxxx/smsc-balance-exporter/releases/download/${SMSC_EXPORTER_VERSION}/smsc_balance_exporter_${SMSC_EXPORTER_VERSION}_${SMSC_EXPORTER_PLATFORM}_${SMSC_EXPORTER_ARCH}.tar.gz
tar xvzf smsc_balance_exporter_${SMSC_EXPORTER_VERSION}_${SMSC_EXPORTER_PLATFORM}_${SMSC_EXPORTER_ARCH}.tar.gz
mv ./smsc_balance_exporter_${SMSC_EXPORTER_VERSION}_${SMSC_EXPORTER_PLATFORM}_${SMSC_EXPORTER_ARCH} /usr/local/bin/smsc_balance_exporter

# Add service to systemctl
cat <<EOF >> /etc/systemd/system/smsc_balance_exporter.service
[Unit]
Description=Smsc Balance Exporter
Wants=network-online.target
After=network-online.target

[Service]
Environment="SMSC_LOGIN=<your_login> SMSC_PASSWORD=<your_password>"
User=balance_exporter
Group=balance_exporter
Type=simple
ExecStart=/usr/local/bin/smsc_balance_exporter

[Install]
WantedBy=multi-user.target
EOF

# Reload systemctl configuration and run service
systemctl daemon-reload
systemctl start smsc_balance_exporter
systemctl status smsc_balance_exporter
```
