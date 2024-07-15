package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	for i := 0; i < 5; i++ {
		message := fmt.Sprintf("Hello %d\n", i)
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error writing to connection:", err)
			return
		}
		fmt.Printf("Sent message: %s\n", message) // Debug statement
		time.Sleep(1 * time.Second)
	}
}