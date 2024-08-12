FROM golang:1.22-alpine as builder
RUN apk add --no-cache git
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM scratch
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/dataBase/ipToCountryDB.csv /dataBase/ipToCountryDB.csv

# Copy the .env file
COPY .env .
EXPOSE 8080
CMD ["./main"]