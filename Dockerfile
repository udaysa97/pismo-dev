FROM golang:1.21-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -tags dynamic  -o ./service cmd/main.go

# Non-Root user
RUN addgroup -g 7777 prismo &&\
    adduser -D -h /home/prismo prismo -s /bin/false -u 7777 -G prismo

USER prismo
RUN mkdir /home/prismo/app
WORKDIR /home/prismo/app

COPY --from=build /app/service ./service

ENV REDIS_POOL_SIZE=32
ENV FLAG_CUSTOM_CIRCUIT_BREAKER_DQL=true
ENV DD_TRACE_PROPAGATION_STYLE_EXTRACT=Datadog
ENV DD_TRACE_PROPAGATION_STYLE_INJECT=Datadog
ENV DD_TRACE_PROPAGATION_STYLE=Datadog

ENTRYPOINT [ "/usr/bin/dumb-init", "--" ]
CMD [ "/home/prismo/app/service" ]