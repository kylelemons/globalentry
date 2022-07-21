# Usage

For remote appointments use a command like this:

```shell
go run github.com/kylelemons/globalentry/cmd/schedulewatcher --remote
```

If you have a location code (either from Chrome Developer Tools or the [code list]):

```shell
go run github.com/kylelemons/globalentry/cmd/schedulewatcher --location=5446 # for San Francisco
```

You can also watch for both at the same time:

```shell
go run github.com/kylelemons/globalentry/cmd/schedulewatcher --remote --location=5446
```
