# Builder runs natively on the build host and cross-compiles for the target platform
FROM --platform=$BUILDPLATFORM quay.io/projectquay/golang:1.25 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -v -o /out/gbot \
    -ldflags "-X github.com/khadzhamov/mkbot/cmd.AppVersion=${VERSION}"

# CA certificates for HTTPS (Telegram API); path differs between distros
RUN cp -L /etc/ssl/certs/ca-certificates.crt /out/ca-certificates.crt 2>/dev/null \
    || cp -L /etc/pki/tls/certs/ca-bundle.crt /out/ca-certificates.crt

FROM scratch
WORKDIR /
COPY --from=builder /out/kbot /kbot
COPY --from=builder /out/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/kbot"]