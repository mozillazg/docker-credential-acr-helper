FROM golang:1.19-buster as BUILD
WORKDIR /go/src/github.com/mozillazg/docker-credential-acr-helper
ADD . .

WORKDIR /go/src/github.com/mozillazg/docker-credential-acr-helper/examples/go-containerregistry-auth
ENV GOPROXY https://goproxy.cn
RUN CGO_ENABLED=0 go build -o get-acr-digest main.go && cp get-acr-digest /tmp/get-acr-digest

FROM alpine:3.15
WORKDIR /app/
USER 65534

COPY --from=BUILD /tmp/get-acr-digest ./
CMD ["/app/get-acr-digest"]
