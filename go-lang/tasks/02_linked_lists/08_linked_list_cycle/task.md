# 2.3. Поиск цикла в связном списке (Floyd's Cycle Detection)

Раздел: `02_linked_lists / 08_linked_list_cycle`

## Условие

Дан односвязный список. Необходимо определить, есть ли в нём цикл.

## Требования

- Реализовать алгоритм Floyd's Cycle Detection (медленный и быстрый указатели).
- Сложность по времени `O(n)`.
- Сложность по памяти `O(1)`.

## Сигнатура

```go
type ListNode struct {
    Val  int
    Next *ListNode
}

func hasCycle(head *ListNode) bool { return false }
```

## Примеры

```
Input:  head = [3,2,0,-4], pos = 1 (цикл начинается с узла со значением 2)
Output: true

Input:  head = [1,2], pos = 0 (цикл начинается с узла со значением 1)
Output: true

Input:  head = [1], pos = -1 (нет цикла)
Output: false
```
