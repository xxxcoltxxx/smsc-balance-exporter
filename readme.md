# Balance exporter for https://smsc.ru service

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
* `interval` - Interval (in seconds) for querying balance. (Default: `3600`)

## Running with docker

```sh
docker run \
    -e SMSC_LOGIN=<your-login> \
    -e SMSC_PASSWORD=<your-password> \
    -p 9601:9601 \
    --restart=unless-stopped \
    --name smsc-balance-exporter \
    -d \
    xxxcoltxxx/smsc-balance-exporter
```

## Running with docker-compose

Create configuration file. For example, file named `docker-compose.yaml`:

```yaml
version: "3"

services:
  smsc-balance-exporter:
    image: xxxcoltxxx/smsc-balance-exporter
    restart: unless-stopped
    environment:
      SMSC_LOGIN: <your-login>
      SMSC_PASSWORD: <your-password>
    ports:
      - 9601:9601
```

Run exporter:
```sh
docker-compose up -d
```

Show service logs:
```sh
docker-compose logs -f smsc-balance-exporter
```

## Running with systemctl

Set variables you need:
```sh
SMSC_EXPORTER_VERSION=v0.1.3-beta.2
SMSC_EXPORTER_PLATFORM=linux
SMSC_EXPORTER_ARCH=amd64
SMSC_LOGIN=<your_login>
SMSC_PASSWORD=<your_password>
```

Download release:
```sh
wget https://github.com/xxxcoltxxx/smsc-balance-exporter/releases/download/${SMSC_EXPORTER_VERSION}/smsc_balance_exporter_${SMSC_EXPORTER_VERSION}_${SMSC_EXPORTER_PLATFORM}_${SMSC_EXPORTER_ARCH}.tar.gz
tar xvzf smsc_balance_exporter_${SMSC_EXPORTER_VERSION}_${SMSC_EXPORTER_PLATFORM}_${SMSC_EXPORTER_ARCH}.tar.gz
mv ./smsc_balance_exporter_${SMSC_EXPORTER_VERSION}_${SMSC_EXPORTER_PLATFORM}_${SMSC_EXPORTER_ARCH} /usr/local/bin/smsc_balance_exporter
```

Add service to systemctl. For example, file named `/etc/systemd/system/smsc_balance_exporter.service`:
```sh
[Unit]
Description=Smsc Balance Exporter
Wants=network-online.target
After=network-online.target

[Service]
Environment="SMSC_LOGIN=${SMSC_LOGIN}"
Environment="SMSC_PASSWORD=${SMSC_PASSWORD}"
Type=simple
ExecStart=/usr/local/bin/smsc_balance_exporter

[Install]
WantedBy=multi-user.target
```

Reload systemctl configuration and run service
```sh
systemctl daemon-reload
systemctl restart smsc_balance_exporter
```

Check service status:
```sh
systemctl status smsc_balance_exporter
```

Check service logs:
```sh
journalctl -fu smsc_balance_exporter
```
