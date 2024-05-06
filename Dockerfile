FROM golang:1.22 as builder

RUN apt-get update -y
RUN apt-get install -y golang libc6
# /usr/sbin/update-ca-certificates

WORKDIR /workspace

COPY / /workspace/
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build .

ARG REGION
ARG SECRETACCESSKEY
ARG ACCESSKEYID
ARG URL
ARG LOGILICA_TOKEN

ENV REGION $REGION
ENV SECRETACCESSKEY $SECRETACCESSKEY
ENV ACCESSKEYID $ACCESSKEYID
ENV URL $URL
ENV LOGILICA_TOKEN $LOGILICA_TOKEN


FROM debian:11
RUN apt-get update -y && apt-get install -y ca-certificates
RUN /usr/sbin/update-ca-certificates
WORKDIR /
COPY --from=builder /workspace/event-listener .
EXPOSE 8080
ENTRYPOINT ["/event-listener"]