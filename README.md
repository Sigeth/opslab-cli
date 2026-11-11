# opslab-cli

The `opslab` command line of opslab (working name):

- **Operators:** install and manage the platform on a Linux machine or cluster (`opslab install`, from M2).
- **Challenge authors:** write, test and publish challenges (`opslab challenge …`).

> Status: M1. Available today: `opslab version`, `opslab challenge lint`, `opslab challenge manifest`.

## Usage

```sh
opslab challenge lint path/to/challenge   # same validation as the operator at install time
opslab challenge manifest path/to/challenge | kubectl apply -f -   # install a local challenge (development, unsigned)
```

## Development

Requirements: Go (version in `go.mod`), make.

```sh
make lint test build   # bin/opslab (linux/amd64)
make install           # into $GOBIN
```

## License

[MIT](LICENSE)
