package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Please provide host:port")
		return
	}
	address := os.Args[1]

	tlsConfig := &tls.Config{}
	if len(os.Args) == 3 && os.Args[2] == "--tls-skip-verify" {
		// Skip certificate verification (useful for self-signed certs; avoid in production)
		tlsConfig.InsecureSkipVerify = true
	}

	conn, err := tls.Dial("tcp", address, tlsConfig)
	if err != nil {
		fmt.Printf("Failed to establish TLS connection: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Successfully connected to TLS server %s\n", address)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")

		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Failed to read from stdin: %v\n", err)
			continue
		}

		if _, err = fmt.Fprintf(conn, text); err != nil {
			fmt.Printf("Failed to write to server: %v\n", err)
			continue
		}

		serverReader := bufio.NewReader(conn)

		msgStrBuilder := strings.Builder{}
		msgStrBuilder.WriteString("->: ")

		for {
			b, err := serverReader.ReadByte()
			if err != nil {
				if err == io.EOF {
					fmt.Println("Connection closed by server")
					return
				}
				fmt.Printf("Failed to read from server: %v\n", err)
				break
			}

			msgStrBuilder.WriteByte(b)
			if serverReader.Buffered() == 0 {
				fmt.Print(msgStrBuilder.String())
				break
			}
		}
	}
}
