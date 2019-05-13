FROM golang:latest as builder
RUN mkdir /goapp
ADD . /goapp
ENV CONN="mongodb://localhost:27017"
WORKDIR /goapp

ARG service

RUN go build -o /bin/main ./$service


FROM debian:latest

COPY --from=builder /bin /app/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ 

EXPOSE 3000

CMD ["/app/main"]
