## Quick start

Just for win10, use `Release` to get different architecture versions.

```sh
# Select the directory you want to detect, ep: C:\Users\huchengze\Desktop\getDuplicateFiles
.\bin\getDuplicateFiles C:\Users\huchengze\Desktop\getDuplicateFiles
```

## Release

```sh
make
```

## Test

```sh
# Make sure the directory in main_test.go(FindFilesInSameSize) exists
go test
```
