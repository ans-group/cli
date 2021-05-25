FROM golang:1.16.4-alpine3.13 AS builder
RUN apk --no-cache add git
COPY . /build/
WORKDIR /build
ENV CGO_ENABLED 0
RUN go build -o ukfast -ldflags "-s -X 'main.VERSION=$(git describe --tags)' -X 'main.BUILDDATE=$(date +'%Y-%m-%dT%H:%M:%S')'"

FROM alpine:3.13
RUN apk --no-cache add ca-certificates bash bash-completion
COPY --from=builder /build/ukfast /bin/ukfast
RUN echo "source /usr/share/bash-completion/bash_completion" >> ~/.bashrc
RUN echo "source <(ukfast completion bash)" >> ~/.bashrc
ENTRYPOINT ["/bin/ukfast"]
