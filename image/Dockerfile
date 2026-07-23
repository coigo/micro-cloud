FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags '-s' -o /app/build .

# ----

FROM docker:dind AS runner

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/build /build

COPY ./start.sh /start.sh
RUN chmod +x /start.sh
RUN dos2unix /start.sh

EXPOSE 50051

ENTRYPOINT ["/start.sh"]