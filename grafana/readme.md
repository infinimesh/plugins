# Grafana backend plugin

This plugin lets your grafana instance connect with infinimesh via our API. We have Grafana builtin in this container, but you can use external services, too. We use Redis Timeseries as backend, the service connects to our stream API, iterate over /objects and build a timeseries structure for every device found. This plugin can be used as an blueprint, or hosted somewhere. We fully comply to GDPR and CCPA, and our users get full data ownership.  

## Connector Setup

Simply set the username, password and api url environment variables in `docker-compose.yml`:
```
USERNAME=Your Infinimesh User Name  
PASSWORD=Your Infinimesh Password
API_URL=FQDN API Endpoint (like https://api.infinimesh.cloud)
```
The API Endpoint must be reachable via DNS, and the CA must be known.  
  
Then run everything via:

```
docker-compose up --build
```

## Grafana Setup

First visit `localhost:3000` and sign in with the default grafana admin credentials (username=admin and password=admin). Follow the instructions to change the password accordingly.

Next add the redis data source by clicking under Configuration -> Data Sources -> Add Data Source -> Redis. Change the address setting to `redis://redis:6379`, and click on `Save and Test`.

If the above works successfully, we should be able to begin visualizing some data. To get started, we can import the sample dashboard provided in this repository (`sample-dashboard.json`) by clicking on Create -> Import -> Upload JSON File.

## Redis Timeseries

We are using [Redislab Timeseries](https://oss.redislabs.com/redistimeseries/), for any documentation how to build dashboards using the redislab timeseries datasource please visit the link.
