# Infinimesh plugins
We publish here plugins to connect infinimesh to external backends. To enable as much as possible modularization we have split the plugins into two main streams:  
* ### generic packages
   * [pkg](pkg)  
   pkg contains shared code to connect to our API, retrieve token and iterate over /objects to find devices in the desired namespace  
   * [redisstream](redistream)  
   shared code for generic cache and stream, based on redis. This package can be included into future plugins.
   
* ### plugins
   * [Elastic](Elastic)  
   Connect Infinimesh IoT seamless into [Elastic](https://elastic.co).
   * [Timeseries](timeseries)  
   [Redis-timeseries](https://oss.redislabs.com/redistimeseries/) with [Grafana](https://grafana.com/) for timeseries-analysis and rapid prototyping, can be used in production when configured as Redis cluster and ready to be hostet via [Redis-Cloud](https://redislabs.com/redis-enterprise-cloud/overview/). 
   * [SAPHana](SAPHana)  
   all code to connect infinimesh IoT Platform to any [SAP Hana](https://www.sap.com/products/hana.html) instance
   * [Snowflake](Snowflake)  
   all code to connect infinimesh IoT Platform to any [Snowflake](https://www.snowflake.com/) instance.  
   * [Cloud Connect](CloudConnect)  
   all code to connect infinimesh IoT Platform to Public Cloud Provider AWS, GCP and Azure. This plugin enables customer to use their own cloud infrastructure and extend infinimesh to other services, like [Scalytics](https://www.scalytics.io), using their own cloud native data pipelines and integration tools. 
  
More plugins will follow, please refer to the plugin directory for any developer friendly documentation.
  
## building plugins
checkout and build docker based environments starting in the / directory of plugins, like:  
```
git clone https://github.com/infinimesh/plugins.git  
cd plugins  
docker-compose -f timeseries/docker-compose.yml --project-directory . up --build
```
Please read the notes in the different plugin directories how to set ```username``` / ```password``` and API Endpoint (if not using [infinimesh.cloud](https://console.infinimesh.cloud)).
