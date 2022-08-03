Poller Service

The Poller Service is meant to download carbon intensity reports from Singularity's Carbonara API and write them to clickhouse. It has two service methods to accomplish this: Update and GetEmissionsForRegion. Both will use Singularity's "Search endpoint" which will return carbon intensity data in series of 5 minute reports

Update Endpoint:

Update is used to "update" clickhouse with new carbon intensity data available between the last report stored(if any at all) in clickhouse and the current time. It will accomplish this for all regions where carbon intensity data is available through Singularity. The steps by which it downloads new carbon intensity reports per region are:
	1. check the date of the last report available in clickhouse. if no such date is available because no reports are found, use the default date available on Singularity's documentation. This is done by Ensure Past Data.

	2. Make API calls to Singularity to download carbon intensity data. Read that data into reports and store in clickhouse

	3. Create an array of dates(hours, days, weeks, and months) to obtain aggregate carbon intensity data. These dates will be created from the reports downloaded(in 5 minute intervals). For example the time stamps from singularity reports: {"2021-07-07 23:00:00 +0000 UTC", "2021-07-07 23:05:00 +0000 UTC", "2021-07-07 23:10:00 +0000 UTC" ... "2021-07-07 23:55:00 +0000 UTC"} would be converted into a date object with the start and end times {"2021-07-07 23:00:00 +0000 UTC", "2021-07-07 23:55:00 +0000 UTC"}. The same would be done for the days and months within given the time frame. This part is done by the function getDates.

	4. Those dates are used to query clickhouse for obtain aggregate data for generated, consumed, and marginal carbon intensity. Those reports are then saved in clickhouse. This part is done by the function getAggregateData

Get Emissions For Region:

This is used to return carbon intensity reports for a given time frame and region. It returns reports in 5 minute intervals

Environments:
Sandbox

Dependencies:
Clickhouse -> for storing Carbon Intensity reports
Carbonara API by Singularity -> For retrieving Carbon Intensity reports

Operation Analysis:
Calc service will not be able to compute power reports(KW) for facilities with power meters if carbon intensity reports are not correctly formatted in clickhouse

Updating the service:
Service configuration can be updated using standard Kubernetes configuration. Secrets are stored in AWS Secret Manager.

Metrics:
Singularity: Server and No data errors
Clickhouse: read and write errors

Deploying the Poller service in an Environment:

1. Configure a clickhouse user for each environment

2. Make sure that the secrets are stored in AWS secrets manager

3. Secrets either stored as carbon/poller for the poller service Singularity API or carbon/clickhouse for clickhouse

3. Connect to a environment cluster using ckutil cloud kube-connect <env name>

4. For each env use `.deploy carbon <branch_name> to <env_name>` in the corresponding slack channel

5. Since there is a large amount of carbon intensity reports to initially backfill when the service is deployed for the first time in an environment, the process below can be run to ensure that the the data is backfilled so that a cronjob won't expire before all reports are retrieved and written to clickhouse.

1. Get the pod that is currently running:
		kubectl get pods -n carbon:

		example response:

		NAME                      READY   STATUS    RESTARTS   AGE
		poller-5ff4565c7d-nzhjv   1/1     Running   0          2m50s


2. Run port-forward and get the logs:
		kubectl -n carbon port-forward {pod_id_from_above} 12500 &
		kubectl -n carbon logs -f {pod_id_from_above}


3. Make an api request to the pod in a new tab:
	 	grpcurl -plaintext -max-time=1200 localhost:12500 poller.Poller.Update


Testing the Poller service Locally:

1. Build server:
		scripts/setup

2. Run server:
		scripts/server
3. Run client:
		go build -o bin/poller-cli github.com/crossnokaye/carbon/services/poller/cmd/poller-cli

4. Call the method Update using client: 
		./bin/poller-cli --url="grpc://localhost:12500" poller update

5. Call update using grpcurl:
		1. brew install grpcurl
		2. grpcurl -plaintext localhost:12500 poller.Poller.Update

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

2698324fc48b :) select * from carbondb.carbon_reports

