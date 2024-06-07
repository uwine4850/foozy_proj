FROM golang:1.22.3
WORKDIR /usr/src/app
COPY go.mod go.sum ./
COPY . ./

CMD ["go", "run", "src/main.go"]