FROM golang:1.15.2-buster as builder
WORKDIR portgw

COPY . ./
ARG version=dev
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-X main.version=$version" -o portgw ./cmd/portgw/main.go

FROM scratch
COPY --from=builder /go/portgw/portgw .
CMD ["./portgw"]