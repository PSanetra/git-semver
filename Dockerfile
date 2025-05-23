FROM golang:1.24-alpine as build

WORKDIR /src

COPY . /src

RUN go test -v -vet=off ./...

RUN GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -o git-semver -ldflags="-s -w" cli/main.go

FROM alpine:3.21

RUN apk --no-cache add git git-lfs openssh-client

COPY --from=build /src/git-semver /usr/local/bin

ENTRYPOINT ["git", "semver"]
