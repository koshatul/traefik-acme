# traefik-acme

Extract certificates from acme.json created by traefik.

## Why

Using traefik to do the work of certbot is good, but sometimes you have other services that need access to the certificate, this CLI tool extracts them out so you can use them outside of traefik.

## Usage

```text
Usage:
  traefik-acme <domain> [flags]
  traefik-acme [command]

Available Commands:
  help        Help about any command
  version     Print the version

Flags:
  -a, --acme string   Location of acme.json file (default "/etc/traefik/acme.json")
  -c, --cert string   Location to write out certificate (default "cert.pem")
  -d, --debug         Debug output
      --exit-code     Exit with exit-code 99 if files updated
      --force         Force writing to file even if not updated
  -h, --help          help for traefik-acme
  -k, --key string    Location to write out key file (default "key.pem")
      --version       version for traefik-acme

Use "traefik-acme [command] --help" for more information about a command.
```

Running from command line.

```shell
traefik-acme -a /config/acme.json -c /etc/service/cert.pem -k /etc/service/key.pem servicename.domain.com
```

If you want to use it in a script (for cron)

```shell
traefik-acme --exit-code -a /config/acme.json -c /etc/service/cert.pem -k /etc/service/key.pem servicename.domain.com
if [ $? == 99 ]; then
    systemctl reload service
fi
```

### Docker


## Development

```shell
make test
```

```shell
ginkgo ./src/...
```
