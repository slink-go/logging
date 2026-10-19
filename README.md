# slink-go/logging
A simplifying wrapper for zerolog.

Import:
```shell
go get github.com/slink-go/logging@v0.0.2
```
Use:
```go
  l := logging.GetLogger("logger-name")
  l.Warn("some %s message %d", "test", 1)
```

## Лицензия

Двойная лицензия:

- [AGPL-3.0-only](LICENSE) — свободное использование, copyleft при распространении;
- [Коммерческая лицензия](LICENSE.COMMERCIAL) — использование без copyleft-обязательств: kami@slink.ws.
