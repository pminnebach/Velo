### About velo.go

This app downloads the bike availability statistics from the velo availability_map and pushes it to an InfluxDB server.

### Docker

#### Build docker image
```
docker build -t velo .
```

#### Run docker image
Substitute the environment variables with your own
```
 docker run -d \
    --name=velo
    -e "INFLUXDB_HOST=http://localhost:8086" \
    -e "INFLUXDB_DATABASE=test" \
    -e "INFLUXDB_USERNAME=root" \
    -e "INFLUXDB_PASSWORD=root" \
    -e "VELO_URL=https://..." \
    velo
```
