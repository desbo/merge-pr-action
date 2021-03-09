FROM golang:1.16-alpine AS builder

WORKDIR /go/src/github.com/permutive/github-actions/merge-pr
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/action

FROM scratch

COPY --from=builder /bin/action /bin/action
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT [ "/bin/action" ]