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

A CLI to manage the installation and compilation of Terraform providers for an ARM-based Mac.

## Table of Contents

- [m1-terraform-provider-helper](#m1-terraform-provider-helper)
  - [Table of Contents](#table-of-contents)
  - [Motivation](#motivation)
  - [Installation](#installation)
    - [Dependencies](#dependencies)
  - [Usage](#usage)
    - [Debugging Installation Problems](#debugging-installation-problems)
    - [Terraform Lockfile handling](#terraform-lockfile-handling)
    - [Providing custom build commands](#providing-custom-build-commands)
    - [Providing custom provider repository](#providing-custom-provider-repository)
    - [Logging](#logging)
    - [Timeouts](#timeouts)
    - [Plugin Directory](#plugin-directory)
  - [Development](#development)
    - [Testing](#testing)
    - [Build](#build)
    - [Release](#release)
  - [License](#license)
## Motivation

While using my then-new MacBook with an M1 chip, I often encountered issues in client projects when
working with Terraform. Some Terraform providers hadn't adapted to the new `darwin_arm64` architecture at
all, or else the provider was pinned to an older, incompatible version. In both cases, there was no
pre-compiled binary for `darwin_arm64`; you had to compile it yourself. (There's a nice write-up on how to
compile in a
[Terraform Issue](https://github.com/hashicorp/terraform/issues/27257#issuecomment-754777716).) As I was
constantly switching back and forth between own-compiled binaries and pre-built ones, I wanted an elegant
solution that managed all the details by itself.

## Installation

```sh
brew install kreuzwerker/taps/m1-terraform-provider-helper
```

### Dependencies

Since Go is used to build the providers, you need to have a working Go setup in the local directory where you run m1-terraform-provider-helper
commands. Although Go is installed by Homebrew as a dependency
of m1-terraform-provider-helper, the Go binary won't necessarily be in your PATH. (For example, if you use
asdf or a similar version manager for Go, the version manager's shim likely comes before Homebrew's Go
binary in your PATH.) Ensure that the command `go version` succeeds before using this tool.

## Usage

```
A CLI to manage the installation of Terraform providers for an ARM-based Mac

Usage:
  m1-terraform-provider-helper [command]

Available Commands:
  activate    Activate the m1-terraform-provider-helper
  completion  Generate the autocompletion script for the specified shell
  deactivate  Deactivate the m1-terraform-provider-helper
  help        Help about any command
  install     Download (and compile) a Terraform provider for an ARM-based Mac
  list        List all available providers and their versions
  lockfile    Commands to work with Terraform lockfiles
  status      Show the status of the m1-terraform-provider-helper installation
  version     Display the current version

Flags:
  -h, --help   help for m1-terraform-provider-helper

Use "m1-terraform-provider-helper [command] --help" for more information about a command.
```

Example:
You want to install version `v2.10.0` of `terraform-provider-vault` because you're using it in a project. Let's assume it has no pre-built binary for an ARM-based Mac:

```sh
m1-terraform-provider-helper activate # (In case you have not activated the helper)
m1-terraform-provider-helper install hashicorp/vault -v v2.10.0 # Install and compile
```

### Debugging Installation Problems

The `install` commands relies on an internal `buildCommands` map to find the correct build command for an provider. If the command is not correct, you can provide a custom build command by using the `--build-command` flag. See [Providing custom build commands](#providing-custom-build-commands) for more details.
In order to find the correct build command, please take a look at the documentation of the provider you are trying to install.

The `m1-terraform-provider-helper` downloads the source code of the provider to `$HOME/.m1-terraform-provider-helper`, which means you can actually play around with the source code and try to compile it yourself.

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

### Providing custom provider repository
You can override the built-in querying mechanism of the terraform registry by using the `--custom-provider-repository-url` flag.

**Explanation**:
The `install` commands relies on an internal queries the default terraform registry url (which you can also override), to 
determine the url of the git repository of the desired provider. However, for some providers there is no url 
as they are, e.g. already *archived*. 

For example for the mysql provider the command would be
```sh
m1-terraform-provider-helper install hashicorp/terraform-provider-mysql -v v1.9.0 --custom-provider-repository-url "https://github.com/hashicorp/terraform-provider-mysql"
```

### Logging

You can enable additional log output by setting the `TF_HELPER_LOG` environment variable to `info` or `debug` log level.

### Timeouts

The `m1-terraform-provider-helper` does make HTTP calls to the terraform provider registry. The default timeout is 10 seconds. You can change that timeout by using the `TF_HELPER_REQUEST_TIMEOUT` environment variable. For example `TF_HELPER_REQUEST_TIMEOUT=15` for a timeout of 15 seconds.

### Plugin Directory

The destination and name of the compiled provider depends on the terraform version:

* For Terraform `<0.13` it is `~/.terraform.d/plugins/darwin_arm64/terraform-provider-template_v2.2.0` (based on https://developer.hashicorp.com/terraform/language/v1.1.x/configuration-0-11/providers#plugin-names-and-versions)
* For all Terraform versions `>=0.13` it is `~/.terraform.d/plugins/registry.terraform.io/${providerName}/${version}/darwin_arm64/terraform-provider-template_2.2.0_x5` (based on https://developer.hashicorp.com/terraform/cli/config/config-file#implied-local-mirror-directories)

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
