FROM golang:1.23.1 as builder
WORKDIR /app
COPY . .
RUN go mod download
ENV GOOS=linux
ENV CGO_ENABLED=0
RUN go build -o main ./cmd/auth-service

FROM gcr.io/distroless/base
USER 1000:1000
COPY --from=builder /app/main /app/main
EXPOSE 8081
WORKDIR /app
CMD [ "./main" ]