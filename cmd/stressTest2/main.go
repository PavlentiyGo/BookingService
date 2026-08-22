package main

import (
	"encoding/json"
	"os"

	vegeta "github.com/tsenart/vegeta/lib"

	"time"
)

func main() {

	req := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    "emaillll@gmail.com",
		Password: "password",
	}
	body, _ := json.Marshal(req)

	targeter := func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		tgt.Method = "POST"
		tgt.Body = body
		tgt.URL = "http://localhost:8080/login"
		return nil
	}
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	rate := vegeta.Rate{Freq: 50, Per: time.Second}
	duration := 30 * time.Second
	for res := range attacker.Attack(targeter, rate, duration, "Dynamic Load Test") {
		metrics.Add(res)
	}
	metrics.Close()
	reporter := vegeta.NewTextReporter(&metrics)
	reporter.Report(os.Stdout)
}
