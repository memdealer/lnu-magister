# Analysis Data Folder

Ця папка містить дані та скрипти для аналізу продуктивності HTTP/1.1, HTTP/2, REST та RPC з використанням мультиплексування. Дані та аналітика представлені у вигляді `.csv` файлів і Jupyter Notebook.

## Структура

### **CSV файли**
1. **`erlang_test_data.csv`**
    - Дані для тестування формули Ерланга B.

2. **`response_time_data.csv`**
    - Оброблені дані про час відповіді.

3. **`response_time_data_raw.csv`**
    - Сирі дані про час відповіді.

4. **`response_time_data_sorted.csv`**
    - Відсортовані дані про час відповіді.

5. **`response_time_large_data.csv`**
    - Великий набір даних для аналізу часу відповіді.

6. **`throughput_unfiltered_data.csv`**
    - Нефільтровані дані про пропускну спроможність.

---

### **Jupyter Notebook**
1. **`calc_erlang_b_calc.ipynb`**
    - Обчислення ймовірності втрати за формулою Ерланга B.

2. **`calc_http1_rest_http2_rpc_multiplexing_large.ipynb`**
    - Аналіз продуктивності HTTP/1.1, HTTP/2, REST і RPC на великому наборі даних.

3. **`calc_http1_rest_http2_rpc_multiplexing_sorted.ipynb`**
    - Аналіз продуктивності HTTP/1.1, HTTP/2, REST і RPC на відсортованих даних.

4. **`calc_resp_time.ipynb`**
    - Аналіз часу відповіді для різних протоколів.

5. **`calc_throughput_calc.ipynb`**
    - Обчислення пропускної спроможності.

6. **`calc_throughput_calc_raw.ipynb`**
    - Аналіз сирих даних пропускної спроможності.

---

### **Додаткові файли**
- **`README.md`**
    - Цей файл, що описує вміст папки.

- **`requirements.txt`**
    - Список залежностей Python для запуску Jupyter Notebook.

---

## Як використовувати
1. **Встановіть залежності**:
   ```bash
   pip install -r requirements.txt
