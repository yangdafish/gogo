package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
    "strconv"
    "github.com/gorilla/mux"
)


var ClientName string
var ClientUrl string
var ConnectedHostsMap map[string]string
var Port int


func WhisperMessageHandler(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("message")
	targetName := string(r.URL.Query().Get("name"))

	if targetLocation, ok := ConnectedHostsMap[targetName]; ok {
		fmt.Printf("[%s] Message sent to %s at %s :: %s\n", ClientName, targetName, targetLocation, message)
	
		whisperRequest := BuildHttpWhisperRequest(targetLocation, ClientName, message)
		SendMessage(whisperRequest)
	} else {
		panicMessage := fmt.Sprintf("Did not find target %s in connected hosts", targetName)
		http.Error(w, panicMessage, http.StatusBadRequest)
	}
}

func RecieveMessageHandler(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("message")
	senderName := r.URL.Query().Get("name")
	if senderUrl, ok := ConnectedHostsMap[senderName]; ok {
		fmt.Printf("[%s] Message received from %s at %s :: %s\n", ClientName, senderName, senderUrl, message)
	} else {
		panicMessage := fmt.Sprintf("Did not find source host %s in connected hosts", senderName)
		http.Error(w, panicMessage, http.StatusBadRequest)
	}
}

func ConnectNodeHandler(w http.ResponseWriter, r *http.Request) {
	requestNodeName := r.URL.Query().Get("name")
	requestNodeUrl := r.URL.Query().Get("url")
	
	if ConnectedHostsMap != nil {
		ConnectedHostsMap[requestNodeName] = requestNodeUrl
		fmt.Printf("[%s] Connected to %s at %s\n", ClientName, requestNodeName, requestNodeUrl)
	} else {
		http.Error(w, "Host Map Not Initialized", http.StatusBadRequest)
	}

	nameJson, err := json.Marshal(ClientName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(nameJson)
}

func RequestConnectNode(nodeUrl string) {
	req := BuildHttpConnectNodeRequest(nodeUrl)
	response := SendMessage(req)

	if response.StatusCode == http.StatusOK {
		defer response.Body.Close()

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err) // panic vs exception?
		}
	 
		responseString, _ := strconv.Unquote(string(responseData))
		ConnectedHostsMap[responseString] = nodeUrl
		fmt.Printf("[%s] Connected to %s at %s\n", ClientName, responseString, nodeUrl)
	} else {
		panicMessage := fmt.Sprintf("Recieved RequestConnectNode response with status code %d", response.StatusCode)
		panic(panicMessage)
	}
}

func SendMessage(req *http.Request) (res *http.Response) {
	res, sendReqErr := http.DefaultClient.Do(req)
	if sendReqErr != nil {
		panic(sendReqErr)
	}
	return res
}

func BuildHttpConnectNodeRequest(targetUrl string) (req *http.Request) {
	requestUrl := fmt.Sprintf("http://%s/connect", targetUrl)
	req, createReqErr := http.NewRequest(
		"GET",
		requestUrl,
		nil,
	)
	if createReqErr != nil {
		panic(createReqErr)
	}

	q := req.URL.Query()
	q.Add("name", ClientName)
	q.Add("url", ClientUrl)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = q.Encode()

	return req
}

func BuildHttpWhisperRequest(targetUrl string, requestSender string, requestMessage string) (req *http.Request) {
	requestUrl := fmt.Sprintf("http://%s/recieve", targetUrl)
	req, createReqErr := http.NewRequest(
		"GET",
		requestUrl,
		nil,
	)
	if createReqErr != nil {
		panic(createReqErr)
	}

	q := req.URL.Query()
	q.Add("name", requestSender)
	q.Add("message", requestMessage)
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = q.Encode()

	return req
}

func main() {
	ConnectedHostsMap = make(map[string]string)

	flag.StringVar(&ClientName, "name", "bar", "a string var")
	flag.IntVar(&Port, "port", 8000, "an int")
	bootNodesPtr := flag.String("bootnodes", "", "a comma delimited string")
	flag.Parse()

	portSetting := fmt.Sprintf(":%d", Port)
	ClientUrl = fmt.Sprintf("localhost%s", portSetting)

	if len(*bootNodesPtr) > 0 {
		delimitedNodes := strings.Split(*bootNodesPtr, ",")

		for _, node := range delimitedNodes {
			RequestConnectNode(node)
		}
	}

	router := mux.NewRouter()
	router.HandleFunc("/whisper", WhisperMessageHandler).Methods("GET")
	router.HandleFunc("/connect", ConnectNodeHandler).Methods("GET")
	router.HandleFunc("/recieve", RecieveMessageHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(portSetting, router))
}
