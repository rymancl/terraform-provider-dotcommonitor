# Terraform provider for Dotcom-Monitor

![Status: Tech Preview](https://img.shields.io/badge/status-experimental-yellow) 
[![Releases](https://img.shields.io/github/v/release/rymancl/terraform-provider-dotcommonitor.svg)](https://github.com/rymancl/terraform-provider-dotcommonitor/releases) [![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg)](https://github.com/semantic-release/semantic-release)

The Terraform provider for [Dotcom-Monitor](https://www.dotcom-monitor.com) makes it easy to ensure performance, functionality, and uptime of websites, web applications, servers, and APIs. Using this provider allows you to follow a "monitoring-as-code" observability approach.

## Quick Links
* [Provider Documentation](https://registry.terraform.io/providers/rymancl/dotcommonitor/latest/docs)
* [Dotcom-Monitor API](https://wiki.dotcom-monitor.com/knowledge-base/getting-started-with-the-api)

## Development & Releases
This provider is under active development. **Feature enhancement releases that contain breaking changes should be expected.** Once `v1.0.0` is released, standard semantic versioning will be followed in regards to the introduction of breaking changes.

## Contributing
Contributions are welcomed! Please understand that the experimental nature of this repository means that contributing code may be a bit of a moving target. If you have an idea for an enhancement or bug fix, and want to take on the work yourself, please first [create an issue](https://github.com/rymancl/terraform-provider-dotcommonitor/issues/new) so that we can discuss the implementation with you before you proceed with the work.

Please review the [contribution guide](_about/CONTRIBUTING.md) to begin.

## Requirements
* [Terraform](https://www.terraform.io/downloads.html) >=0.13
* [Go](https://golang.org/doc/install) >=1.16 (to build the provider plugin)
