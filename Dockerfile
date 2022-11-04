FROM golang:1.19-alpine as build

WORKDIR /src

COPY . /src

RUN go test -v -vet=off ./...

RUN GOOS=linux GARCH=amd64 go build -o git-semver -ldflags="-s -w" cli/main.go

FROM alpine:3.16

RUN apk --no-cache add git openssh-client

COPY --from=build /src/git-semver /usr/local/bin

ENTRYPOINT ["git", "semver"]
