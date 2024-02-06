package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zelas91/goph-keeper/internal/server/models"
	"log"
	"net/http"
	"net/url"
)

var (
	secretKey   = ""
	buildCommit = "N/A"
	buildDate   = "N/A"
)

func main() {
	//fmt.Printf("client build data (%s) version (%s) ---(%s)---\n", buildDate, buildCommit, secretKey)
	//crypt, err := crypto.NewEncrypt(secretKey)
	//fmt.Println(crypt, err)
	//client.NewClient("localhost:8080").Start()
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/api/download"}

	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("jwt", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDcyMzE1MjgsIkxvZ2luIjoiemVsYXMifQ.ptUz9v03OdhkPkaQiXiG4vL_VDFrRb0wFyKyimE70mU")
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Fatal("Не удалось подключиться к серверу:", err)
	}
	defer conn.Close()

	fileData := models.BinaryFile{
		FileName: "goph-keeper-main.zip",
	}

	if err = conn.WriteJSON(fileData); err != nil {
		log.Fatal("Ошибка преобразования в JSON:", err)
	}
	var answer models.AnswerBinaryFile
	if err = conn.ReadJSON(&answer); err != nil {
		fmt.Println("error", err)
		return

	}
	if !answer.Confirm {
		return
	}
	for {
		mt, b, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("ERROR ", mt, err)
			return
		}
		fmt.Println(mt, b)

	}

	//file, err := os.Open("C:/Users/zelas/Downloads/goph-keeper-main.zip")
	//if err != nil {
	//	log.Fatal("Не удалось открыть файл:", err)
	//}
	//defer file.Close()
	//

	//buffer := make([]byte, 1024)
	//for {
	//	n, err := file.Read(buffer)
	//	if err != nil {
	//		if err != io.EOF {
	//			log.Fatal("Ошибка чтения файла:", err)
	//		}
	//		conn.WriteMessage(websocket.TextMessage, []byte("Binary data transfer completed"))
	//		break // Достигнут конец файла
	//	}
	//	err = conn.WriteMessage(websocket.BinaryMessage, buffer[:n])
	//	if err != nil {
	//		log.Fatal("Ошибка отправки сообщения:", err)
	//	}
	//}
	//_, _, err = conn.ReadMessage()
	//if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
	//	fmt.Println("Завершился с ошибкой", err)
	//	return
	//}
	//fmt.Println("Файл успешно отправлен.")

}
