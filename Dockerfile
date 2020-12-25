FROM scratch
COPY artifacts/build/release/linux/amd64/traefik-acme /traefik-acme
ENTRYPOINT [ "/traefik-acme" ]
