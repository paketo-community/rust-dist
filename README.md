# Rust Dist Cloud Native Buildpack

The Rust Dist CNB provides the Rust stand-alone distribution. The buildpack
installs the Rust stand-alone distribution onto the $PATH which makes it available for
subsequent buildpacks. These buildpacks can then use that distribution to build
Rust projects. The Rust Cargo CNB is an example of a buildpack that utilizes the Rust
build tools.

## Integration

The Rust Dist CNB provides Rust as a dependency. Downstream buildpacks, like
[Rust Cargo CNB](https://github.com/dmikusa/rust-cargo-cnb) can require the rust
dependency by generating a [Build Plan
TOML](https://github.com/buildpacks/spec/blob/master/buildpack.md#build-plan-toml)
file that looks like the following:

```toml
[[requires]]

  # The name of the Rust dependency is "rust". This value is considered
  # part of the public API for the buildpack and will not change without a plan
  # for deprecation.
  name = "rust"

  # The version of the Rust dependency is not required. In the case it
  # is not specified, the buildpack will provide the latest stable version, 
  # which can be seen in the buildpack.toml file.
  version = "1.43.1"
```

## Usage

To package this buildpack for consumption:

```bash
$ ./scripts/package.sh
```

This builds the buildpack's Go source using GOOS=linux by default. You can supply another value as the first argument to package.sh.
