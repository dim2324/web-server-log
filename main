package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "os"
    "time"
)

func main() {
    // Парсинг аргументов командной строки
    filename := flag.String("file", "testdata/logs.csv", "путь к файлу с логами")
    workers := flag.Int("workers", 3, "количество воркеров")
    timeout := flag.Duration("timeout", 10*time.Second, "таймаут обработки")
    flag.Parse()

    // Проверяем существование файла
    if _, err := os.Stat(*filename); os.IsNotExist(err) {
        log.Fatalf("Файл %s не найден", *filename)
    }

    // Создаем контекст с таймаутом
    ctx, cancel := context.WithTimeout(context.Background(), *timeout)
    defer cancel()

    fmt.Println("🚀 Начинаем обработку логов...")
    fmt.Printf("📁 Файл: %s\n", *filename)
    fmt.Printf("👥 Количество воркеров: %d\n", *workers)
    fmt.Printf("⏱️  Таймаут: %v\n\n", *timeout)

    startTime := time.Now()

    // Читаем логи из файла
    logChan, err := readLogs(*filename)
    if err != nil {
        log.Fatalf("Ошибка чтения логов: %v", err)
    }

    // Обрабатываем логи с использованием worker pool
    processedChan := processLogs(ctx, logChan, *workers)

    // Фильтруем только ошибки (статус код >= 400)
    errorChan := filterLogs(processedChan, 400)

    // Подсчитываем статистику по ошибкам
    errorStats := calculateStats(errorChan)

    // Выводим результаты
    fmt.Println("=== 📊 РЕЗУЛЬТАТЫ ОБРАБОТКИ ===")
    fmt.Printf("⏱️  Время выполнения: %v\n", time.Since(startTime))
    fmt.Printf("📊 Всего обработано запросов: %d\n", errorStats.TotalRequests)
    fmt.Printf("❌ Количество ошибок (4xx, 5xx): %d\n", errorStats.ErrorCount)
    
    if errorStats.TotalRequests > 0 {
        fmt.Printf("📈 Среднее время ответа для ошибочных запросов: %.2f мс\n", errorStats.AverageRespTime)
    }

    if len(errorStats.RequestsByIP) > 0 {
        printTopIPs(errorStats.RequestsByIP, 3)
    } else {
        fmt.Println("\n📭 Нет данных для отображения топ IP")
    }

    // Дополнительно: обрабатываем все запросы для полной статистики
    fmt.Println("\n=== 📊 ПОЛНАЯ СТАТИСТИКА ===")
    
    // Перечитываем файл для полной статистики
    fullLogChan, err := readLogs(*filename)
    if err != nil {
        log.Fatalf("Ошибка чтения логов: %v", err)
    }
    
    fullProcessedChan := processLogs(ctx, fullLogChan, *workers)
    fullStats := calculateStats(fullProcessedChan)
    
    fmt.Printf("📊 Всего запросов: %d\n", fullStats.TotalRequests)
    fmt.Printf("❌ Всего ошибок: %d\n", fullStats.ErrorCount)
    
    if fullStats.TotalRequests > 0 {
        fmt.Printf("📈 Процент ошибок: %.2f%%\n", float64(fullStats.ErrorCount)/float64(fullStats.TotalRequests)*100)
        fmt.Printf("⏱️  Среднее время ответа: %.2f мс\n", fullStats.AverageRespTime)
    }
    
    if len(fullStats.RequestsByIP) > 0 {
        printTopIPs(fullStats.RequestsByIP, 5)
    }

    fmt.Println("\n✅ Обработка завершена!")
}
