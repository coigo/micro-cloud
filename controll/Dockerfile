FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags '-s' -o /app/build .

# ----

FROM alpine:3.19.9 AS runner

COPY --from=builder /app/build /build

EXPOSE 50051
ENTRYPOINT ["/build"]