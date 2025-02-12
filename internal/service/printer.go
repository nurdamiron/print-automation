// internal/service/printer.go
package service

import (
    "fmt"
    "io"
    "sync"
    "time"
)

type PrintJob struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
}

type PrintQueue struct {
    Jobs []PrintJob `json:"jobs"`
}

type PrinterService struct {
    connections map[string]*PrinterConnection
    printQueue  map[string][]PrintJob // key - printer ID
    mu          sync.RWMutex
}

func NewPrinterService() *PrinterService {
    return &PrinterService{
        connections: make(map[string]*PrinterConnection),
        printQueue:  make(map[string][]PrintJob),
    }
}

// PrintDocument отправляет документ на печать
func (s *PrinterService) PrintDocument(printerID string, document io.Reader) error {
    s.mu.RLock()
    connection, exists := s.connections[printerID]
    s.mu.RUnlock()

    if !exists {
        return fmt.Errorf("printer not found")
    }

    if !connection.IsConnected() {
        return fmt.Errorf("printer not connected")
    }

    // Создаем новое задание печати
    job := PrintJob{
        ID:        fmt.Sprintf("job_%d", time.Now().Unix()),
        Name:      "Document",
        Status:    "pending",
        CreatedAt: time.Now(),
    }

    // Добавляем задание в очередь
    s.mu.Lock()
    if _, exists := s.printQueue[printerID]; !exists {
        s.printQueue[printerID] = make([]PrintJob, 0)
    }
    s.printQueue[printerID] = append(s.printQueue[printerID], job)
    s.mu.Unlock()

    // В реальном приложении здесь бы происходила отправка документа на принтер
    // Сейчас просто имитируем успешную отправку
    go func() {
        time.Sleep(2 * time.Second) // Имитация процесса печати
        s.mu.Lock()
        for i, j := range s.printQueue[printerID] {
            if j.ID == job.ID {
                s.printQueue[printerID][i].Status = "completed"
                break
            }
        }
        s.mu.Unlock()
    }()

    return nil
}


func (s *PrinterService) ConnectPrinter(id, name, ip string, port int) (*PrinterInfo, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    // Проверяем, не подключен ли уже принтер
    if conn, exists := s.connections[id]; exists && conn.IsConnected() {
        return &conn.info, nil
    }

    // Создаем новое подключение
    printerInfo := PrinterInfo{
        IP:       ip,
        Name:     name,
        Port:     port,
        Protocol: "RAW",
    }

    connection := NewPrinterConnection(printerInfo)
    err := connection.Connect()
    if err != nil {
        return nil, fmt.Errorf("failed to connect to printer: %v", err)
    }

    // Получаем начальный статус
    status, err := connection.GetStatus()
    if err != nil {
        connection.Disconnect()
        return nil, fmt.Errorf("failed to get printer status: %v", err)
    }

    // Обновляем информацию о принтере
    connection.info.Status = status
    connection.info.IsOnline = true

    // Сохраняем подключение
    s.connections[id] = connection

    return &connection.info, nil
}

func (s *PrinterService) GetPrinterStatus(id string) (*PrinterInfo, error) {
    s.mu.RLock()
    connection, exists := s.connections[id]
    s.mu.RUnlock()

    if !exists {
        return nil, fmt.Errorf("printer not found")
    }

    if !connection.IsConnected() {
        return nil, fmt.Errorf("printer not connected")
    }

    status, err := connection.GetStatus()
    if err != nil {
        return nil, fmt.Errorf("failed to get printer status: %v", err)
    }

    connection.info.Status = status
    return &connection.info, nil
}

func (s *PrinterService) DisconnectPrinter(id string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    connection, exists := s.connections[id]
    if !exists {
        return fmt.Errorf("printer not found")
    }

    err := connection.Disconnect()
    if err != nil {
        return fmt.Errorf("failed to disconnect printer: %v", err)
    }

    delete(s.connections, id)
    return nil
}

func (s *PrinterService) DiscoverPrinters() ([]PrinterInfo, error) {
    return DiscoverPrinters("255.255.255.0") // Используем стандартную маску подсети
}

// StartStatusMonitoring запускает мониторинг статуса принтера
func (s *PrinterService) StartStatusMonitoring(id string, interval time.Duration) error {
    s.mu.RLock()
    connection, exists := s.connections[id]
    s.mu.RUnlock()

    if !exists {
        return fmt.Errorf("printer not found")
    }

    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()

        for range ticker.C {
            if !connection.IsConnected() {
                return
            }

            status, err := connection.GetStatus()
            if err != nil {
                fmt.Printf("Failed to get printer status: %v\n", err)
                continue
            }

            connection.info.Status = status
        }
    }()

    return nil
}



// GetPrintQueue возвращает текущую очередь печати для принтера
func (s *PrinterService) GetPrintQueue(printerID string) (*PrintQueue, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    connection, exists := s.connections[printerID]
    if !exists {
        return nil, fmt.Errorf("printer not found")
    }

    if !connection.IsConnected() {
        return nil, fmt.Errorf("printer not connected")
    }

    jobs, exists := s.printQueue[printerID]
    if !exists {
        return &PrintQueue{Jobs: make([]PrintJob, 0)}, nil
    }

    return &PrintQueue{Jobs: jobs}, nil
}