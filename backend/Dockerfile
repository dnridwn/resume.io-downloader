FROM golang:1.22.1-alpine

ENV APP_HOME "go/src/backend"
RUN mkdir -p "$APP_HOME"

WORKDIR "$APP_HOME"
COPY . .

RUN go get .
CMD ["go", "run", "."]
