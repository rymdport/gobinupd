# gobinupd
A tool to update installed tools in the `go/bin` folder.
It updates all your tools in parallel to speed up the update process.

## Install

To install `gobinupd`, run the following command:

```sh
go install github.com/rymdport/gobinupd@latest
```

Windows users might wish to run the tool directly, without installing it, to avoid having it try to update itself.

```sh
go run github.com/rymdport/gobinupd@latest
```

## Usage

Simply run `gobinupd` to automatically update all installed tools.
No path needs to be specified.

**Flags:**
- Pass `--release` to build without debug information for smaller binaries.
- Pass `--verbose` to print more information about the update process.
