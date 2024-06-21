#base go image

FROM golang:1.22-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest && \
    ln -s /go/bin/linux_amd64/migrate /usr/local/bin/migrate

RUN CGO_ENABLED=0 go build -o prmv ./cmd/web

RUN chmod +x ./prmv

#build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app /app
#add go tools binares
COPY --from=builder /go/bin/migrate /app/migrate

RUN chmod +x /app/entrypoint.sh

WORKDIR /app

CMD ["/app/prmv"]