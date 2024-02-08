FROM golang:1.21 as builder

WORKDIR /go/src/server

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GGOOS=linux

RUN go env -w GO111MODULE="on"

COPY ../.. .
RUN go build -tags=nomsgpack -a -installsuffix nocgo -o /app cmd/main.go

FROM debian:buster-slim

RUN apt update && apt install -y curl

COPY --from=builder /app /opt/server/
COPY settings.yml /opt/server/

ENTRYPOINT ["/opt/server/app"]
