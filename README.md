# Infinimesh plugins
We publish here plugins to connect infinimesh to external backends. To enable as much as possible modularization we have split the plugins into two main streams:  
* ### generic packages
   * [pkg](pkg)  
   pkg contains shared code to connect to our API, retrieve token and iterate over /objects to find devices in the desired namespace  
   * [redisstream](redistream)  
   shared code for generic cache and stream, based on redis. This package can be included into future plugins.
   
* ### plugins  
   * [Timeseries](timeseries)  
   redis-timeseries with grafana  
   * [SAPHana](SAPHana)  
   all code to connect infinimesh IoT platform to any SAPHana instance
   * [Snowflake](Snowflake)  
   all code to connect infinimesh IoT platform to any Snowflake instance using the Snowpipe API  
   * [Scalytics](Scalytics)  
   all code to connect infinimesh IoT Platform to [Scalytics](https://www.scalytics.io) using data ingest pipelines
  
More plugins will follow, please refer to the plugin directory for any developer friendly documentation.
  
## building plugins
checkout and build docker based environments starting in the / directory of plugins, like:  
```
docker-compose -f timeseries/docker-compose.yml --project-directory . up --build
```
Please read the notes in the different plugin directories
