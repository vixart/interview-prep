# 4.6.13. Версионирование документов (конкурентная обработка обновлений)

Раздел: `concurrency / 13_document_versioning`

## Условие

На вход сервису поступают обновления документов из очереди (Kafka).
Документы могут поступать в произвольном порядке и с дубликатами.

### Структура документа

- `Url` — уникальный идентификатор документа;
- `PubDate` — время заявляемой публикации документа;
- `FetchTime` — время получения обновления (идентификатор версии);
- `Text` — текст документа;
- `FirstFetchTime` — изначально `nil`, необходимо заполнить.

### Требования

Для группы документов с одинаковым `Url` необходимо:

1. `Text` и `FetchTime` брать из документа с МАКСИМАЛЬНЫМ `FetchTime` (последняя версия).
2. `PubDate` брать из документа с МИНИМАЛЬНЫМ `FetchTime` (первая версия).
3. `FirstFetchTime` установить равным минимальному `FetchTime`.

### Особенности

- Документы приходят в произвольном порядке (не по `FetchTime`).
- Возможны дубликаты сообщений.
- Если получили более свежую версию — обновляем результат.
- Если получили более старую версию — обновляем только `PubDate` и `FirstFetchTime`.
- Сервис работает в многопоточной среде (конкурентный доступ).

## Пример

```
Приходит документ 1: Url="doc1", FetchTime=100, Text="v1", PubDate=50
Результат: {Url="doc1", FetchTime=100, Text="v1", PubDate=50, FirstFetchTime=100}

Приходит документ 2: Url="doc1", FetchTime=200, Text="v2", PubDate=60
Результат: {Url="doc1", FetchTime=200, Text="v2", PubDate=50, FirstFetchTime=100}

Приходит документ 3: Url="doc1", FetchTime=50, Text="v0", PubDate=40
Результат: {Url="doc1", FetchTime=200, Text="v2", PubDate=40, FirstFetchTime=50}
```

## Заготовка

```go
package main

type Document struct {
    Url            string
    PubDate        uint64
    FetchTime      uint64
    Text           string
    FirstFetchTime *uint64
}

type Processor interface {
    // Process обрабатывает входящий документ и возвращает актуальную версию.
    // Возвращает nil, если документ является дубликатом и не требует обновления
    Process(doc Document) (*Document, error)
}
```
