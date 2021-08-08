# traefik-acme

[![ci](https://github.com/koshatul/traefik-acme/actions/workflows/ci.yml/badge.svg)](https://github.com/koshatul/traefik-acme/actions/workflows/ci.yml)

Extract certificates from acme.json created by traefik.

## Why

Using traefik to do the work of certbot is good, but sometimes you have other services that need access to the certificate, this CLI tool extracts them out so you can use them outside of traefik.

## Usage

```text
Usage:
  traefik-acme <domain> [flags]

Flags:
  -a, --acme string                   Location of acme.json file (default "/etc/traefik/acme.json")
  -c, --cert string                   Location to write out certificate (default "cert.pem")
  -r, --certificate-resolver string   Certificate Resovler name from traefik config (default "acme")
  -d, --debug                         Debug output
      --exit-code                     Exit with exit-code 99 if files updated
      --force                         Force writing to file even if not updated
  -h, --help                          help for traefik-acme
  -k, --key string                    Location to write out key file (default "key.pem")
  -v, --version                       version for traefik-acme
```

Running from command line, `myresolver` is the name of the certificate resolver defined in your traefik configuration.  Example from the [documentation](https://doc.traefik.io/traefik/https/acme/) `[certificatesResolvers.myresolver.acme]` would mean using `myresolver` as the parameter.

```shell
traefik-acme -r myresolver -a /config/acme.json -c /etc/service/cert.pem -k /etc/service/key.pem servicename.domain.com
```

If you want to use it in a script (for cron)

```shell
traefik-acme --exit-code -r mycertresolver -a /config/acme.json -c /etc/service/cert.pem -k /etc/service/key.pem servicename.domain.com
if [ $? == 99 ]; then
    systemctl reload service
fi
```

### Traefik v1

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

```shell
docker run --rm \
 -v "/docker/traefik/config/:/input" \
 -v "/docker/myservice/certs:/output" \
 --workdir /output \
 koshatul/traefik-acme:latest --acme "/input/acme.json" domain.example.com
```

The example expects the `acme.json` to be in `/docker/traefik/config` and to write the `cert.pem` and `key.pem` to `/docker/myservice/certs`.

## Development

```shell
make test
```

```shell
ginkgo ./src/...
```
