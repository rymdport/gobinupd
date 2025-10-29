# gobinupd
A tool to update installed tools in the `go/bin` folder.

## Install

To install `gobinupd`, run the following command:

```sh
go install github.com/rymdport/gobinupd@latest
```

## Usage

Simply run `gobinupd` to automatically update all installed tools.
- Optionally, pass `--release` to build without debug information for smaller binaries.
- Optionally, pass `--verbose` to print more information about the update process.
