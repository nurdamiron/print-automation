package services

import (
    "fmt"
    "io"
    "net"
    "os"
    "time"
)

// SendToPrinterRaw отправляет локальный файл (PDF/PS/PCL) на принтер через RAW-порт (9100).
func SendToPrinterRaw(filePath, printerIP string, printerPort int) error {
    // Открываем файл
    f, err := os.Open(filePath)
    if err != nil {
        return fmt.Errorf("не удалось открыть файл для печати: %w", err)
    }
    defer f.Close()

    // Формируем адрес
    addr := fmt.Sprintf("%s:%d", printerIP, printerPort)

    // Пытаемся подключиться с таймаутом
    conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
    if err != nil {
        return fmt.Errorf("не удалось подключиться к принтеру [%s]: %w", addr, err)
    }
    defer conn.Close()

    // Копируем данные файла в соединение
    _, err = io.Copy(conn, f)
    if err != nil {
        return fmt.Errorf("ошибка при передаче данных на принтер: %w", err)
    }

    // Возможен дополнительный завершающий символ, но зачастую HP принтеры завершают задание автоматически
    return nil
}

func CheckPrinterConnection(ip string, port int) error {
    addr := fmt.Sprintf("%s:%d", ip, port)
    conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
    if err != nil {
        return fmt.Errorf("принтер недоступен: %w", err)
    }
    conn.Close()
    return nil
}