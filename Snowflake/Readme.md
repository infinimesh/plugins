## Snowflake DB connector
We develop a Snowflake DB plugin which connects to Infinimesh directly and writes data directly into a Snowflake table.

## Connector Setup

Simply set the relevant Infinimesh and Snowflake environment variables in `docker-compose.yml`:

```
# infinimesh env variables
USERNAME=CHANGEME
PASSWORD=CHANGEME

# snowflake env variables
SNOWFLAKE_INITDB=true
SNOWFLAKE_ACCOUNT=CHANGEME
SNOWFLAKE_USER=CHANGEME
SNOWFLAKE_PASSWORD=CHANGEME
```

`SNOWFLAKE_INITDB` is an option for creating the required set of relations when the service starts. Note that the operations here may not be optimal and further tuning may be required

Then we can run everything from the repository root via:

```
docker-compose -f Snowflake/docker-compose.yml --project-directory . up --build
```

## Developer notes

Device states are inserted into a SQL table `infinimesh.device_states` with the following schema

```
uid VARCHAR(1000)
timestamp TIMESTAMP
version VARCHAR(1000)
data VARIANT
PRIMARY KEY (uid, timestamp)
```

Note that `data` is a JSON type that holds all device data, represented as `map[string]interface{}` in code
