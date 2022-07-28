Understanding grid emissions across the country has been significant in the decarbonization movement. It’s especially relevant within facilities that are trying to reduce power consumption with corresponding carbon emissions. Since Atlas is providing them a method of power reduction, there should also exist a method of displaying emission reduction. Since we are not relying on a way to directly measure CO2 emissions, we can define a reduction in terms of site’s energy consumption and CO2 intensity data available through Singularity.   

Singularities Carbonara API provides all necessary data on fuel consumption rates/trends, carbon intensity rates/trends and the influx/outflux of energy across a region and within their power plants.
# carbon

build server
scripts/setup
run server: scripts/server

run client: go build -o bin/poller-cli github.com/crossnokaye/carbon/services/poller/cmd/poller-cli

call update using client: ./bin/poller-cli --url="grpc://localhost:12500" poller update
call uopdate using grpcurl:
brew install grpcurl
grpcurl -plaintext localhost:12500 poller.Poller.Update

Connect to clickhouse locally:
exec into docker container:
```bash
$ docker ps
CONTAINER ID   IMAGE                                 COMMAND                  CREATED       STATUS      PORTS                                                      NAMES
2698324fc48b   yandex/clickhouse-server:21.11.10.1   "/entrypoint.sh"         4 weeks ago   Up 6 days   0.0.0.0:8123->8123/tcp, 9009/tcp, 0.0.0.0:8088->9000/tcp   carbon_clickhouse
ec6f88377f32   redis:alpine                          "docker-entrypoint.s…"   6 weeks ago   Up 6 days   0.0.0.0:6379->6379/tcp                                     iam-redis


docker exec -it 2698324fc48b /bin/sh

```
connnect to clickhouse
```
# clickhouse-client --password atlas -u atlas

2698324fc48b :) select * from carbondb.carbon_reports
```