package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	vegeta "github.com/tsenart/vegeta/lib"

	"time"
)

var uuids = [50]string{
	"00000000-0000-0000-0000-000000000001",
	"00000000-0000-0000-0000-000000000002",
	"00000000-0000-0000-0000-000000000003",
	"00000000-0000-0000-0000-000000000004",
	"00000000-0000-0000-0000-000000000005",
	"00000000-0000-0000-0000-000000000006",
	"00000000-0000-0000-0000-000000000007",
	"00000000-0000-0000-0000-000000000008",
	"00000000-0000-0000-0000-000000000009",
	"00000000-0000-0000-0000-000000000010",
	"00000000-0000-0000-0000-000000000011",
	"00000000-0000-0000-0000-000000000012",
	"00000000-0000-0000-0000-000000000013",
	"00000000-0000-0000-0000-000000000014",
	"00000000-0000-0000-0000-000000000015",
	"00000000-0000-0000-0000-000000000016",
	"00000000-0000-0000-0000-000000000017",
	"00000000-0000-0000-0000-000000000018",
	"00000000-0000-0000-0000-000000000019",
	"00000000-0000-0000-0000-000000000020",
	"00000000-0000-0000-0000-000000000021",
	"00000000-0000-0000-0000-000000000022",
	"00000000-0000-0000-0000-000000000023",
	"00000000-0000-0000-0000-000000000024",
	"00000000-0000-0000-0000-000000000025",
	"00000000-0000-0000-0000-000000000026",
	"00000000-0000-0000-0000-000000000027",
	"00000000-0000-0000-0000-000000000028",
	"00000000-0000-0000-0000-000000000029",
	"00000000-0000-0000-0000-000000000030",
	"00000000-0000-0000-0000-000000000031",
	"00000000-0000-0000-0000-000000000032",
	"00000000-0000-0000-0000-000000000033",
	"00000000-0000-0000-0000-000000000034",
	"00000000-0000-0000-0000-000000000035",
	"00000000-0000-0000-0000-000000000036",
	"00000000-0000-0000-0000-000000000037",
	"00000000-0000-0000-0000-000000000038",
	"00000000-0000-0000-0000-000000000039",
	"00000000-0000-0000-0000-000000000040",
	"00000000-0000-0000-0000-000000000041",
	"00000000-0000-0000-0000-000000000042",
	"00000000-0000-0000-0000-000000000043",
	"00000000-0000-0000-0000-000000000044",
	"00000000-0000-0000-0000-000000000045",
	"00000000-0000-0000-0000-000000000046",
	"00000000-0000-0000-0000-000000000047",
	"00000000-0000-0000-0000-000000000048",
	"00000000-0000-0000-0000-000000000049",
	"00000000-0000-0000-0000-000000000050",
}
var reqCounter int64

func getRandomDate() string {
	now := time.Now()
	randDay := atomic.AddInt64(&reqCounter, 1) % 8
	finalDate := now.AddDate(0, 0, int(randDay))
	return finalDate.Format("2006-01-02")
}
func getRandomUUID() string {
	randInd := atomic.AddInt64(&reqCounter, 1) % 50
	return uuids[randInd]
}

func main() {

	client := http.Client{}
	req := struct {
		Role string `json:"role"`
	}{
		Role: "user",
	}

	body, _ := json.Marshal(req)
	request, _ := http.NewRequest("POST", "http://localhost:8080/dummyLogin", bytes.NewBuffer(body))
	resp, _ := client.Do(request)
	defer resp.Body.Close()
	var token struct {
		Token string `json:"token"`
	}

	json.NewDecoder(resp.Body).Decode(&token)
	targeter := func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		roomID := getRandomUUID()
		dateStr := getRandomDate()
		if tgt.Header == nil {
			tgt.Header = make(http.Header)
		}
		tgt.Header.Add("Authorization", "Bearer "+token.Token)
		tgt.Method = "GET"
		tgt.URL = fmt.Sprintf("http://localhost:8080/rooms/%s/slots/list?date=%s", roomID, dateStr)
		return nil
	}
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	rate := vegeta.Rate{Freq: 100, Per: time.Second}
	duration := 30 * time.Second
	for res := range attacker.Attack(targeter, rate, duration, "Dynamic Load Test") {
		metrics.Add(res)
	}
	metrics.Close()
	reporter := vegeta.NewTextReporter(&metrics)
	reporter.Report(os.Stdout)
}
