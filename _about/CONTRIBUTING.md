# Contributing

## Building and installing
In order to contribute to the provider, you have to build and install it manually. Make sure you have a supported version of Go installed and working. Check out or download this repository, then open a terminal and change to its directory.

### Installing the provider to `terraform.d/plugins`
The provider must be built/installed to your [plugins directory](https://www.terraform.io/docs/extend/how-terraform-works.html#plugin-locations) as appropriate.

Example Windows
```
$ go build -o C:\Users\me\AppData\Roaming\terraform.d\plugins\github.com\rymancl\dotcommonitor\1.0.0\windows_amd64/terraform-provider-dotcommonitor.exe
```

Example Linux
```
$ go build -o ~/.terraform.d/plugins/github.com/rymancl/dotcommonitor/1.0.0/darwin_amd64/terraform-provider-dotcommonitor
```

## Open issues
Any open and unassigned [issues](https://github.com/rymancl/terraform-provider-dotcommonitor/issues) are likely up for grabs. Feel free to ping the [maintainers](./MAINTAINERS.md) to confirm the status of an issue or if you have any question.
