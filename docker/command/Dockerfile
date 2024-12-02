FROM golang:1.23-alpine
RUN apk add build-base
WORKDIR /app
COPY . .
WORKDIR /app/cmd/commander
RUN go env -w CGO_ENABLED=1 && go build -o commander .
CMD [ "./commander" ]
