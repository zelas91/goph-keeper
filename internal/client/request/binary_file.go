package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	error2 "github.com/zelas91/goph-keeper/internal/client/error"
	"github.com/zelas91/goph-keeper/internal/server/models"
)

type BinaryFile struct {
	request *Request
}

func NewBinaryFile(request *Request) *BinaryFile {
	return &BinaryFile{request: request}
}

func (b *BinaryFile) Delete(args []string) error {
	if len(args) < 1 {
		return error2.ErrInvalidCommand
	}
	url := fmt.Sprintf("/file/%s", args[0])
	resp, err := b.request.R().Delete(url)
	if err != nil {
		return fmt.Errorf("request file delete err: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("request file delete error status code = %d", resp.StatusCode())
	}
	return nil
}

func (b *BinaryFile) Upload(args []string) error {
	if len(args) < 2 {
		return error2.ErrInvalidCommand
	}

	file, err := os.Open(args[1])
	if err != nil {
		return fmt.Errorf("open file err:%v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("file close err: %v\n", err)
		}
	}()
	fInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("file info err:%v", err)
	}
	bf := models.BinaryFile{
		FileName: args[0],
		Size:     int(fInfo.Size()),
	}
	conn, err := b.request.WebsocketConnect("/file/upload")
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Printf("websocket connect close err: %v\n", err)
		}
	}()

	if err = conn.WriteJSON(bf); err != nil {
		return fmt.Errorf("faile encode to json: %v", err)
	}
	var answer models.AnswerBinaryFile
	if err = conn.ReadJSON(&answer); err != nil {
		return fmt.Errorf("read  answer server err %v", err)

	}
	if !answer.Confirm {
		return errors.New("server did not confirm the download of file")
	}
	buffer := make([]byte, 1024)
	counter := 0
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("read file err: %v", err)
			}
			if err := conn.WriteMessage(websocket.TextMessage,
				[]byte("Binary data transfer completed")); err != nil {
				return fmt.Errorf("websocket send msg end file err: %v", err)
			}
			break
		}
		err = conn.WriteMessage(websocket.BinaryMessage, buffer[:n])
		if err != nil {
			return fmt.Errorf("websocket send msg err: %v", err)
		}
		counter = counter + len(buffer[:n])
		fmt.Printf("\rupload: %dKB file size = %dKB", counter/1024, bf.Size/1024)
	}
	fmt.Println("")
	_, _, err = conn.ReadMessage()
	if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
		return fmt.Errorf("transfer completed with err: %v", err)
	}
	return nil
}

func (b *BinaryFile) Download(args []string) error {
	if len(args) < 2 {
		return error2.ErrInvalidCommand
	}
	url := fmt.Sprintf("/file/%s", args[0])
	resp, err := b.request.R().Get(url)
	if err != nil {
		return fmt.Errorf("request file information err: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("request file information error status code = %d", resp.StatusCode())
	}
	var bf models.BinaryFile
	if err := json.Unmarshal(resp.Body(), &bf); err != nil {
		return fmt.Errorf("file information decode err: %v", err)
	}

	conn, err := b.request.WebsocketConnect("/file/download")
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Printf("websocket connect close err: %v\n", err)
		}
	}()

	if err = conn.WriteJSON(bf); err != nil {
		return fmt.Errorf("faile encode to json: %v", err)
	}

	var answer models.AnswerBinaryFile
	if err = conn.ReadJSON(&answer); err != nil {
		return fmt.Errorf("read  answer server err %v", err)

	}
	if !answer.Confirm {
		return errors.New("server did not confirm the download of file")
	}
	file, err := os.Create(fmt.Sprintf("%s/%s", args[1], bf.FileName))
	if err != nil {
		return fmt.Errorf("create file err:%v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("file close err: %v\n", err)
		}
	}()
	counter := 0
	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				break
			}
			return fmt.Errorf("failed to read message: %v", err)
		}
		if mt == websocket.BinaryMessage {
			_, err := file.Write(msg)
			if err != nil {
				return fmt.Errorf("write file err: %v", err)
			}
		}
		counter = counter + len(msg)
		fmt.Printf("\rdownload: %dKB file size = %dKB", counter/1024, bf.Size/1024)
	}
	fmt.Print("\n")
	return nil
}

func (b *BinaryFile) Files(args []string) error {
	resp, err := b.request.R().Get("/file")
	if err != nil {
		return err
	}
	str, err := prettyJSON(resp.Body())
	if err != nil {
		return err
	}
	fmt.Println(str)
	return nil
}
