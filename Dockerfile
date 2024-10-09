FROM golang:alpine as builder
ENV CGO_ENABLED 0
ENV GOOS linux
WORKDIR /build
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /build/app .

FROM alpine
WORKDIR /app
COPY --from=builder /build/app /app/app
CMD ["./app"]

