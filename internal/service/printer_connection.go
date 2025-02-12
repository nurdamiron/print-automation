// internal/service/printer_connection.go
package service

import (
    "fmt"
    "io"
    "net"
    "sync"
    "time"
    "strings"
)

type PrinterInfo struct {
    IP          string `json:"ip"`
    Name        string `json:"name"`
    Model       string `json:"model"`
    Status      string `json:"status"`
    IsOnline    bool   `json:"is_online"`
    Port        int    `json:"port"`
    Protocol    string `json:"protocol"`
}

type PrinterConnection struct {
    conn      net.Conn
    info      PrinterInfo
    isActive  bool
    mu        sync.Mutex
}

// PJL команды для работы с принтером
const (
    PJL_ENQUIRY    = "\x1B%-12345X@PJL INQUIRE\r\n"
    PJL_INFO       = "\x1B%-12345X@PJL INFO STATUS\r\n"
    PJL_ECHO       = "\x1B%-12345X@PJL ECHO\r\n"
)

func NewPrinterConnection(info PrinterInfo) *PrinterConnection {
    return &PrinterConnection{
        info: info,
    }
}

func (pc *PrinterConnection) Connect() error {
    pc.mu.Lock()
    defer pc.mu.Unlock()

    if pc.isActive {
        return fmt.Errorf("already connected")
    }

    address := fmt.Sprintf("%s:%d", pc.info.IP, pc.info.Port)

    conn, err := net.DialTimeout("tcp", address, 5*time.Second)
    if err != nil {
        return fmt.Errorf("failed to connect to printer at %s: %v", address, err)
    }

    pc.conn = conn
    pc.isActive = true

    if err := pc.sendEcho(); err != nil {
        pc.Disconnect()
        return fmt.Errorf("printer connection test failed: %v", err)
    }

    return nil
}

func (pc *PrinterConnection) Disconnect() error {
    pc.mu.Lock()
    defer pc.mu.Unlock()

    if !pc.isActive {
        return nil
    }

    if pc.conn != nil {
        err := pc.conn.Close()
        pc.conn = nil
        pc.isActive = false
        return err
    }

    return nil
}

func (pc *PrinterConnection) IsConnected() bool {
    pc.mu.Lock()
    defer pc.mu.Unlock()
    return pc.isActive
}

func (pc *PrinterConnection) sendCommand(cmd string) (string, error) {
    if !pc.isActive || pc.conn == nil {
        return "", fmt.Errorf("printer not connected")
    }

    pc.conn.SetDeadline(time.Now().Add(5 * time.Second))

    _, err := pc.conn.Write([]byte(cmd))
    if err != nil {
        return "", fmt.Errorf("failed to send command: %v", err)
    }

    buffer := make([]byte, 1024)
    n, err := pc.conn.Read(buffer)
    if err != nil {
        return "", fmt.Errorf("failed to read response: %v", err)
    }

    return string(buffer[:n]), nil
}

func (pc *PrinterConnection) sendEcho() error {
    response, err := pc.sendCommand(PJL_ECHO)
    if err != nil {
        return err
    }

    if !strings.Contains(response, "ECHO") {
        return fmt.Errorf("invalid echo response from printer")
    }

    return nil
}

func (pc *PrinterConnection) GetStatus() (string, error) {
    response, err := pc.sendCommand(PJL_INFO)
    if err != nil {
        return "", err
    }

    return parseStatusResponse(response), nil
}

// Добавляем новые методы для печати
func (pc *PrinterConnection) Print(document io.Reader) error {
    pc.mu.Lock()
    defer pc.mu.Unlock()

    if !pc.isActive || pc.conn == nil {
        return fmt.Errorf("printer not connected")
    }

    // Отправляем команду начала печати
    startCmd := "\x1B%-12345X@PJL ENTER LANGUAGE=PCL\r\n"
    if _, err := pc.conn.Write([]byte(startCmd)); err != nil {
        return fmt.Errorf("failed to send start command: %v", err)
    }

    // Копируем данные документа на принтер
    // В реальном приложении здесь нужно добавить форматирование и контроль
    _, err := io.Copy(pc.conn, document)
    if err != nil {
        return fmt.Errorf("failed to send document: %v", err)
    }

    // Отправляем команду завершения печати
    endCmd := "\x1B%-12345X@PJL EOJ\r\n"
    if _, err := pc.conn.Write([]byte(endCmd)); err != nil {
        return fmt.Errorf("failed to send end command: %v", err)
    }

    return nil
}


// Добавляем метод для отмены печати
func (pc *PrinterConnection) CancelPrint() error {
    pc.mu.Lock()
    defer pc.mu.Unlock()

    if !pc.isActive || pc.conn == nil {
        return fmt.Errorf("printer not connected")
    }

    cancelCmd := "\x1B%-12345X@PJL CANCEL\r\n"
    if _, err := pc.conn.Write([]byte(cancelCmd)); err != nil {
        return fmt.Errorf("failed to send cancel command: %v", err)
    }

    return nil
}


func parseStatusResponse(response string) string {
    response = strings.ToLower(response)
    
    if strings.Contains(response, "ready") {
        return "READY"
    } else if strings.Contains(response, "busy") {
        return "BUSY"
    } else if strings.Contains(response, "paper jam") {
        return "PAPER_JAM"
    } else if strings.Contains(response, "out of paper") {
        return "OUT_OF_PAPER"
    } else if strings.Contains(response, "toner low") {
        return "TONER_LOW"
    }
    
    return "UNKNOWN"
}

func DiscoverPrinters(subnetMask string) ([]PrinterInfo, error) {
    printers := make([]PrinterInfo, 0)
    
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return nil, fmt.Errorf("failed to get network interfaces: %v", err)
    }

    for _, addr := range addrs {
        if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                ports := []int{9100, 515, 631}
                
                for _, port := range ports {
                    baseIP := ipnet.IP.To4()
                    for i := 1; i < 255; i++ {
                        testIP := fmt.Sprintf("%d.%d.%d.%d", 
                            baseIP[0], baseIP[1], baseIP[2], i)
                        
                        if printer := testPrinterConnection(testIP, port); printer != nil {
                            printers = append(printers, *printer)
                        }
                    }
                }
            }
        }
    }

    return printers, nil
}

func testPrinterConnection(ip string, port int) *PrinterInfo {
    address := fmt.Sprintf("%s:%d", ip, port)
    conn, err := net.DialTimeout("tcp", address, 1*time.Second)
    if err != nil {
        return nil
    }
    defer conn.Close()

    _, err = conn.Write([]byte(PJL_ENQUIRY))
    if err != nil {
        return nil
    }

    buffer := make([]byte, 1024)
    conn.SetReadDeadline(time.Now().Add(2 * time.Second))
    n, err := conn.Read(buffer)
    if err != nil {
        return nil
    }

    response := string(buffer[:n])
    if strings.Contains(response, "PJL") || strings.Contains(response, "READY") {
        return &PrinterInfo{
            IP:       ip,
            Port:     port,
            Status:   "READY",
            IsOnline: true,
            Protocol: "RAW",
            Name:     fmt.Sprintf("Printer at %s", ip),
            Model:    "Auto-detected printer",
        }
    }

    return nil
}