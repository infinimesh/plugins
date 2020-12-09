# Infinimesh plugins
We publish here plugins to connect infinimesh to external backends. To enable as much as possible modularization we have split the plugins into two main streams:  
* pkg  
pkg contains shared code to connect to our API, retrieve token and iterate over /objects to find devices in the desired namespace  
* named plugins  
names plugins describing the external backends and system, divided into their respective name:  
..* timeseries  
redis-timeseries with grafana  
..* SAPHana
all code to connect infinimesh IoT platform to any SAPHana instance
..* Snowflake  
all code to connect infinimesh IoT platform to any Snowflake instance using the Snowpipe API  
  
More plugins will follow, please refer to the named plugins for any developer friendly documentation
