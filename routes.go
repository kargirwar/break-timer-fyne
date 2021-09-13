package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Status    string      `json:"status"`
	Msg       string      `json:"msg"`
	ErrorCode string      `json:"error-code"`
	Data      interface{} `json:"data"`
}

func getQueryParams(r *http.Request) (string, error) {
	params := r.URL.Query()

	rules, present := params["rules"]
	if !present || len(rules) == 0 {
		e := errors.New("Rules not provided")
		return "", e
	}

	return rules[0], nil
}

func setTimers(w http.ResponseWriter, r *http.Request) {
	rules, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendError(r.Context(), w, err, ERR_INVALID_USER_INPUT)
		return
	}

	timerCh <- string(rules)
	ioutil.WriteFile(getOsFilePath(SETTINGS_FILE), rules, 0644)

	sendSuccess(r.Context(), w, nil)
}

func stop(w http.ResponseWriter, r *http.Request) {
	fmt.Println("stop")
	playerCh <- STOP
	sendSuccess(r.Context(), w, nil)
}
