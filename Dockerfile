FROM golang:1.20-alpine

WORKDIR /trailbrake_api

COPY . .

RUN go build -o trailbrake_api

EXPOSE 8080

CMD [ "./trailbrake_api" ]