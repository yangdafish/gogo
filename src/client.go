package main

import "flag"
import "fmt"
import "strings"
import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	// "net/url"
    "github.com/gorilla/mux"
)

var connectedClients []int
var stringSet map[string]struct{}
var name string
var port int
var delimitedNodeString string
var delimitedNodes []string
var connectedHostsMap map[string]string
var sourceUrl string


// func ConstructUrl(location string) (url string){
// 	return 
// }

func WhisperMessage(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	message := r.URL.Query().Get("message")
	targetUrl := fmt.Sprintf("http://%s/recieve", sourceUrl)

	fmt.Printf("[%s] Message sent to %s at %s :: %s\n", name, name, targetUrl, message)

	whisperRequest := BuildWhisperRequest(targetUrl, sourceUrl, message)
	SendMessage(whisperRequest)
}

func RecieveMessage(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("message")
	senderName := r.URL.Query().Get("name")

	fmt.Printf("[%s] Message received from %s at %s :: %s\n", name, senderName, sourceUrl, message)
}

func RequestInfo(nodeUrl string) {
	requestUrl := fmt.Sprintf("http://%s/info", nodeUrl)
	req, createReqErr := http.NewRequest(
		"GET",
		requestUrl,
		nil,
	)
	if createReqErr != nil {
		panic(createReqErr) // NewRequest only errors on bad methods or un-parsable urls
	}
	fmt.Println("Sending Request Info", req)
	response := SendMessage(req)

	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }
 
    responseString := string(responseData)

	connectedHostsMap[nodeUrl] = responseString
	fmt.Printf("[%s] Connected to %s at %s\n", name, responseString, nodeUrl)

	// fmt.Println("Got Response", responseString)
}

func ReturnInfo(w http.ResponseWriter, r *http.Request) {
	nameJson, err := json.Marshal(name)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(nameJson)
}

func BuildWhisperRequest(targetUrl string, requestSender string, requestMessage string) (req *http.Request){
	req, createReqErr := http.NewRequest(
		"GET",
		targetUrl,
		nil,
	)
	if createReqErr != nil {
		panic(createReqErr) // NewRequest only errors on bad methods or un-parsable urls
	}

	q := req.URL.Query()
	q.Add("name", requestSender)
	q.Add("message", requestMessage)
	req.URL.RawQuery = q.Encode()
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return req
}

func SendMessage(req *http.Request) (res *http.Response){
	res, sendReqErr := http.DefaultClient.Do(req)
	if sendReqErr != nil {
		panic(sendReqErr) // NewRequest only errors on bad methods or un-parsable urls
	}
	return res
}

func main() {
	flag.StringVar(&name, "name", "bar", "a string var")
	flag.IntVar(&port, "port", 8000, "an int")
	bootNodesPtr := flag.String("bootnodes", "", "a comma delimited string")
	flag.Parse()

	// fmt.Println("name:", name)
	// fmt.Println("port:", port)
	if len(*bootNodesPtr) > 0 {
		connectedHostsMap = make(map[string]string)
		delimitedNodes := strings.Split(*bootNodesPtr, ",")

		fmt.Println("Foudn Nodes", len(delimitedNodes))
		for _, node := range delimitedNodes {
			fmt.Println("Ya here pla", delimitedNodes)
			RequestInfo(node)
		}
	}

	portSetting := fmt.Sprintf(":%d", port)
	sourceUrl = fmt.Sprintf("localhost%s", portSetting)

	router := mux.NewRouter()
	router.HandleFunc("/whisper", WhisperMessage).Methods("GET")
	router.HandleFunc("/info", ReturnInfo).Methods("GET")
	router.HandleFunc("/recieve", RecieveMessage).Methods("GET")

	fmt.Println("Listenging  on ", sourceUrl)
	log.Fatal(http.ListenAndServe(portSetting, router))
}
