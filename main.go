package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/reteps/go-akinator"
	"golang.org/x/net/websocket"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Data struct {
	Chat        `json:"Chat"`
	Eightball   `json:"Eightball"`
	Main        `json:"Main"`
	LanguageBot `json:"LanguageBot"`
	ToD         `json:"ToD"`
	Hangman     `json:"Hangman"`
	Akinator    *akinator.Client `json:"Akinator"`
}
type Hangman struct {
	Answer  string
	Stage   int
	Guessed []string `json:"Guessed"`
	Word    []string `json:"Word"`
}
type User struct {
	SentMessages        int
	ConsecutiveCommands int
	SentCommands        int
	DeadChatWins        int
	HangmanWins         int
	IsAdmin             bool
	IsBanned            bool
	Nickname            string
}
type Chat struct {
	Users             map[string]*User `json:"Users"`
	LastMessage       Message          `json:"LastMessage"`
	LastCommandSender string
	IsDeadChat        bool
	IsGroupChat       bool
}
type Eightball struct {
	Messages []string
}
type Main struct {
	HelpMessages        map[string]string `json:"HelpMessages"`
	Messages            map[string]string `json:"Messages"`
	AdultContent        bool
	BlacklistedCommands []string `json:"BlacklistedCommands"`
	WhitelistedCommands []string `json:"WhitelistedCommands"`
	MaxConsecutive      int
}
type ToD struct {
	Truths []string `json:"Truths"`
	Dares  []string `json:"Dares"`
}
type LanguageBot struct {
	Image string
	Words []string `json:"Words"`
	On    bool
}
type Message struct {
	Message   string
	From      string
	Chat      string
	Timestamp float64
	IsCommand bool
}

func isSupportedFileFormat(item string, formats []string) bool {
	parts := strings.Split(item, ".")
	return stringInSlice(parts[len(parts)-1], formats)
}

