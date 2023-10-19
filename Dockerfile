FROM golang:1.20
WORKDIR /usr/src/app
COPY go.mod go.sum ./
COPY . ./

CMD ["go", "run", "src/main.go"]