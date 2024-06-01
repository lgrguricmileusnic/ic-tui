FROM golang:1.22-alpine as builder
COPY . /ic-tui
WORKDIR /ic-tui
RUN cd ./cmd/ic-tui && go install
CMD ["ic-tui"]

FROM ubuntu:latest
# RUN apk update && apk upgrade
COPY --from=builder /go/bin/ic-tui .
EXPOSE 8080/tcp
EXPOSE 22/tcp
CMD ["/ic-tui", "-ssh"]