// sends a message - if it is a file location or url it downloads and sends that instead
func send(sendToUrl, message, to string) error {
	if message == "" {
		return nil
	}
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	var fw io.Writer
	rand.Seed(time.Now().UTC().UnixNano())
	reqUID := strconv.FormatInt(int64(rand.Float64()*1679616), 36)
	values := map[string]string{"hashid": to, "reqUID": reqUID, "text": "", "recipients": ""}
	fileFormats := []string{"png", "jpg", "gif", "mp4", "bmp"}
	if len(message) > 7 && (message[:7] == "/Users/" || message[:6] == "/home/") && isSupportedFileFormat(message, fileFormats) {
		r, err := os.Open(message)
		if err != nil {
			return err
		}
		defer r.Close()
		if fw, err = w.CreateFormFile("file-name", r.Name()); err != nil {
			return err
		}
		_, err = io.Copy(fw, r)
		if err != nil {
			return err
		}
	} else if len(message) > 4 && message[:4] == "http" && isSupportedFileFormat(message, fileFormats) {
		response, err := http.Get(message)
		if err != nil {
			return err
		}
		split_message := strings.Split(strings.Split(message, "://")[1], "/")
		fileName := split_message[len(split_message)-1]
		defer response.Body.Close()
		if fw, err = w.CreateFormFile("file-name", fileName); err != nil {
			return err
		}
		_, err = io.Copy(fw, response.Body)
		if err != nil {
			return err
		}
	} else {
		values["text"] = url.QueryEscape(message)
	}
	for key, value := range values {
		err := w.WriteField(key, value)
		if err != nil {
			panic(err)
		}
	}
	w.Close()
	req, err := http.NewRequest("POST", sendToUrl, &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

// Writes the data to a file
func writesettings(data map[string]*Data) error {
	jsondata, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("./data.json", jsondata, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Reads the data from a file
func readsettings() (map[string]*Data, error) {

	file, err := ioutil.ReadFile("./data.json")

	if err != nil {
		return nil, err
	}
	data := map[string]*Data{}
	err = json.Unmarshal(file, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

// Creates a new chat if it doesnt exist
func chatCreator(data map[string]*Data, event Message) map[string]*Data {
	if _, exists := data[event.Chat]; !exists {
		log.Printf("New chat created %s\n", event.Chat)
		data[event.Chat] = &Data{}
		*data[event.Chat] = *data["defaultChat"]
		data[event.Chat].Chat.Users = map[string]*User{}
		for userName, value := range data["defaultChat"].Chat.Users {
			data[event.Chat].Chat.Users[userName] = &User{}
			*data[event.Chat].Chat.Users[userName] = *value
		}
		if event.Chat[:4] == "chat" {
			data[event.Chat].Chat.IsGroupChat = true
		}

	}
	return data
}

// Creates a new user if it doesnt exist
func userCreator(data *Data, event Message) *Data {
	if _, ok := data.Chat.Users[event.From]; !ok {
		log.Printf("New user created %s: %s\n", event.From, event.Message)
		data.Chat.Users[event.From] = &User{}
		*data.Chat.Users[event.From] = *data.Chat.Users["defaultUser"]
		data.Chat.Users[event.From].Nickname = event.From

	}
	return data
}

// Returns an error if a user doesn't have permissions for a command
func hasPermissions(data *Data, event Message) error {
	lessThanMaxCommands := data.Main.MaxConsecutive == 0 || data.Chat.Users[event.From].ConsecutiveCommands < data.Main.MaxConsecutive || !data.Chat.IsGroupChat
	isBlacklistedCommand := stringInSlice(strings.Split(event.Message, " ")[0], data.Main.BlacklistedCommands)
	isNotWhitelistedCommand := len(data.Main.WhitelistedCommands) > 0 && !stringInSlice(strings.Split(event.Message, " ")[0], data.Main.WhitelistedCommands)
	if data.Chat.Users[event.From].IsAdmin {
		return nil
	}
	if !lessThanMaxCommands {
		return errors.New(data.Main.Messages["MaxCommands"])
	}
	if data.Chat.Users[event.From].IsBanned {
		return errors.New(data.Main.Messages["IsBanned"])
	}
	if isBlacklistedCommand || isNotWhitelistedCommand {
		return errors.New(data.Main.Messages["ChatPermissions"])
	}
	return nil

}

// Handles a message event
func handleEvent(funcmap map[string]interface{}, processingFuncs []interface{}, data map[string]*Data, event Message, url string) {
	data = chatCreator(data, event)
	data[event.Chat] = userCreator(data[event.Chat], event)
	for keyword, function := range funcmap {
		if strings.Split(event.Message, " ")[0] == keyword {
			event.IsCommand = true
			var result string
			var err error
			var message string
			err = hasPermissions(data[event.Chat], event)
			if err == nil {
				event.Message = strings.TrimLeft(event.Message[len(keyword):], " ")
				switch function.(type) {
				case func(*Data, Message) (string, error):
					result, err = function.(func(*Data, Message) (string, error))(data[event.Chat], event)
				case func(*Data, Message) (string, *Data, error):
					result, data[event.Chat], err = function.(func(*Data, Message) (string, *Data, error))(data[event.Chat], event)
				case func(map[string]*Data, Message) (map[string]*Data, error):
					data, err = function.(func(map[string]*Data, Message) (map[string]*Data, error))(data, event)
					result = data[event.Chat].Main.Messages["SuccessMessage"]

				case func(*Data, Message) (*Data, error):
					data[event.Chat], err = function.(func(*Data, Message) (*Data, error))(data[event.Chat], event)
					result = data[event.Chat].Main.Messages["SuccessMessage"]
				default:
					err = errors.New(data[event.Chat].Main.Messages["FunctionError"])
				}
			}
			if err != nil {
				message = fmt.Sprintf("%s%s", data[event.Chat].Main.Messages["ErrorMessage"], err.Error())
			} else {
				message = result
			}

			err = send(fmt.Sprintf("http://%s/sendMessage.srv", url), message, event.Chat)

			if err != nil {
				log.Printf("Send Error:%s", err.Error())
			}
			break
		}
	}
	for _, processingFunc := range processingFuncs {
		switch processingFunc.(type) {
		case func(*Data, Message) (string, *Data):
			var message string
			message, data[event.Chat] = processingFunc.(func(*Data, Message) (string, *Data))(data[event.Chat], event)
			err := send(fmt.Sprintf("http://%s/sendMessage.srv", url), message, event.Chat)
			if err != nil {
				log.Printf("Send Error:%s", err.Error())
			}
		default:
			data[event.Chat] = processingFunc.(func(*Data, Message) *Data)(data[event.Chat], event)
		}
	}
	writesettings(data)
}

func main() {
	data, err := readsettings()
	if err != nil {
		log.Fatal(err)
	}
	url := "192.168.1.15:333"
	ws, err := websocket.Dial(fmt.Sprintf("ws://%s/service", url), "", "http://localhost/")
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not connect to ipod on %s.", url))
	}
	log.Println("ios-imessage-bot started.")
	for {
		var event []map[string]interface{}
		websocket.JSON.Receive(ws, &event)
		if event[0]["event"].(string) == "newMSG" && event[0]["messageParts"].([]interface{})[0].(map[string]interface{})["type"].(string) == "text" {
			message := Message{event[0]["messageParts"].([]interface{})[0].(map[string]interface{})["text"].(string), event[0]["particID"].(string), event[0]["recipHashID"].(string), event[0]["date"].(float64), false}
			go handleEvent(funcmap, processingFuncs, data, message, url)
		}
	}
}
