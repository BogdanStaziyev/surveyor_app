FROM golang:1.18-alpine

WORKDIR /app/src

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

COPY internal/infra/database/migrations ../migrations

RUN go build -o ../surveyor_app cmd/server/main.go

WORKDIR /app

RUN rm -rf src

EXPOSE 8080

CMD [ "/app/surveyor_app" ]