FROM golang:1.12
RUN mkdir -p /go/src/app/
WORKDIR /go/src/app/
COPY src/ .

RUN go install

EXPOSE 8080

ENTRYPOINT ["/go/bin/app"]
