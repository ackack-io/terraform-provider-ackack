# Terraform Provider: ackack
This is the official Terraform provider for [ackack.io](https://ackack.io).

This provider requires an API key to use. Don't have one? 
[Sign up for free](https://ackack.io/signup) and get 3 free monitors with Slack and Discord alerts!

## Using the provider

Fill this in for each provider

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24
- [GNU Make](https://www.gnu.org/software/make/)
- [golangci-lint](https://golangci-lint.run/usage/install/#local-installation) (optional)

## Building

1. Clone the repository
1. Enter the repository directory
1. `make build`

## Development

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
