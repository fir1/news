FROM golang:1.20.4-alpine as build-backend

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy && \
    go mod download -x

COPY . .

# Run tests
RUN go test ./...

RUN CGO_ENABLED=0 go build -o /bin/backend cmd/*.go

FROM scratch
COPY --from=build-backend /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-backend /bin/backend /bin/backend
COPY --from=build-backend /app/docs/ /docs/
ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /zoneinfo.zip
ENV ZONEINFO /zoneinfo.zip
EXPOSE 8080

ENTRYPOINT ["/bin/backend"]