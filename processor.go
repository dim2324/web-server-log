package main

import (
    "bufio"
    "context"
    "encoding/csv"
    "fmt"
    "io"
    "log"
    "os"
    "strconv"
    "strings"
    "sync"
)

// LogEntry представляет запись лога
type LogEntry struct {
    Timestamp    string // время в формате "2024-01-15 10:30:00"
    IP           string // IP адрес клиента
    Method       string // HTTP метод (GET, POST и т.д.)
    URL          string // путь запроса
    StatusCode   int    // HTTP статус код
    ResponseTime int    // время ответа в миллисекундах
}

// Statistics представляет статистику по логам
type Statistics struct {
    TotalRequests   int            // общее количество запросов
    ErrorCount      int            // количество ошибок (статус >= 400)
    RequestsByIP    map[string]int // количество запросов с каждого IP
    AverageRespTime float64        // среднее время ответа
    TotalRespTime   int64          // суммарное время ответов (для вычисления среднего)
}

// Парсинг строки CSV в структуру LogEntry
func parseLogLine(line string) (LogEntry, error) {
    parts := strings.Split(line, ",")
    if len(parts) != 6 {
        return LogEntry{}, fmt.Errorf("неверное количество полей: ожидается 6, получено %d", len(parts))
    }

    // Очищаем поля от пробелов
    for i := range parts {
        parts[i] = strings.TrimSpace(parts[i])
    }

    statusCode, err := strconv.Atoi(parts[4])
    if err != nil {
        return LogEntry{}, fmt.Errorf("ошибка парсинга статус кода: %w", err)
    }

    responseTime, err := strconv.Atoi(parts[5])
    if err != nil {
        return LogEntry{}, fmt.Errorf("ошибка парсинга времени ответа: %w", err)
    }

    return LogEntry{
        Timestamp:    parts[0],
        IP:           parts[1],
        Method:       parts[2],
        URL:          parts[3],
        StatusCode:   statusCode,
        ResponseTime: responseTime,
    }, nil
}

// Чтение логов из файла
func readLogs(filename string) (<-chan LogEntry, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("ошибка открытия файла: %w", err)
    }

    out := make(chan LogEntry)

    go func() {
        defer file.Close()
        defer close(out)

        reader := csv.NewReader(bufio.NewReader(file))
        reader.TrimLeadingSpace = true
        reader.FieldsPerRecord = -1 // Разрешить переменное количество полей
        
        // Читаем все записи
        for {
            record, err := reader.Read()
            if err == io.EOF {
                break
            }
            if err != nil {
                log.Printf("Ошибка чтения CSV: %v", err)
                continue
            }

            // Пропускаем заголовок, если он есть
            if len(record) > 0 && strings.ToLower(record[0]) == "timestamp" {
                continue
            }

            // Объединяем поля обратно в строку для парсинга
            line := strings.Join(record, ",")
            entry, err := parseLogLine(line)
            if err != nil {
                log.Printf("Ошибка парсинга строки: %v, строка: %s", err, line)
                continue
            }

            select {
            case out <- entry:
            }
        }
    }()

    return out, nil
}

// Обработка логов с использованием worker pool
func processLogs(ctx context.Context, input <-chan LogEntry, numWorkers int) <-chan LogEntry {
    output := make(chan LogEntry)
    var wg sync.WaitGroup

    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for {
                select {
                case <-ctx.Done():
                    return
                case entry, ok := <-input:
                    if !ok {
                        return
                    }
                    // Здесь можно добавить дополнительную обработку
                    // Например, нормализацию URL, валидацию данных и т.д.
                    
                    select {
                    case <-ctx.Done():
                        return
                    case output <- entry:
                    }
                }
            }
        }(i)
    }

    go func() {
        wg.Wait()
        close(output)
    }()

    return output
}

// Фильтрация записей по статусу-коду
func filterLogs(input <-chan LogEntry, minStatus int) <-chan LogEntry {
    output := make(chan LogEntry)

    go func() {
        defer close(output)
        for entry := range input {
            if entry.StatusCode >= minStatus {
                select {
                case output <- entry:
                }
            }
        }
    }()

    return output
}

// Подсчет статистики
func calculateStats(input <-chan LogEntry) Statistics {
    stats := Statistics{
        RequestsByIP: make(map[string]int),
    }

    for entry := range input {
        stats.TotalRequests++
        
        if entry.StatusCode >= 400 {
            stats.ErrorCount++
        }
        
        stats.RequestsByIP[entry.IP]++
        stats.TotalRespTime += int64(entry.ResponseTime)
    }

    if stats.TotalRequests > 0 {
        stats.AverageRespTime = float64(stats.TotalRespTime) / float64(stats.TotalRequests)
    }

    return stats
}

// Вывод топ IP адресов
func printTopIPs(requestsByIP map[string]int, n int) {
    type ipCount struct {
        IP    string
        Count int
    }

    ips := make([]ipCount, 0, len(requestsByIP))
    for ip, count := range requestsByIP {
        ips = append(ips, ipCount{IP: ip, Count: count})
    }

    // Сортировка по убыванию количества запросов
    for i := 0; i < len(ips)-1; i++ {
        for j := i + 1; j < len(ips); j++ {
            if ips[i].Count < ips[j].Count {
                ips[i], ips[j] = ips[j], ips[i]
            }
        }
    }

    fmt.Printf("\nТоп-%d IP адресов:\n", n)
    fmt.Println("-------------------")
    for i := 0; i < n && i < len(ips); i++ {
        fmt.Printf("%d. %s - %d запросов\n", i+1, ips[i].IP, ips[i].Count)
    }
}
