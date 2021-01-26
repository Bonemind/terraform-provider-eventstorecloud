# 1.5.1 (Jan 26, 2021)

ENHANCEMENTS:

* Add missing server version
* Add missing instance types
* Fix ESC_ORG_ID environment variable key

# 1.5.0 (Dec 11, 2020)

ENHANCEMENTS:

* Publish to registry

# 1.4.0 (Dec 9, 2020)

ENHANCEMENTS:

* Azure support

FIXES:

* Increase the allowed CIDR range for networks and peerings to accommodate GCP and Azure

# 1.3.0 (Nov 13, 2020)

FIXES:

* refresh token was not properly passed to ESC client 

# 1.2.0 (September 30, 2020)

ENHANCEMENTS:

* `resource/eventstorecloud_managed_cluster`: `disk_size` can now be updated to expand the disk size of a cluster. ([#6](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/6))
* `resource/eventstorecloud_managed_cluster`: deleted managed clusters are correctly removed from state. ([#6](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/6))
* `resource/eventstorecloud_network`: deleted networks are correctly removed from state. ([#6](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/6))
* `resource/eventstorecloud_peering`: `provider_metadata` now contains `gcp_project_name` and `gcp_project_id`. `gcp_project_id` is suitable for passing to `google_compute_network_peering` resources. ([#6](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/6))
* `resource/eventstorecloud_peering`: deleted peerings are correctly removed from state. ([#6](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/6))

FIXES:

* ESC Client auth token directory is now created if it does not exist

## 1.1.0 (July 23, 2020)

ENHANCEMENTS:

* `resource/eventstorecloud_peering`: Add computed property `provider_metadata`. ([#4](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/4))
* `resource/eventstorecloud_managed_cluster`: Correct `three-node-multi-zone` topology validation. ([#2](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/2))
* `resource/eventstorecloud_managed_cluster`: Add `projection_level` field. ([#5](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/5))
* Meaningful error messages from the Event Store Cloud API are now presented. ([#5](https://github.com/EventStore/terraform-provider-eventstorecloud/pull/5))

## 1.0.0 (June 28, 2020)

FEATURES:

* Initial release of provider
* **New Data Source:** `eventstorecloud_project`
* **New Resource:** `eventstorecloud_project`
* **New Resource:** `eventstorecloud_network`
* **New Resource:** `eventstorecloud_peering`
* **New Resource:** `eventstorecloud_managedcluster`
