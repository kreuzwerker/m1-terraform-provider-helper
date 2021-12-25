# m1-terraform-provider-helper

A CLI to help with managing the installation and compilation of terraform providers when running a new M1 Mac. 

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
  status      Shows the status of the m1 provider installations

Flags:
  -h, --help   help for m1-terraform-provider-helper
```

Example:
You want to install the `terraform-provider-vault` in version `v2.10.0` because you are using it in a project and let's assume it has no pre-build binary for Mac M1:

```sh
m1-terraform-provider-helper activate # (In case you have not activated the helper)
m1-terraform-provider-helper install hashicorp/vault -v v2.10.0 # Install and compile
```

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
