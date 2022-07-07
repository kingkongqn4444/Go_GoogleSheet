FROM golang:alpine

COPY . .

CMD ["./be-wallet"]
