FROM golang:1.24-alpine AS build

RUN adduser --uid 10000 --disabled-password antivax

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build

FROM scratch
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/antivax /antivax
USER antivax
CMD ["/antivax"]
