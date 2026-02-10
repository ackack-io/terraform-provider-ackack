# Examples

This directory contains examples demonstrating how to use the ackack Terraform provider.

## Provider Configuration

See [provider/provider.tf](provider/provider.tf) for provider configuration options.

```hcl
provider "ackack" {
  # Set via ACKACK_API_KEY environment variable or directly
  api_key = "ak_your_api_key"
}
```

## Resources

- **[ackack_monitor](resources/ackack_monitor)** - Create uptime monitors (HTTP, DNS, SSL, TCP)
- **[ackack_alert](resources/ackack_alert)** - Configure alert notifications for monitors
- **[ackack_system](resources/ackack_system)** - Group monitors into logical systems
- **[ackack_report](resources/ackack_report)** - Generate uptime and incident reports

## Data Sources

- **[ackack_monitor](data-sources/ackack_monitor)** - Read a single monitor by ID
- **[ackack_monitors](data-sources/ackack_monitors)** - List all monitors

## Running the Examples

1. Set your API key:
   ```bash
   export ACKACK_API_KEY="ak_your_api_key"
   ```

2. Navigate to an example directory:
   ```bash
   cd examples/resources/ackack_monitor
   ```

3. Initialize and apply:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

## Document Generation

The document generation tool looks for files in the following locations by default:

* **provider/provider.tf** - example file for the provider index page
* **data-sources/`full data source name`/data-source.tf** - example file for the named data source page
* **resources/`full resource name`/resource.tf** - example file for the named resource page
* **resources/`full resource name`/import.sh** - example import command for the named resource page
