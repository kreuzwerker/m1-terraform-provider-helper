<a href="https://terraform.io">
    <img src="https://raw.githubusercontent.com/kreuzwerker/m1-terraform-provider-helper/main/assets/terraform-logo.png" alt="Terraform logo" title="Terraform" align="right" height="100" />
</a>
<a href="https://kreuzwerker.de">
    <img src="https://raw.githubusercontent.com/kreuzwerker/m1-terraform-provider-helper/main/assets/xw-logo.png" alt="Kreuzwerker logo" title="Kreuzwerker" align="right" height="100" />
</a>

# m1-terraform-provider-helper

[![Release](https://img.shields.io/github/v/release/kreuzwerker/m1-terraform-provider-helper)](https://github.com/kreuzwerker/m1-terraform-provider-helper/releases)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/kreuzwerker/m1-terraform-provider-helper/blob/main/LICENSE)  
[![Go Status](https://github.com/kreuzwerker/m1-terraform-provider-helper/workflows/tests/badge.svg)](https://github.com/kreuzwerker/m1-terraform-provider-helper/actions)
[![Lint Status](https://github.com/kreuzwerker/m1-terraform-provider-helper/workflows/golangci-lint/badge.svg)](https://github.com/kreuzwerker/m1-terraform-provider-helper/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/kreuzwerker/m1-terraform-provider-helper)](https://goreportcard.com/report/github.com/kreuzwerker/m1-terraform-provider-helper)  

A CLI to help with managing the installation and compilation of terraform providers when running a new M1 Mac. 

## Table of Contents

- [m1-terraform-provider-helper](#m1-terraform-provider-helper)
  - [Table of Contents](#table-of-contents)
  - [Motivation](#motivation)
  - [Installation](#installation)
  - [Usage](#usage)
    - [Terraform Lockfile handling](#terraform-lockfile-handling)
    - [Providing custom build commands](#providing-custom-build-commands)
    - [Logging](#logging)
  - [Development](#development)
    - [Testing](#testing)
    - [Build](#build)
    - [Release](#release)
  - [License](#license)
## Motivation

While using my new Macbook with M1 chip I often encountered issues in client projects when working with Terraform projects. Either some terraform providers have no adapted to the new `darwin_arm64` at all or the version of the used provider is pinned to an older version. In both cases, there is no pre-compiled binary for `darwin_arm64` => you have to compile it yourself. There is a nice writeup on how to compile in a [Terraform Issue](https://github.com/hashicorp/terraform/issues/27257#issuecomment-754777716). As I am constantly switching forth and back between using own-compiled binaries and pre-build, I wanted to have an elegant solution which manages all the details by itself.

## Installation

```sh
brew install kreuzwerker/taps/m1-terraform-provider-helper
```

## Usage

```
A CLI to manage the installation of terraform providers for the Mac M1 chip

Usage:
  m1-terraform-provider-helper [command]

Available Commands:
  activate    Activate the usage of m1 provider helper
  deactivate  Deactivate the usage of M1 provider helper
  help        Help about any command
  install     Downloads (and compiles) a terraform provider for the M1 chip
  list        Lists all available providers and their versions
  lockfile    Commands to work with terraform lockfiles
  status      Shows the status of the m1 provider installations
  version     Display the current version

Flags:
  -h, --help   help for m1-terraform-provider-helper
```

Example:
You want to install the `terraform-provider-vault` in version `v2.10.0` because you are using it in a project and let's assume it has no pre-build binary for Mac M1:

```sh
m1-terraform-provider-helper activate # (In case you have not activated the helper)
m1-terraform-provider-helper install hashicorp/vault -v v2.10.0 # Install and compile
```

### Terraform Lockfile handling

tl;dr: Use `m1-terraform-provider-helper lockfile upgrade` to add the checksum of all used local providers to your projects `.terraform.lock.hcl`. Use the `--help` flag to see all available options for specifying input and output directories.

Most Terraform projects have a `.terraform.lock.hcl` file for pinning depedencies (https://www.terraform.io/language/files/dependency-lock). When using the `m1-terraform-provider-helper` and installing a provider locally, all following `terraform init` commands will lead to an error:

```
Error: Failed to install provider

Error while installing hashicorp/azurerm v2.1.0: the current package for
registry.terraform.io/hashicorp/azurerm 2.1.0 doesn't match any of the
checksums previously recorded in the dependency lock file.
```

The reason is that the checksums inside the existing lockfile are the checksum of the previously installed `darwin_amd64` provider. Now we are using our own `darwin_arm64` compiled provider, which has a different checksum. In order to make `terraform init` work again, we have to add the checksum of the local provider to the lockfile.

This is done via the `m1-terraform-provider-helper lockfile upgrade` command. It also two flags which you can use to specify the input/output lockfile:
* `--input-lockfile-path` 
* `--output-path`


### Providing custom build commands

You can override the built-in build command handling by using the `--custom-build-command` flag.

**Explanation**:
The `install` commands relies on an internal `buildCommands` map to find the correct build command for an provider. For some important providers we have hardcoded different commands, but the default (and fallback) is `make build`. If that does not work for the provider you want to install, you can also pass a custom build command using the `--custom-build-command` flag.

Please refer to the documentation of the provider to find out the build command.

### Logging

You can enable additional log output by setting the `TF_HELPER_LOG` environment variable to `info` or `debug` log level.

## Development

### Testing

To run tests execute:

```sh
make test
make lint
```

in the project's root directory.

### Build

To build the app execute:

```sh
make build
```

in the project's root directory. This will generate the executable `dist/m1-terraform-provider-helper` file that you can run.

### Release

**IMPORTANT**: Before releasing any version, you have to manually edit the `cmd/version.go` file and change the `version` constant to the new version you'll release.

If you want to generate the changelog and see it only (it will neither commit, tag nor push)
run one of the following commands:

```sh
make patch
make minor
make major
```

If you want it automated prepend `TAG=1` to the command as follows:

```sh
# TAG=1 indicates to tag and generate the changelog
TAG=1 make minor
git push origin main --tags 
```

## License

Distributed under the MIT License. See `LICENSE.txt` for more information.
