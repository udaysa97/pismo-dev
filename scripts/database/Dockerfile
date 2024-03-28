FROM golang:1.19-alpine

RUN apk add --update --no-cache dumb-init git

WORKDIR /app
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Define env vars
ARG ENV
ENV ENV=${ENV:-staging}

COPY . /app/scripts/database

ENTRYPOINT [ "/usr/bin/dumb-init", "--", "/app/scripts/database/entrypoint.sh" ]