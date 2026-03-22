# web-server-log
Multithreaded web server log handler

- Reads entries from a CSV file with logs.
- Processes them in parallel using a worker pool (minimum 3 goroutines).
- Filters entries by status code (e.g., only 4xx and 5xx errors).
- Calculates statistics: total number of requests, number of errors, average response time, top IP addresses.
- Outputs results to the console.


Expected project structure:
log-processor/
├── main.go # program entry point
├── processor.go # log processing logic (readLogs, processLogs, filterLogs, calculateStats)
├── testdata/
│ └── logs.csv # test log file (minimum 10-15 entries)
├── go.mod # Go module file
├── README.md # launch instructions with command examples

# Log Processor - Многопоточный обработчик логов веб-сервера

## 📋 Описание
Многопоточный обработчик логов веб-сервера, который читает записи из CSV файла, обрабатывает их параллельно с использованием пула воркеров и выводит статистику.

## ✨ Возможности
- Чтение логов из CSV файла
- Параллельная обработка с использованием worker pool (минимум 3 горутины)
- Фильтрация записей по статус-коду (ошибки 4xx и 5xx)
- Подсчет статистики:
  - Общее количество запросов
  - Количество ошибок
  - Среднее время ответа
  - Топ IP адресов
- Контекст с таймаутом для отмены операций
- Обработка ошибок без остановки программы

## 🛠 Технологии
- Go 1.21+
- Горутины для параллельной обработки
- Каналы для передачи данных
- Паттерн Worker Pool
- Контекст для отмены операций
- Мьютексы для синхронизации

## 📁 Структура проекта
