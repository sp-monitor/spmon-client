package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
)


func main() {
	go startModules()
	startUDPServer()
}

func startModules() {
	modulePath := "./modules/spmon-module-cpu"
    absPath, err := filepath.Abs(modulePath)
    if err != nil {
        fmt.Println("Error getting absolute path:", err)
        return
    }

    cmd := exec.Command(absPath, "-h", "localhost", "-p", "9559")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    err = cmd.Run()
    if err != nil {
        fmt.Println("Error running executable:", err)
        return
    }
}

func startUDPServer() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:9559")
	if err != nil {
			fmt.Println("Error resolving UDP address:", err)
			return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
			fmt.Println("Error listening on UDP:", err)
			return
	}

	defer conn.Close()

	fmt.Println("UDP server listening on", addr)


	// Handle SIGINT signal to gracefully close the UDP connection
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT)

	go func() {
			<-sigint
			fmt.Println("Received SIGINT signal, closing UDP connection...")
			conn.Close()
			os.Exit(0)
	}()

	for {
			buf := make([]byte, 1024)
			n, addr, err := conn.ReadFromUDP(buf)
			if err != nil {
					fmt.Println("Error reading from UDP:", err)
					continue
			}
			fmt.Printf("Received %d bytes from %s: %s\n", n, addr, string(buf[:n]))
			go sendToApp(buf)
	}
}