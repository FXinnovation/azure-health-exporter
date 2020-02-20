# azure-health-exporter

Prometheus exporter exposing Azure [Resource health](https://docs.microsoft.com/en-us/azure/service-health/resource-health-overview) status as `up` metrics.

This exporter does not export metrics from Azure Monitor API, please use the [azure_metrics_exporter](https://github.com/RobustPerception/azure_metrics_exporter) for that. Note that the `azure-health-exporter` is following the same label naming used in the `azure_metrics_exporter`, in order to ease metric joins in PromQL queries.

## Getting Started

### Azure account requirements

This exporter call Azure API from an existing subscription with these requirements:

* An application must be registered (e.g., Azure Active Directory -> App registrations -> New application registration)

### API rate limit

Note that Azure imposes a very low rate limit for the Resource Health API calls. The Resource Health provider is returning an `X-Ms-Ratelimit-Remaining-Subscription-Resource-Requests` header, which means, as per [documentation](https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/request-limits-and-throttling#remaining-requests), that the "service has overridden the default limit". The limit has been observed to be 100 requests per 10 minutes. To minimize impact, the exporter has been designed to perform only one Resource Health request per scrape. The rate limit remaining count is also exposed as a metric.

### Prerequisites

To run this project, you will need a [working Go environment](https://golang.org/doc/install).

### Installing

```bash
go get -u github.com/FXinnovation/azure-health-exporter
```

## Building

Build the sources with

```bash
make build
```

## Run the binary

```bash
./azure-health-exporter
```

The exporter expects these environment variables to configure the Azure API
connection.

Environment Variable | Description
---------------------| -----------
AZURE_SUBSCRIPTION_ID | Found under properties in the Azure portal for your application/service
AZURE_TENANT_ID | Found under `Azure Active Directory > Properties` and listed as `Directory ID`
AZURE_CLIENT_ID | Also listed as `Application Id`, is obtained by registering an application under 'Azure Active Directory'
AZURE_CLIENT_SECRET | Is generated by selecting your application/service under Azure Active Directory, selecting 'keys', and generating a new key

By default, the exporter optional config file is expected in `config/config.yml`.

Use -h flag to list available options.

## Testing

### Running unit tests

```bash
make test
```

## Configuration

Configuration is usually done in `config/config.yml`.

An example can be found in
[config/config_example.yml](https://github.com/FXinnovation/azure-health-exporter/blob/master/config/config_example.yml).

Configuration element | Description
--------------------- | -----------
resource_configurations | (Mandatory) A list of configuration elements to select resources to monitor for health
resource_types | (Mandatory) A list of resource type to filter resources (must be part of the [supported type list](https://docs.microsoft.com/en-us/azure/service-health/resource-health-checks-resource-types))
resource_tags | (Mandatory) A map of resource tag name and value to filter resources
expose_azure_tag_info | (Optional, default to `false`) Whether or not to expose the `azure_tag_info` metric

## Docker image

You can run images published in [dockerhub](https://hub.docker.com/r/fxinnovation/azure-health-exporter).

You can also build a docker image using:

```bash
make docker
```

The resulting image is named `fxinnovation/azure-health-exporter:<git-branch>`.

The image exposes port 9613 and expects an optional config in `/opt/azure-health-exporter/config.yml`.
To configure it, you must pass the environment variables, and you can bind-mount a config from your host:

```bash
docker run -p 9613:9613 -v /path/on/host/config/config.yml:/opt/azure-health-exporter/config/config.yml -e AZURE_SUBSCRIPTION_ID="my_subscription_id" -e AZURE_TENANT_ID="my_tenant_id" -e AZURE_CLIENT_ID="my_client_id" -e AZURE_CLIENT_SECRET="my_client_secret" fxinnovation/azure-health-exporter:<git-branch>
```

## Exposed metrics

Metric | Description
------ | -----------
azure_resource_health_availability_up | [Resource health](https://docs.microsoft.com/en-us/azure/service-health/resource-health-overview) availability that relies on signals from different Azure services to assess whether a resource is healthy. This UP metric is 0 if availability status is `Unavailable`, and is 1 otherwise.
azure_tag_info | Tags of the Azure resource, exposed only if `expose_azure_tag_info` config is set to true
azure_resource_health_ratelimit_remaining_requests | Azure subscription scoped Resource Health requests remaining (based on `X-Ms-Ratelimit-Remaining-Subscription-Resource-Requests` header)

Example:

```
# HELP azure_resource_health_availability_up Resource health availability that relies on signals from different Azure services to assess whether a resource is healthy
# TYPE azure_resource_health_availability_up gauge
azure_resource_health_availability_up{resource_group="my_group",resource_name="my_name",resource_type="Microsoft.Storage/storageAccounts",subscription_id="xxx"} 1
# HELP azure_tag_info Tags of the Azure resource
# TYPE azure_tag_info gauge
azure_tag_info{resource_group="my_group",resource_name="my_name",resource_type="Microsoft.Storage/storageAccounts",subscription_id="xxx",tag_monitoring="enabled"} 1
# HELP azure_resource_health_ratelimit_remaining_requests Azure subscription scoped Resource Health requests remaining (based on X-Ms-Ratelimit-Remaining-Subscription-Resource-Requests header)
# TYPE azure_resource_health_ratelimit_remaining_requests gauge
azure_resource_health_ratelimit_remaining_requests{subscription_id="xxx"} 98
```

## Contributing

Refer to [CONTRIBUTING.md](https://github.com/FXinnovation/azure-health-exporter/blob/master/CONTRIBUTING.md).

## License

Apache License 2.0, see [LICENSE](https://github.com/FXinnovation/azure-health-exporter/blob/master/LICENSE).
