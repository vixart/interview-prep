# 2.5. Удаление дубликатов из отсортированного связного списка

Раздел: `02_linked_lists / 10_remove_duplicates_sorted_list`

## Условие

Дан отсортированный по возрастанию связный список.
Необходимо удалить все дубликаты так, чтобы каждый элемент встречался только один раз.

## Требования

- Сложность по времени `O(n)`.
- Сложность по памяти `O(1)`.

## Сигнатура

```go
type ListNode struct {
    Val  int
    Next *ListNode
}

func deleteDuplicates(head *ListNode) *ListNode { return nil }
```

## Примеры

```
Input:  head = [1,1,2]
Output: [1,2]

Input:  head = [1,1,2,3,3]
Output: [1,2,3]

Input:  head = [1,1,1,1]
Output: [1]
```
