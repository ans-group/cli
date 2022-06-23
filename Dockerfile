FROM golang:1.18-alpine3.16 AS builder
RUN apk --no-cache add git
COPY . /build/
WORKDIR /build
ENV CGO_ENABLED 0
RUN go build -o ans -ldflags "-s -X 'main.VERSION=$(git describe --tags)' -X 'main.BUILDDATE=$(date +'%Y-%m-%dT%H:%M:%S')'"

FROM alpine:3.16
RUN apk --no-cache add ca-certificates bash bash-completion
COPY --from=builder /build/ans /bin/ans
RUN echo "source /usr/share/bash-completion/bash_completion" >> ~/.bashrc
RUN echo "source <(ans completion bash)" >> ~/.bashrc
ENTRYPOINT ["/bin/ans"]
