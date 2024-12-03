package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const GracePeriod = 3 * time.Second

func main() {
	if err := Start(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func Start() error {

	// listen to OS signals and gracefully shutdown server
	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		close(stopped)
	}()

	//Server Setup
	fmt.Println("Starting TCP Server")
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8080))
	if err != nil {
		return err //unalbe to start listener
	}
	// Create a context with a cancel function
	ctx, cancel := context.WithCancel(context.Background())
	go StartAcceptingConnections(ctx, listener)

	//Start test client
	go StartTestClient("1")

	//wait for signal to shutdown
	<-stopped
	fmt.Println("Starting graceful shutdown")

	//calling cancel to notify services to shutdown
	cancel()

	//start client after cancel should not be able to connect
	go StartTestClient("2")

	//allowing other services a grace perdiod to shutdown
	time.Sleep(GracePeriod)
	fmt.Println("Finished graceful shutdown")
	return nil
}

func StartAcceptingConnections(ctx context.Context, listener net.Listener) {
	defer listener.Close()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := listener.Accept() //accepting new connection
			if err != nil {
				fmt.Println("Error accepting connection:", err)
			} else {
				go handleConnection(ctx, conn)
			}
		}
	}
}

func handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for {
		select {
		case <-ctx.Done():
			scanner.Scan()
			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading from connection:", err)
				fmt.Fprintf(conn, "%s\n", "Transaction Error")
				return //exiting connection on error
			} else {
				fmt.Fprintf(conn, "%s\n", "Transaction Cancelled")
			}

		default:
			scanner.Scan()
			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading from connection:", err)
				fmt.Fprintf(conn, "%s\n", "Transaction Error")
				return //exiting connection on error
			} else {
				fmt.Fprintf(conn, "%s\n", "Transaction Accepted")
			}
		}
	}
}

// FOR TESTING
func StartTestClient(id string) {
	fmt.Printf("Starting Client %s\n", id)
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error dialing into TCP connection:", err)
		return
	}
	defer conn.Close()

	for {
		_, err := conn.Write([]byte("PAYMENT|10\n"))
		if err != nil {
			fmt.Println("Error writing to server:", err)
			return
		}

		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}
		fmt.Printf("Client %s received: %s", id, response)
		time.Sleep(1 * time.Second)
	}
}
