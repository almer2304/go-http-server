package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func main(){
	listener, err := net.Listen("tcp", ":1000")
	if err != nil{
		log.Fatal("Gagal terhubung ke TCP: ", err)
	}

	defer listener.Close()
	fmt.Println("Server berhasil terhubung ke port 1000")

	for {
		conn, err := listener.Accept()
		if err != nil{
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn){
	defer conn.Close()

	reader := bufio.NewReader(conn)

	requestLine, err := reader.ReadString('\n')
	if err != nil{
		return
	}
	fmt.Printf("Berhasil menerima request %s", requestLine)

	parts := strings.Fields(requestLine)
	method := parts[0]
	path := parts[1]

	if method != "GET"{
		responseHandler(conn, "405 Method Not Allowed", "text/plain", []byte("Method Not Allowed"))
		return
	}

	if path == "/"{
		path = "/index.html"
	}

	filePath := filepath.Join("public", filepath.Clean(path))
	fileData, err := os.ReadFile(filePath)
	if err != nil{
		log.Printf("File tidak ditemukan: %s", filePath)
		responseHandler(conn, "404", "text/plain", []byte("<h1>404 Not Found</h1><p>File yang Anda cari tidak ada di server ini.</p>"))
		return
	}

	contentType := getContentType(filePath)
	responseHandler(conn, "200 OK", contentType, fileData)
}

func responseHandler(conn net.Conn, status string, contentType string, body []byte){
	formatHeader := fmt.Sprintf(
		"HTTP/1.1 %s\r\n"+
		"Content-Type: %s\r\n"+
		"Content-Length: %d\r\n"+
		"Connection: close\r\n"+
		"\r\n",
		status,
		contentType,
		len(body),
	)

	conn.Write([]byte(formatHeader))
	conn.Write(body)
}

func getContentType(contentType string) string{
	ext := filepath.Ext(contentType)
	switch ext{
	case ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	default:
		return "text/plain"

	}
}