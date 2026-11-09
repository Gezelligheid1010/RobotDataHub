FROM golang:1.23-bookworm AS build
WORKDIR /src
COPY . .
RUN go mod tidy && CGO_ENABLED=0 go build -o /app ./cmd/server
FROM gcr.io/distroless/static-debian12
COPY --from=build /app /app
USER 65532
EXPOSE 8080
ENTRYPOINT ["/app"]
