FROM golang:1.8

ENV INFLUXDB_HOST="http://localhost:8086" INFLUXDB_DATABASE="-" INFLUXDB_USERNAME="root" INFLUXDB_PASSWORD="root" VELO_URL="-"

WORKDIR /go/src/app
COPY . .

RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run"]
