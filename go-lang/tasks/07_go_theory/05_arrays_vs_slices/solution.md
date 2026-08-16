# Решение

## Разбор

1. Предсказать вывод всех `Println`.
2. Объяснить:
   - массив — значимый тип, при передаче в функцию копируется целиком,
     поэтому `modifyArray` не меняет оригинал;
   - слайс — заголовок (указатель на массив, len, cap); копируется заголовок,
     но данные общие, поэтому `modifySlice` меняет оригинал;
   - `slice := array[:]` смотрит на тот же массив, поэтому изменение через слайс
     видно и в `array`.

## Ожидаемый вывод

```
Before modifyArray: [1 2 3]
Inside modifyArray: [10 2 3]
After modifyArray: [1 2 3]
Before modifySlice: [1 2 3]
Inside modifySlice: [10 2 3]
After modifySlice: [10 2 3]
Final array: [10 2 3]
```
