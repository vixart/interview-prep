# 2.4. Слияние двух отсортированных связных списков

Раздел: `02_linked_lists / 09_merge_two_sorted_lists`

## Условие

Даны два отсортированных по возрастанию связных списка.
Необходимо слить их в один отсортированный список.

## Требования

- Сложность по времени `O(n + m)`.
- Сложность по памяти `O(1)` — не создавать новые узлы.

## Сигнатура

```go
type ListNode struct {
    Val  int
    Next *ListNode
}

func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode { return nil }
```

## Примеры

```
Input:  list1 = [1,2,4], list2 = [1,3,4]
Output: [1,1,2,3,4,4]

Input:  list1 = [], list2 = []
Output: []

Input:  list1 = [], list2 = [0]
Output: [0]
```
