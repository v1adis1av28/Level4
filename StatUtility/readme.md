# GC Metrics Exporter

Утилита предоставлящаюя метрики о памяти и сборке мусора через эндпоинт в формате Prometheus.

## Метрики

- `memstats_alloc` — объём выделенной и используемой памяти (Alloc).
- `total_allocation` — общее количество выделенной памяти с начала работы.
- `system_bytes` — объём памяти, полученной от ОС.
- `total_gc_counts` — количество выполненных сборок мусора.
- `last_gc_time_seconds` — время последней сборки мусора.

## Запуск

```
git clone https://github.com/v1adis1av28/level4
cd StatUtility
go run cmd/main.go
```

Сервер будет запущен на `http://localhost:8080`.

## Примеры запросов

- Получить метрики: [http://localhost:8080/metrics](http://localhost:8080/metrics)
- Открыть pprof: [http://localhost:8080/debug/pprof/](http://localhost:8080/debug/pprof/)
