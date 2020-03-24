# go-check
Check for unused or updatable golang modules.

go-check runs go list -u -m -json all and displays the updates for your modules.

## Installation
```
go get -u github.com/vincegio/go-check
```
## Usage
go-check should be run inside a folder with go modules set up.
### Updates
```
go-check updates [flags]

Flags:
  -d, --direct        Direct packages
  -h, --help          help for updates
  -u, --interactive   Interactive update

Global Flags:
  -v, --verbose   verbose output
```
