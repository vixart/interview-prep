# Задачи с собеседований по Go

Сборник задач с собеседований по Go — с условиями и решениями.

В каждой папке:
- `task.md` — условие задачи;
- `solution.go` — решение кодом (запускается `go run solution.go`, внутри тесты в `main`);
- `solution.md` — разбор/ответ (для задач «что выведет код», ревью и дизайна);
- `fixed.go` / `schema.sql` — исправленный код для ревью-задач и SQL для задач по БД.

## 01_arrays_and_slices — Массивы и слайсы

- **1.1. Слияние двух отсортированных массивов (merge in-place)** — [задача](01_arrays_and_slices/01_merge_sorted_arrays/task.md), [код](01_arrays_and_slices/01_merge_sorted_arrays/solution.go)
- **1.2. Удаление всех нулей из слайса** — [задача](01_arrays_and_slices/02_remove_zeros/task.md), [код](01_arrays_and_slices/02_remove_zeros/solution.go)
- **1.3. Произведение всех элементов кроме i-го** — [задача](01_arrays_and_slices/03_product_except_self/task.md), [код](01_arrays_and_slices/03_product_except_self/solution.go)
- **1.4. Find All Numbers Disappeared in an Array** — [задача](01_arrays_and_slices/04_find_disappeared_numbers/task.md), [код](01_arrays_and_slices/04_find_disappeared_numbers/solution.go)
- **1.5. Вращение массива** — [задача](01_arrays_and_slices/05_rotate_array/task.md), [код](01_arrays_and_slices/05_rotate_array/solution.go)

## 02_linked_lists — Связные списки

