package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

type Message struct{
	from	string
	payload	[]byte
}

type Server struct { 
	listenAddr	string
	ln 			net.Listener
	quitch 		chan struct{}
	msgch 		chan Message
}


func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr:	listenAddr,
		quitch:		make(chan struct{}),
		msgch:      make(chan Message),
	}
	
}


func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return err
	}
	defer ln.Close()
	fmt.Println("Server started on", s.listenAddr)
	s.ln = ln

	go s.acceptLoop()

	<-s.quitch
	close(s.msgch)

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			select {
			case <-s.quitch:
				fmt.Println("Server is shutting down.")
				return
			default:
				fmt.Println("Accept error:", err)
				continue
			}
		}

		fmt.Println("new connection to server from addr: ", conn.RemoteAddr())

		go s.readLoop(conn)

	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer func() {
		fmt.Println("Closing connection from:", conn.RemoteAddr())
		conn.Close()
	}()

	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("connection closed by client: ", conn.RemoteAddr())
				return
			}
			
			fmt.Println("read error:", err)
			return
		
		}
		
		
		s.msgch <- Message{
			from: 		conn.RemoteAddr().String(),
			payload:	buf[:n],
		}
		receivedData := buf[:n]
		fmt.Printf("Received %d bytes from %s\n", n, conn.RemoteAddr())
		fmt.Printf("Received message from %s: %s\n", conn.RemoteAddr(), string(receivedData))
	}
}

func main() {
	server := NewServer(":3000")
	fmt.Println("Starting server...")
	server.Start()

	go func() {
		if err := server.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() { 
		for msg := range server.msgch{

			fmt.Println(string(msg.payload))
			fmt.Printf("received message from connection: (%s):%s\n ", msg.from, string(msg.payload))
		}
		}()
	select {}
}	
