FROM golang:1.23-alpine
RUN apk add build-base
WORKDIR /app
COPY . .
WORKDIR /app/cmd/queryer
RUN go env -w CGO_ENABLED=1 && go build -o queryer .
CMD [ "./queryer" ]
