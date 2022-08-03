Calc Service

Description:

Historical Emissions Endpoint:


Deploying the Poller service in an Environment:

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
	 	grpcurl -plaintext -d '{"org_id": "52858b15-16ce-4998-b317-a1ce68c348c3", "facility_id": "a5746ffa-2073-455e-b811-322ad3c3c4b7", "location_id": "cf153258-c08f-4ff0-9b01-d51d452e40e5", "duration": [{"start_time": "2020-01-01T00:00:00Z", "end_time": "2020-01-02T00:00:00Z"}], "interval": "hourly"}' localhost:12200 calc.Calc.HistoricalCarbonEmissions

-max-time=1200 
