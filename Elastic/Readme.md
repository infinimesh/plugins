## ElasticConnect
ElasticConnect is a cloud native connector to integrate Infinimesh IoT directly with [Elastic](https://elastic.co). With this plugin data engineers and data science community can directly stream IoT data into the processing engine to build AI driven IoT use-cases, like autonomous driving or moving pattern detection.

## Connector Setup

Simply set the relevant Infinimesh environment variables in `docker-compose.yml`:

```
USERNAME=CHANGEME
PASSWORD=CHANGEME
```

Then we can run everything from the repository root via:

```
docker-compose -f Elastic/docker-compose.yml --project-directory . up --build
```

## Developer notes

By default, this creates and writes documents into an index `infinimesh`. We can verify that the connector is working as expected using the following queries:

```
curl localhost:9200/_cat/indices
curl localhost:9200/infinimesh/_search?q=*:*
```
