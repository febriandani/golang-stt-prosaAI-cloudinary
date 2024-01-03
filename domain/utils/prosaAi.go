package utils

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

type ResultTranscript struct {
	TranscriptData string
}

func SpeechToTextApi(ctx context.Context, filename string) ResultTranscript {
	// Setup
	url := viper.GetString("PROSA.URL")
	apiKey := viper.GetString("PROSA.API_KEY")
	// Authenticate via HTTP Header
	headers := make(map[string][]string)
	headers["x-api-key"] = []string{apiKey}

	// Connect to the WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(url, headers)
	if err != nil {
		log.Println("error stt 1 ", err)
		panic(err)
	}
	defer conn.Close()

	// Configure the session
	config := map[string]interface{}{
		"label":           "Streaming Test",
		"model":           "stt-general-online",
		"include_partial": false,
	}
	err = conn.WriteJSON(config)
	if err != nil {
		log.Println("error stt 2 ", err)
		panic(err)
	}

	// Concurrently send audio data and receive results
	done := make(chan struct{})
	go func() {
		defer close(done)
		sendAudio(filename, conn)
	}()
	result := receiveMessage(conn)

	return result

}

func sendAudio(filename string, conn *websocket.Conn) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	chunkSize := 16000
	buffer := make([]byte, chunkSize)

	for {
		n, err := file.Read(buffer)
		if err != nil {
			break
		}

		err = conn.WriteMessage(websocket.BinaryMessage, buffer[:n])
		if err != nil {
			panic(err)
		}
	}

	// Signify the end of audio stream
	err = conn.WriteMessage(websocket.BinaryMessage, []byte{})
	if err != nil {
		panic(err)
	}
}

func receiveMessage(conn *websocket.Conn) ResultTranscript {

	var result ResultTranscript

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("error read msg ", err)
			panic(err)
		}

		var data map[string]interface{}
		err = json.Unmarshal(message, &data)
		if err != nil {
			log.Println("error unmarshal msg ", err)
			panic(err)
		}

		messageType := data["type"].(string)

		switch messageType {
		case "result":
			transcript := data["transcript"].(string)
			log.Println("result 1 : ", transcript)
			result.TranscriptData = transcript
			return result
			// Process final transcript
		case "partial":
			transcript := data["transcript"].(string)
			// Process partial transcript
			log.Println("partial : ", transcript)
		}

		// Print received data
	}
}
