package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// HTTP request context
type HttpRequestContext struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Response       interface{}
	ResponseType   int
	HTMLEscape     bool
}

// response type
const (
	RspTypeJSON  = 1
	RspTypeJSONP = 2
)

// json response
type JsonResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// NewContext
func NewContext(w http.ResponseWriter, r *http.Request) *HttpRequestContext {
	ctx := &HttpRequestContext{
		Request:        r,
		ResponseWriter: w,
		HTMLEscape:     true,
	}
	return ctx
}

// set json response
func (ctx *HttpRequestContext) SetJsonResponse(res *JsonResponse) error {
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(ctx.ResponseWriter)
	encoder.SetEscapeHTML(ctx.HTMLEscape)
	return encoder.Encode(res)
}

// start server
func StartServer(srv *http.Server, restartWait int) {
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR2)

	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(restartWait)*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

// parse json configuration
func ParseJsonConf(conf interface{}, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("open error:", err)
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(conf)
	if err != nil {
		log.Println("error:", err)
		return err
	}
	log.Printf("conf:%+v\n", conf)
	return nil
}
