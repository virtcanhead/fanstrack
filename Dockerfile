FROM alpine:3.9

RUN apk add --no-cache ca-certificates

ADD fanstrack /

CMD ["/fanstrack"]
