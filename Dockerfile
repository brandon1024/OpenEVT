###
### Builder Stage
###
FROM docker.io/golang:1.27.4 AS builder

COPY . .
RUN CGO_ENABLED=0 go install -ldflags="-s -w" ./cmd/openevt

###
### Final Stage
###
FROM scratch

COPY --from=builder /go/bin/openevt /usr/bin/openevt

USER 65534:65534
CMD ["/usr/bin/openevt"]
