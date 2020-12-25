ARG PLATFORM=linux/amd64
FROM --platform=${PLATFORM} alpine:3.12
COPY artifacts/build/release/linux/amd64/traefik-acme /
ENTRYPOINT [ "/traefik-acme" ]
