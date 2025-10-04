FROM golang:1.24-trixie AS builder

# RUN apk add --update gcc musl-dev
RUN mkdir -p /myapp
ADD . /myapp
WORKDIR /myapp

RUN  go install github.com/swaggo/swag/cmd/swag@latest  &&  go generate . && GOOS=linux CGO_ENABLED=1  go build -ldflags='-extldflags=-static'  -o myapp


FROM containers.intersystems.com/intersystems/iris-community:latest-cd

COPY --from=builder /myapp/myapp /myapp
COPY health.sh /health.sh
COPY irisinit.sh /irisinit.sh

HEALTHCHECK --interval=30s --timeout=30s --start-period=20s --retries=3 CMD [ "/health.sh" ]

CMD [ "--after", "/irisinit.sh" ]