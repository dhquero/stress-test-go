FROM golang:1.24.2 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./cmd/stress -o stresstest

FROM scratch
WORKDIR /app
COPY --from=build /app/cmd/stress/stresstest .
ENTRYPOINT ["./stresstest"]
