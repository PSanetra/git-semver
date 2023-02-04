FROM alpine:3.17

RUN apk --no-cache add git git-lfs openssh-client

COPY ./git-semver /usr/local/bin/git-semver

ENTRYPOINT ["git", "semver"]