- **2.1. Удаление N-го элемента с конца односвязного списка** — [задача](02_linked_lists/06_remove_nth_from_end/task.md), [код](02_linked_lists/06_remove_nth_from_end/solution.go)
- **2.2. Разворот связного списка** — [задача](02_linked_lists/07_reverse_linked_list/task.md), [код](02_linked_lists/07_reverse_linked_list/solution.go)
- **2.3. Поиск цикла в связном списке (Floyd's Cycle Detection)** — [задача](02_linked_lists/08_linked_list_cycle/task.md), [код](02_linked_lists/08_linked_list_cycle/solution.go)
- **2.4. Слияние двух отсортированных связных списков** — [задача](02_linked_lists/09_merge_two_sorted_lists/task.md), [код](02_linked_lists/09_merge_two_sorted_lists/solution.go)
- **2.5. Удаление дубликатов из отсортированного связного списка** — [задача](02_linked_lists/10_remove_duplicates/task.md), [код](02_linked_lists/10_remove_duplicates/solution.go)

## 03_strings — Строки

- **3.1. Анаграммы** — [задача](03_strings/11_anagram/task.md), [код](03_strings/11_anagram/solution.go)
- **3.2. Реверс слов (с сохранением палиндромов)** — [задача](03_strings/12_reverse_words/task.md), [код](03_strings/12_reverse_words/solution.go)
- **3.3. Валидация скобочной последовательности** — [задача](03_strings/13_valid_parentheses/task.md), [код](03_strings/13_valid_parentheses/solution.go)
- **3.4. Самая длинная подстрока без повторяющихся символов** — [задача](03_strings/14_longest_substring/task.md), [код](03_strings/14_longest_substring/solution.go)
- **3.5. Поиск всех анаграмм в строке** — [задача](03_strings/15_find_all_anagrams/task.md), [код](03_strings/15_find_all_anagrams/solution.go)

## 04_hash_tables — Хеш-таблицы

- **4.1. Two Sum — поиск суммы в массиве** — [задача](04_hash_tables/16_two_sum/task.md), [код](04_hash_tables/16_two_sum/solution.go)
- **4.2. Группировка серверов по стабильности** — [задача](04_hash_tables/17_group_by_stability/task.md), [код](04_hash_tables/17_group_by_stability/solution.go)
- **4.3. Фильтрация продавцов по городам** — [задача](04_hash_tables/18_filter_sellers/task.md), [код](04_hash_tables/18_filter_sellers/solution.go)
- **4.4. Определение чемпионов соревнований** — [задача](04_hash_tables/19_find_champions/task.md), [код](04_hash_tables/19_find_champions/solution.go)
- **4.5. Первый уникальный символ в строке** — [задача](04_hash_tables/20_first_unique_char/task.md), [код](04_hash_tables/20_first_unique_char/solution.go)

## 05_two_pointers — Два указателя

- **5.1. K ближайших элементов в отсортированном массиве** — [задача](05_two_pointers/21_k_closest/task.md), [код](05_two_pointers/21_k_closest/solution.go), [разбор](05_two_pointers/21_k_closest/solution.md)
- **5.2. Контейнер с максимальным количеством воды** — [задача](05_two_pointers/22_container_with_most_water/task.md), [код](05_two_pointers/22_container_with_most_water/solution.go)
- **5.3. Минимальная разница между k элементами** — [задача](05_two_pointers/23_min_difference_k_elements/task.md), [код](05_two_pointers/23_min_difference_k_elements/solution.go), [разбор](05_two_pointers/23_min_difference_k_elements/solution.md)
- **5.4. Наименьшая разница по модулю между двумя массивами** — [задача](05_two_pointers/24_min_absolute_difference/task.md), [код](05_two_pointers/24_min_absolute_difference/solution.go), [разбор](05_two_pointers/24_min_absolute_difference/solution.md)
- **5.5. Удаление дубликатов из отсортированного массива in-place** — [задача](05_two_pointers/25_remove_duplicates_in_place/task.md), [код](05_two_pointers/25_remove_duplicates_in_place/solution.go)

## 06_sliding_window — Скользящее окно

- **6.1. Максимум в скользящем окне** — [задача](06_sliding_window/26_max_sliding_window/task.md), [код](06_sliding_window/26_max_sliding_window/solution.go), [разбор](06_sliding_window/26_max_sliding_window/solution.md)
- **6.2. Минимальное окно с подстрокой** — [задача](06_sliding_window/27_minimum_window_substring/task.md), [код](06_sliding_window/27_minimum_window_substring/solution.go), [разбор](06_sliding_window/27_minimum_window_substring/solution.md)
- **6.3. Оптимальное планирование отпуска** — [задача](06_sliding_window/28_optimal_vacation/task.md), [код](06_sliding_window/28_optimal_vacation/solution.go), [разбор](06_sliding_window/28_optimal_vacation/solution.md)
- **6.4. Максимальная сумма подмассива длины k** — [задача](06_sliding_window/29_max_sum_subarray/task.md), [код](06_sliding_window/29_max_sum_subarray/solution.go)
- **6.5. Подмассив с заданной суммой (минимальная длина)** — [задача](06_sliding_window/30_subarray_min_len_sum/task.md), [код](06_sliding_window/30_subarray_min_len_sum/solution.go), [разбор](06_sliding_window/30_subarray_min_len_sum/solution.md)

## 07_go_theory — Теория Go: указатели, defer, слайсы, мапы, интерфейсы

- **4.1.1. Указатели: как работают и где ловят новичков** — [задача](07_go_theory/01_pointers/task.md), [разбор](07_go_theory/01_pointers/solution.md)
- **4.1.2. Defer: порядок выполнения и подводные камни** — [задача](07_go_theory/02_defer_order/task.md), [разбор](07_go_theory/02_defer_order/solution.md)
- **4.1.3. Defer и указатели: когда значения меняются неожиданно** — [задача](07_go_theory/03_defer_and_pointers/task.md), [разбор](07_go_theory/03_defer_and_pointers/solution.md)
- **4.1.4. Указатели на элементы слайса и append** — [задача](07_go_theory/04_pointer_to_slice_elem/task.md), [разбор](07_go_theory/04_pointer_to_slice_elem/solution.md)
- **4.1.5. Массивы vs слайсы: в чём разница** — [задача](07_go_theory/05_arrays_vs_slices/task.md), [разбор](07_go_theory/05_arrays_vs_slices/solution.md)
- **4.1.6. Ещё один defer: значение получателя фиксируется в момент объявления** — [задача](07_go_theory/06_defer_value_receiver/task.md), [разбор](07_go_theory/06_defer_value_receiver/solution.md)
- **4.1.7. Выравнивание типов: как Go раскладывает поля структуры в памяти** — [задача](07_go_theory/07_type_alignment/task.md), [разбор](07_go_theory/07_type_alignment/solution.md)
- **4.2.1. Конкатенация строк: почему это медленно и как ускорить** — [задача](07_go_theory/08_string_concat/task.md), [разбор](07_go_theory/08_string_concat/solution.md)
- **4.2.2. Длина строки: руны, байты и UTF-8** — [задача](07_go_theory/09_string_length_runes/task.md), [разбор](07_go_theory/09_string_length_runes/solution.md)
- **4.2.3. Итерация по строке: правильный способ обхода** — [задача](07_go_theory/10_string_iteration/task.md), [разбор](07_go_theory/10_string_iteration/solution.md)
- **4.2.4. Работа с рунами: вывод символов корректно** — [задача](07_go_theory/11_string_mutation_runes/task.md), [разбор](07_go_theory/11_string_mutation_runes/solution.md)
- **4.2.5. Конвертация `[]byte` ↔ `string` без аллокаций** — [задача](07_go_theory/12_bytes_string_conversion/task.md), [разбор](07_go_theory/12_bytes_string_conversion/solution.md)
- **4.3.1. Append: как меняется capacity и когда создаётся новый массив** — [задача](07_go_theory/13_append_capacity/task.md), [разбор](07_go_theory/13_append_capacity/solution.md)
- **4.3.2. Append и изменение слайса — вариант A** — [задача](07_go_theory/14_append_variant_a/task.md), [разбор](07_go_theory/14_append_variant_a/solution.md)
- **4.3.3. Append и изменение слайса — вариант B** — [задача](07_go_theory/15_append_variant_b/task.md), [разбор](07_go_theory/15_append_variant_b/solution.md)
- **4.3.4. Генератор слайса из N уникальных случайных чисел** — [задача](07_go_theory/16_uniq_random_generator/task.md), [код](07_go_theory/16_uniq_random_generator/solution.go), [разбор](07_go_theory/16_uniq_random_generator/solution.md)
- **4.3.5. Capacity слайса: предсказание len и cap** — [задача](07_go_theory/17_slice_capacity/task.md), [разбор](07_go_theory/17_slice_capacity/solution.md)
- **4.3.6. Ещё один append: запись в чужую ёмкость** — [задача](07_go_theory/18_append_shared_backing/task.md), [разбор](07_go_theory/18_append_shared_backing/solution.md)
- **4.3.7. Модификация слайса через указатель на элемент** — [задача](07_go_theory/19_modify_slice_pointer/task.md), [разбор](07_go_theory/19_modify_slice_pointer/solution.md)
- **4.3.8. Обновление слайса через функцию: что изменится** — [задача](07_go_theory/20_append_in_func/task.md), [разбор](07_go_theory/20_append_in_func/solution.md)
- **4.3.9. «Магия» capacity при создании срезов** — [задача](07_go_theory/21_slice_magic/task.md), [разбор](07_go_theory/21_slice_magic/solution.md)
- **4.4.1. Конкурентное обновление мапы: GetOrCreate** — [задача](07_go_theory/22_map_get_or_create/task.md), [код](07_go_theory/22_map_get_or_create/solution.go), [разбор](07_go_theory/22_map_get_or_create/solution.md)
- **4.4.2. Исправление бага: безопасная работа с мапой (общая мапа в канале)** — [задача](07_go_theory/23_map_in_channel/task.md), [разбор](07_go_theory/23_map_in_channel/solution.md)
- **4.4.3. Ограничение размера мапы (реализация LRU)** — [задача](07_go_theory/24_map_limit_lru/task.md), [код](07_go_theory/24_map_limit_lru/solution.go), [разбор](07_go_theory/24_map_limit_lru/solution.md)
- **4.4.4. MergeToMap — объединение значений мапы без дубликатов** — [задача](07_go_theory/25_merge_to_map/task.md), [код](07_go_theory/25_merge_to_map/solution.go), [разбор](07_go_theory/25_merge_to_map/solution.md)
- **4.4.5. sync.Map: когда использовать вместо обычной мапы (GetOrCompute)** — [задача](07_go_theory/26_get_or_compute/task.md), [код](07_go_theory/26_get_or_compute/solution.go), [разбор](07_go_theory/26_get_or_compute/solution.md)
- **4.4.6. Случайный порядок обхода мапы** — [задача](07_go_theory/27_map_random_order/task.md), [разбор](07_go_theory/27_map_random_order/solution.md)
- **4.5.1. Приведение типов: type assertion и type switch** — [задача](07_go_theory/28_type_assertion/task.md), [разбор](07_go_theory/28_type_assertion/solution.md)
- **4.5.2. Вернуть ошибку из функции, не подключая дополнительных пакетов** — [задача](07_go_theory/29_error_interface/task.md), [код](07_go_theory/29_error_interface/solution.go), [разбор](07_go_theory/29_error_interface/solution.md)
- **4.5.3. Исправление бага: реализация интерфейса и получатели методов** — [задача](07_go_theory/30_method_sets/task.md), [разбор](07_go_theory/30_method_sets/solution.md)
- **4.5.4. Ещё один баг с интерфейсами: type assertion к `interface{}`** — [задача](07_go_theory/31_nil_interface_assertion/task.md), [разбор](07_go_theory/31_nil_interface_assertion/solution.md)
- **4.5.5. Nil-интерфейсы: почему сравнение с nil не работает** — [задача](07_go_theory/32_nil_interface_compare/task.md), [разбор](07_go_theory/32_nil_interface_compare/solution.md)

## 08_concurrency — Конкурентность

- **4.6.1. Таймаут для запроса и конкурентная среда** — [задача](08_concurrency/01_request_timeout/task.md), [код](08_concurrency/01_request_timeout/solution.go), [разбор](08_concurrency/01_request_timeout/solution.md)
- **4.6.2. Кастомный таймаут: контроль медленной зависимости** — [задача](08_concurrency/02_custom_timeout/task.md), [код](08_concurrency/02_custom_timeout/solution.go), [разбор](08_concurrency/02_custom_timeout/solution.md)
- **4.6.3. Исправление бага с каналами: deadlock и утечки** — [задача](08_concurrency/03_channel_deadlock_fix/task.md), [код](08_concurrency/03_channel_deadlock_fix/solution.go), [разбор](08_concurrency/03_channel_deadlock_fix/solution.md)
- **4.6.4. Топ-K документов для пользователя** — [задача](08_concurrency/04_top_documents/task.md), [код](08_concurrency/04_top_documents/solution.go), [разбор](08_concurrency/04_top_documents/solution.md)
- **4.6.5. Исправление бага в парковке: ограничение числа мест** — [задача](08_concurrency/05_parking_semaphore/task.md), [код](08_concurrency/05_parking_semaphore/solution.go), [разбор](08_concurrency/05_parking_semaphore/solution.md)
- **4.6.6. Аналитика в реальном времени** — [задача](08_concurrency/06_realtime_analytics/task.md), [код](08_concurrency/06_realtime_analytics/solution.go), [разбор](08_concurrency/06_realtime_analytics/solution.md)
- **4.6.7. Параллельные HTTP-запросы к списку URL** — [задача](08_concurrency/07_parallel_url_requests/task.md), [код](08_concurrency/07_parallel_url_requests/solution.go), [разбор](08_concurrency/07_parallel_url_requests/solution.md)
- **4.6.8. Web crawler: статистика загрузок за последние 10 минут** — [задача](08_concurrency/08_web_crawler_stats/task.md), [код](08_concurrency/08_web_crawler_stats/solution.go), [разбор](08_concurrency/08_web_crawler_stats/solution.md)
- **4.6.9. Rate Limiter — ограничение RPS** — [задача](08_concurrency/09_rate_limiter/task.md), [код](08_concurrency/09_rate_limiter/solution.go), [разбор](08_concurrency/09_rate_limiter/solution.md)
- **4.6.10. Client-side балансировщик нагрузки** — [задача](08_concurrency/10_load_balancer/task.md), [код](08_concurrency/10_load_balancer/solution.go), [разбор](08_concurrency/10_load_balancer/solution.md)
- **4.6.11. Параллельный поиск на нескольких серверах** — [задача](08_concurrency/11_parallel_search/task.md), [код](08_concurrency/11_parallel_search/solution.go), [разбор](08_concurrency/11_parallel_search/solution.md)
- **4.6.12. MultiProcess — параллельный запуск асинхронных задач** — [задача](08_concurrency/12_async_jobs/task.md), [код](08_concurrency/12_async_jobs/solution.go), [разбор](08_concurrency/12_async_jobs/solution.md)
- **4.6.13. Версионирование документов (конкурентная обработка обновлений)** — [задача](08_concurrency/13_document_versioning/task.md), [код](08_concurrency/13_document_versioning/solution.go), [разбор](08_concurrency/13_document_versioning/solution.md)

## 09_concurrency_patterns — Паттерны конкурентности

- **4.7.1. Fan-in: объединение нескольких каналов в один** — [задача](09_concurrency_patterns/01_fan_in/task.md), [код](09_concurrency_patterns/01_fan_in/solution.go), [разбор](09_concurrency_patterns/01_fan_in/solution.md)
- **4.7.2. Sharded Cache — сегментированный in-memory кеш** — [задача](09_concurrency_patterns/02_sharded_cache/task.md), [код](09_concurrency_patterns/02_sharded_cache/solution.go), [разбор](09_concurrency_patterns/02_sharded_cache/solution.md)
- **4.7.3. Worker Pool — пул воркеров для обработки изображений** — [задача](09_concurrency_patterns/03_worker_pool/task.md), [код](09_concurrency_patterns/03_worker_pool/solution.go), [разбор](09_concurrency_patterns/03_worker_pool/solution.md)
- **4.7.4. Pipeline — конвейер обработки финансовых транзакций** — [задача](09_concurrency_patterns/04_pipeline/task.md), [код](09_concurrency_patterns/04_pipeline/solution.go), [разбор](09_concurrency_patterns/04_pipeline/solution.md)
- **4.7.5. Semaphore — ограничение числа одновременных соединений** — [задача](09_concurrency_patterns/05_semaphore/task.md), [код](09_concurrency_patterns/05_semaphore/solution.go), [разбор](09_concurrency_patterns/05_semaphore/solution.md)

## 10_code_review — Code review

- **4.8.1. Ревью кода: сбор курсов валют из банков** — [задача](10_code_review/01_bank_rates/task.md), [разбор](10_code_review/01_bank_rates/solution.md)
- **4.8.2. Ревью кода: аналитика посещений** — [задача](10_code_review/02_visits_analytics/task.md), [разбор](10_code_review/02_visits_analytics/solution.md)
- **4.8.3. Ревью кода: deadlock в каналах** — [задача](10_code_review/03_channel_deadlock/task.md), [исправленный код](10_code_review/03_channel_deadlock/fixed.go), [разбор](10_code_review/03_channel_deadlock/solution.md)
- **4.8.4. Ревью кода: отмена HTTP-запросов и утечки ресурсов** — [задача](10_code_review/04_http_cancel/task.md), [разбор](10_code_review/04_http_cancel/solution.md)
- **4.8.5. Ревью кода: worker pool с ошибками синхронизации** — [задача](10_code_review/05_worker_pool/task.md), [исправленный код](10_code_review/05_worker_pool/fixed.go), [разбор](10_code_review/05_worker_pool/solution.md)
- **4.8.6. Ревью кода: race condition в кеше** — [задача](10_code_review/06_cache_race/task.md), [исправленный код](10_code_review/06_cache_race/fixed.go), [разбор](10_code_review/06_cache_race/solution.md)
- **4.8.7. Ревью кода: JSON API для управления пользователями** — [задача](10_code_review/07_json_api/task.md), [разбор](10_code_review/07_json_api/solution.md)
- **4.8.8. Ревью кода: рефакторинг работы с каналами** — [задача](10_code_review/08_channel_refactor/task.md), [исправленный код](10_code_review/08_channel_refactor/fixed.go), [разбор](10_code_review/08_channel_refactor/solution.md)
- **4.8.9. Ревью кода: кеш с проблемами конкурентности** — [задача](10_code_review/09_cache_concurrency/task.md), [исправленный код](10_code_review/09_cache_concurrency/fixed.go), [разбор](10_code_review/09_cache_concurrency/solution.md)

## 11_system_design — Проектирование

- **4.9.1. Проектирование базы данных для чата** — [задача](11_system_design/01_chat_db/task.md), [SQL](11_system_design/01_chat_db/schema.sql), [разбор](11_system_design/01_chat_db/solution.md)
- **4.9.2. Медленный сервис аналитики** — [задача](11_system_design/02_slow_service/task.md), [разбор](11_system_design/02_slow_service/solution.md)
