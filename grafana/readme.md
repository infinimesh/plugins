# Grafana backend plugin

This plugin lets your grafana instance connect with infinimesh via our API.

## Connector Setup

Simply set the username, password and api url environment variables in `docker-compose.yml`, then run everything via:

```
docker-compose up --build
```

## Grafana Setup

First visit `localhost:3000` and sign in with the default grafana admin credentials (username=admin and password=admin). Follow the instructions to change the password accordingly.

Next add the redis data source by clicking under Configuration -> Data Sources -> Add Data Source -> Redis. Change the address setting to `redis://redis:6379`, and click on `Save and Test`.

If the above works successfully, we should be able to begin visualizing some data. To get started, we can import the sample dashboard provided in this repository (`sample-dashboard.json`) by clicking on Create -> Import -> Upload JSON File.
