package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

//error codes
const ERR_UNRECOVERABLE = "unrecoverable-error"
const ERR_INVALID_USER_INPUT = "invalid-user-input"

//statuses
const SUCCESS = "success"
const ERROR = "error"

func mw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "content-type")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func sendError(ctx context.Context, w http.ResponseWriter, err error, code string) {
	defer TimeTrack(ctx, time.Now())

	res := &Response{
		Status:    "error",
		Msg:       err.Error(),
		ErrorCode: code,
	}
	str, _ := json.Marshal(res)
	fmt.Fprintf(w, string(str))
}

func sendSuccess(ctx context.Context, w http.ResponseWriter, data interface{}) {
	defer TimeTrack(ctx, time.Now())

	res := &Response{
		Status: "ok",
		Data:   data,
	}
	str, err := json.Marshal(res)
	if err != nil {
		e := errors.New("Unrecoverable error")
		sendError(ctx, w, e, ERR_UNRECOVERABLE)
		return
	}
	fmt.Fprint(w, string(str))
}

func TimeTrack(ctx context.Context, start time.Time) {
	elapsed := time.Since(start)

	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(1)

	// Retrieve a function object this functions parent.
	funcObj := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path).
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

	//log.Print(fmt.Sprintf("%s took %s", name, elapsed))
	Dbg(ctx, fmt.Sprintf("%s took %s", name, elapsed))
}

func Dbg(ctx context.Context, v string) {
	_, fl, line, _ := runtime.Caller(1)
	log.WithFields(log.Fields{
		"fl": fl + ":" + strconv.Itoa(line),
	}).Debug(v)
}
