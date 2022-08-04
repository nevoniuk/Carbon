Calc Service

Description:

Historical Emissions Endpoint:


Deploying the Calc service in an Environment:

1. Configure a clickhouse user for each environment

2. Make sure that the secrets are stored in AWS secrets manager(calc service only uses clickhouse secrets)

3. Connect to a environment cluster using ckutil cloud kube-connect janeway

4. For each env use `.deploy carbon <branch_name> to <env_name>` in the corresponding slack channel

Make a client call to calc service to test in the given env:

1. Get the pod that is currently running:
		kubectl get pods -n carbon:

		example response:

		NAME                      READY   STATUS    RESTARTS   AGE
		calc-5ff4565c7d-nzhjv   1/1     Running   0          2m50s


2. Run port-forward and get the logs:
		kubectl -n carbon port-forward {pod_id_from_above} 12200 &
		kubectl -n carbon logs -f {pod_id_from_above}


3. Make an api request to the pod in a new tab:
	 	grpcurl -plaintext -d '{"org_id": "52858b15-16ce-4998-b317-a1ce68c348c3", "facility_id": "a5746ffa-2073-455e-b811-322ad3c3c4b7", "location_id": "cf153258-c08f-4ff0-9b01-d51d452e40e5", "duration": [{"start_time": "2020-09-23T00:00:00-00:00", "end_time": "2020-09-24T00:00:00-00:00"}], "interval": "hourly"}' localhost:12200 calc.Calc.HistoricalCarbonEmissions
-max-time=1200 
"2006-01-02T15:04:05-07:00"

Testing the Calc service Locally:

1. Build server:
		scripts/setup

2. Run server:
		scripts/server
3. Run client:
		go build -o bin/calc-cli github.com/crossnokaye/carbon/services/calc/cmd/calc-cli

4. Call the method HistoricalCarbonEmissions using client: 
		./bin/calc-cli --url="grpc://localhost:12200" poller update

5. Call update using grpcurl:
		1. brew install grpcurl
		2. grpcurl -plaintext -d '{"org_id": "52858b15-16ce-4998-b317-a1ce68c348c3", "facility_id": "a5746ffa-2073-455e-b811-322ad3c3c4b7", "location_id": "cf153258-c08f-4ff0-9b01-d51d452e40e5", "duration": [{"start_time": "2020-01-01T00:00:00Z", "end_time": "2020-01-02T00:00:00Z"}], "interval": "hourly"}' localhost:12200 calc.Calc.HistoricalCarbonEmissions

Connect to clickhouse locally to ensure that carbon intensity reports were written:

1. Exec into docker container:

	run docker ps

CONTAINER ID   IMAGE                                 COMMAND                  CREATED       STATUS      PORTS                                                      NAMES
2698324fc48b   yandex/clickhouse-server:21.11.10.1   "/entrypoint.sh"         4 weeks ago   Up 6 days   0.0.0.0:8123->8123/tcp, 9009/tcp, 0.0.0.0:8088->9000/tcp   carbon_clickhouse
ec6f88377f32   redis:alpine                          "docker-entrypoint.s…"   6 weeks ago   Up 6 days   0.0.0.0:6379->6379/tcp                                     iam-redis

then

docker exec -it 2698324fc48b /bin/sh

3. connnect to clickhouse

# clickhouse-client --password atlas -u atlas

4. query for carbon intensity reports

2698324fc48b :) select * from carbondb.carbon_intensity_reports
