package main

import (
    "fmt"
    "log"
    "net"
    "os"
    "runtime"
)

func main() {
    var username string

    // Ask user to enter a username
    fmt.Print("Enter username: ")
    fmt.Scanln(&username)

    // Get local IP address
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        log.Fatal(err)
    }
    var localIP net.IP
    for _, addr := range addrs {
        if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                localIP = ipnet.IP
                break
            }
        }
    }

    // Scan devices in local network
    fmt.Printf("Scanning devices in local network (%s.x):\n", localIP)
    for i := 1; i <= 255; i++ {
        ip := fmt.Sprintf("%s.%d", localIP, i)
        go func(ip string) {
            if deviceInfo := getDeviceInfo(ip); deviceInfo != nil {
                fmt.Printf("Device found: IP=%s, Hostname=%s, OS=%s\n", deviceInfo.ip, deviceInfo.hostname, deviceInfo.os)
            }
        }(ip)
    }

    // Wait for scan to complete
    fmt.Println("Scan complete. Press Enter to exit.")
    fmt.Scanln()
}

type deviceInfo struct {
    ip       string
    hostname string
    os       string
}

func getDeviceInfo(ip string) *deviceInfo {
    // Get hostname
    hostnames, err := net.LookupAddr(ip)
    if err != nil {
        return nil
    }
    hostname := hostnames[0]

    // Get operating system information
    osInfo := runtime.GOOS

    return &deviceInfo{
        ip:       ip,
        hostname: hostname,
        os:       osInfo,
    }
}
