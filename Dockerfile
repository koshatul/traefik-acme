FROM golang:1.14-alpine as builder

RUN apk add --no-cache make curl bash git zip unzip util-linux gcc

WORKDIR /go/src/app
COPY . .

RUN CGO_ENABLED=0 make release
# RUN go get -d -v ./...
# RUN go install -v ./...

FROM scratch

COPY --from=builder /go/src/app/artifacts/build/release/linux/amd64/traefik-acme /traefik-acme

ENTRYPOINT [ "/traefik-acme" ]
