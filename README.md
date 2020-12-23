# Infinimesh plugins
We publish here plugins to connect infinimesh to external backends. To enable as much as possible modularization we have split the plugins into two main streams:  
* ### generic packages
   * [pkg](pkg)  
   pkg contains shared code to connect to our API, retrieve token and iterate over /objects to find devices in the desired namespace  
   * [redisstream](redistream)  
   shared code for generic cache and stream, based on redis. This package can be included into future plugins.
   
* ### Plugins and connectors
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
  
## Building plugins
checkout and build docker based environments starting in the / directory of plugins, like:  
```
git clone https://github.com/infinimesh/plugins.git  
cd plugins  
docker-compose -f timeseries/docker-compose.yml --project-directory . up --build
```
Please read the notes in the different plugin directories how to set ```username``` / ```password``` and API Endpoint (if not using [infinimesh.cloud](https://console.infinimesh.cloud)).  

## Deploy to any Kubernetes / OpenShift  
We recommend to use [kompose](https://kompose.io/) to translate the dockerfiles into kubernetes ready deployments. As example:  
```
# verify that it works via docker-compose  
docker-compose -f Elastic/docker-compose.yml --project-directory . up --build  
  
# convert to k8s yaml  
kompose -f Elastic/docker-compose.yml convert  
  
# prepare env - this makes sure that when we run `docker build` the image is accessible via minikube  
eval $(minikube docker-env)  
  
# build images and change the image name so that the k8s cluster doesn't try to pull it from some registry  
docker build -f ./redisstream/Dockerfile -t redisstream:0.0.1 . # change the image in producer-pod.yaml to redisstream:0.0.1  
docker build -f ./Elastic/Dockerfile -t elastic:0.0.1 . # change the image in consumer-pod.yaml to elastic:0.0.1  
  
# apply each yaml file
kubectl apply -f xxx.yaml  
  
# verify that it's working, eg via logs  
kubectl logs producer  

```

