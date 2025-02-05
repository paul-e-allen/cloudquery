# Table: gcp_compute_addresses

This table shows data for GCP Compute Addresses.

https://cloud.google.com/compute/docs/reference/rest/v1/addresses#Address

The primary key for this table is **self_link**.

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_source_name|String|
|_cq_sync_time|Timestamp|
|_cq_id|UUID|
|_cq_parent_id|UUID|
|project_id|String|
|address|String|
|address_type|String|
|creation_timestamp|String|
|description|String|
|id|Int|
|ip_version|String|
|ipv6_endpoint_type|String|
|kind|String|
|name|String|
|network|String|
|network_tier|String|
|prefix_length|Int|
|purpose|String|
|region|String|
|self_link (PK)|String|
|status|String|
|subnetwork|String|
|users|StringArray|