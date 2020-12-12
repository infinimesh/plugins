## SAP Hana connector
We develop a SAP Hana plugin which connects Infinimesh directly with Hana. This plugin provides a redis cache for max 60 seconds to mitigate network latency.

### Developer notes  
Ideally we would want to use the SAP HANA JSON Docstore, but the SAP HANA cloud offering doesn't support it as of time of writing. Here we opt for an alternative where we store retrieved data as a JSON string in a table. Also see (https://answers.sap.com/questions/13207475/)
