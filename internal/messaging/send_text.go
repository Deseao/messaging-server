package messaging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type smsMessage struct {
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
}

func Send(accountID, authToken, fromNumber, toNumber string) {
	payload := smsMessage{
		From: fromNumber,
		To:   toNumber,
		Text: "Wash your hands!",
	}

	byts, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("error:", err)
	}

	reader := bytes.NewReader(byts)

	newRequest, err := http.NewRequest(http.MethodPost, "https://www.freeclimb.com/apiserver/Accounts/"+accountID+"/Messages", reader)
	if err != nil {
		fmt.Println("error:", err)
	}

	newRequest.Header.Add("Content-Type", "application/json")
	newRequest.SetBasicAuth(accountID, authToken)

	rsp, err := http.DefaultClient.Do(newRequest)
	if err != nil {
		fmt.Println("error:", err)
	}

	if rsp.StatusCode != http.StatusAccepted {
		fmt.Println("Response Bad:", rsp.StatusCode)
		defer rsp.Body.Close()
		b, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Printf("%s", b)
	}
}
