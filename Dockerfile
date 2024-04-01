FROM golang:1.22 as builder

RUN apt-get update -y
RUN apt-get install -y golang libc6

WORKDIR /workspace

COPY / /workspace/
RUN go build .

FROM debian:11
WORKDIR /
COPY --from=builder /workspace/event-listener .
EXPOSE 8080
ENTRYPOINT ["/event-listener"]