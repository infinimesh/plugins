## CloudConnect
We develop a cloud native connector to integrate Infinimesh IoT directly into Cloud Native Storage of AWS, GCP and Azure. This plugin enables data engineers and the data science community to directly store IoT data into the data storage solutions, like S3. This plugin works seamless with [Scalytics AI](https://scalytics.io)

## Connector Setup

Simply set the required infinimesh and cloud variables in `docker-compose.yml`.

Then we can simply run:

```
docker-compose -f CloudConnect/docker-compose.yml --project-directory . up --build
```

## Developer Notes

The plugin comes in two parts: (1) writing data to local csv and (2) uploading the csv file to the chosen cloud storage. These two functions are implemented as two separate services.

When writing to local csv, files are sharded by device uids and timestamps (eg. `./<device_id>/<timestamp>.csv`). Timestamp here refers to the number of seconds past unix epoch, round to the next minute. This means that all events for the same device within the same minute timeframe will be written to the same file for bulk uploading.

After writing to local csv, a separate service walks the directory and searches for csv files that were not recently modified via each file's last modified time. If an old file is detected, the file gets uploaded to the relevant cloud, and deleted if the operation is successful.
