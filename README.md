# `docker.io/paketo-community/rust-dist`

The Rust Dist Cloud Native Buildpack provides a Rust toolchain from a zip archive distribution. The buildpack installs the Rust toolchain onto the `$PATH` which makes it available for subsequent buildpacks to consume. Subsequent buildpacks can then use the toolchain to build Rust projects. The Rust Cargo CNB is an example of a buildpack that utilizes a Rust toolchain.

## Behavior

This buildpack will always pass detection.

This buildpack will participate during build if any of the following conditions are met

* Another buildpack requires `rust`

The buildpack will do the following if Rust is requested:

* Contributes Rust to a layer marked `build` and `cache` with all commands on `$PATH`

## Configuration

| Environment Variable | Description                                                                                                                                                                                                                    |
| -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `$BP_RUST_VERSION`   | As a user of the buildpack, you may specify which version of Rust gets installed by setting this environment variable at build time. The version you specify must exist in the `buildpack.toml` file or you will get an error. |

## Bindings
The buildpack optionally accepts the following bindings:

### Type: `dependency-mapping`
| Key                   | Value   | Description                                                                                       |
| --------------------- | ------- | ------------------------------------------------------------------------------------------------- |
| `<dependency-digest>` | `<uri>` | If needed, the buildpack will fetch the dependency with digest `<dependency-digest>` from `<uri>` |

## Usage

In general, [you probably want the rust CNB instead](https://github.com/paketo-community/rust/#tldr). 

If you want to use this particular CNB directly, the easiest way is via image. Run `pack build -b paketocommunity/rust-dist:<version> ...`.

## License
This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0
