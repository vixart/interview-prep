# 4.3.6. Ещё один append: запись в чужую ёмкость

Раздел: `06_one_more_append`

Тип задачи: вопрос с собеседования «что выведет код и почему».

## Условие

```go
package main

import "fmt"

func main() {
    nums := []int{1, 2, 3}

    addNum(nums[0:2])
    fmt.Println(nums)

    addNums(nums[0:2])
    fmt.Println(nums)
}

func addNum(nums []int) {
    nums = append(nums, 4)
}

func addNums(nums []int) {
    nums = append(nums, 5, 6)
}
```

**Задание:** определи, что выведет программа (и скомпилируется ли она), и объясни почему.
