FROM golang:1.11.5-alpine3.9 AS builder
COPY . /build/
WORKDIR /build
ENV GO111MODULE on
ENV CGO_ENABLED 0
RUN go build -mod=vendor -o ukfast

FROM alpine:3.9  
RUN apk --no-cache add ca-certificates
COPY --from=builder /build/ukfast /bin/ukfast
CMD ["/bin/sh"]