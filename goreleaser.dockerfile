FROM alpine:3.19

RUN apk --no-cache add git git-lfs openssh-client

COPY ./git-semver /usr/local/bin/git-semver

ENTRYPOINT ["git", "semver"]
