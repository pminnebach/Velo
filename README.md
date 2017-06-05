### About main.go

The sensitive configuration data is stored in `config.json`.

#### config.json sample

```json
{
  "Host"      : "",
  "Database"      : "",
  "Username"      : "",
  "Password"      : "",
  "Url" 	  : ""
}
```

`Host` Link to the InfluxDB instance. eg: `http://link.to.influxdb:8086`

`Database` Name of the database as configured in InfluxDB

`Username` & `password` to connect to the InfluxDB database

`Url` Link to the public availability_map of Velo Antwerpen

### Docker

#### Build docker images
```
docker build -t velo:latest .
```

#### Run docker image
```
docker run -d velo
```
