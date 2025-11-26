# OrChannel

`OrChannel` — пакет, для объединения нескольких `done`-каналов в один. Возвращаемый канал закрывается, как только любой из переданных каналов закрывается.

## Установка

```bash
go get github.com/v1adislav28/level4/orchannel
```


## Использование

```go
package main

import (
    "fmt"
    "time"
    "github.com/v1adislav28/level4/OrChannel"
)

func main() {
    sig := func(after time.Duration) <-chan interface{} {
        c := make(chan interface{})
        go func() {
            defer close(c)
            time.Sleep(after)
        }()
        return c
    }

    start := time.Now()
    <-OrChannel.Or(
        sig(2*time.Hour),
        sig(5*time.Minute),
        sig(1*time.Second),
        sig(1*time.Hour),
        sig(1*time.Minute),
    )
    fmt.Printf("done after %v\n", time.Since(start))
}
```

## Тесты

Запустить тесты:

```bash
go test 
```

## Примеры

См. `example.go` для примера использования.