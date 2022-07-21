# Usage

For remote appointments use a command like this:

```shell
go run github.com/kylelemons/globalentry/cmd/schedulewatcher@latest --remote
```

If you have a location code (either from Chrome Developer Tools or the [code list]):

```shell
go run github.com/kylelemons/globalentry/cmd/schedulewatcher@latest --location=5446 # for San Francisco
```

You can also watch for both at the same time:

```shell
go run github.com/kylelemons/globalentry/cmd/schedulewatcher@latest --remote --location=5446
```

You can also install the tool and run it locally

```shell
go install github.com/kylelemons/globalentry/cmd/schedulewatcher@latest
schedulewatcher --remote --location=5446 --every=1h
```