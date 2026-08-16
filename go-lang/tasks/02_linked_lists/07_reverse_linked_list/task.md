# 2.2. Разворот связного списка

Раздел: `02_linked_lists / 07_reverse_linked_list`

## Условие

Дан односвязный список. Необходимо развернуть его.

## Требования

- Реализовать итеративное решение (дополнительно — рекурсивный вариант).
- Сложность по времени `O(n)`.
- Сложность по памяти `O(1)`.

## Сигнатура

```go
type ListNode struct {
    Val  int
    Next *ListNode
}

func reverseList(head *ListNode) *ListNode { return nil }
```

## Примеры

```
Input:  head = [1,2,3,4,5]
Output: [5,4,3,2,1]

Input:  head = [1,2]
Output: [2,1]

Input:  head = []
Output: []
```
