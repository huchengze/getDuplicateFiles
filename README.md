## Quick start

Just for win10, use `Release` to get different architecture versions.

```sh
# Select the directory you want to detect, ep: C:\Users\Administrator\Desktop\getDuplicateFiles
.\bin\getDuplicateFiles C:\Users\Administrator\Desktop\getDuplicateFiles
```

## Release

```sh
make
```

## Test

```sh
# Make sure the directory in main_test.go(FindFilesInSameSize) is exist
go test
```
