## SAP Hana connector
We develop a SAP Hana plugin which connects Infinimesh directly with Hana. This plugin provides a redis cache for max 60 seconds to mitigate network latency.

## Connector Setup

Simply set the relevant Infinimesh and SAP HANA environment variables in `docker-compose.yml`:

```
# infinimesh env variables
USERNAME=CHANGEME
PASSWORD=CHANGEME

# saphana env variables
SAPHANA_INITDB=true
SAPHANA_USERNAME=CHANGEME
SAPHANA_PASSWORD=CHANGEME
SAPHANA_HOST=CHANGEME # on SAP HANA cloud, this looks like instanceid.hana.trial-region.hanacloud.ondemand.com
SAPHANA_PORT=443
```

`SAPHANA_INITDB` is an option for creating the required set of relations when the service starts. Note that the operations here may not be optimal and further tuning may be required

Then we can run everything from the repository root via:

```
docker-compose -f SAPHana/docker-compose.yml --project-directory . up --build
```

## Developer notes

Device states are inserted into a SQL table `devices` with the following schema

```
uid varchar(1000)
timestamp timestamp
version varchar(1000)
data varchar(5000)
PRIMARY KEY (uid, timestamp)
```

`data` is a JSON string that holds all device data, which is represented as `map[string]interface{}` in code. Ideally we would prefer to insert data into the SAP HANA JSON Docstore over a SQL table, but the SAP HANA cloud offering doesn't support docstore as of time of writing. Also see [this](https://answers.sap.com/questions/13207475/)
