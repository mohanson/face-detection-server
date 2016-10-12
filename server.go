package main

/*
#cgo CFLAGS: -I./libfaced -I/usr/include
#cgo LDFLAGS: -L./libfaced -lfaced
#include "faced.h"
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"unsafe"
)

const (
	Version = "0.0.1"
	Port    = "8090"
)

type DetectionResultFace struct {
	X      int `json:"x"`
	Y      int `josn:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type DetectionResult struct {
	Size []int                 `json:"size"`
	Face []DetectionResultFace `json:"face"`
}

func FaceDetect(path string) *DetectionResult {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	cresp := C.FaceDetect(cpath)
	defer C.free(unsafe.Pointer(cresp))

	resp := C.GoString(cresp)

	dr := &DetectionResult{}
	json.Unmarshal([]byte(resp), dr)

	return dr
}

// ============================================================================
// HTTP Server
// ============================================================================

type Handler struct {
	mux map[string]func(http.ResponseWriter, *http.Request)
}

func NewHandler() *Handler {
	Handler := &Handler{}
	Handler.mux = make(map[string]func(http.ResponseWriter, *http.Request))
	return Handler
}

func (Handler *Handler) Bind(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	Handler.mux[pattern] = handler
}

func (Handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := Handler.mux[r.URL.String()]; ok {
		log.Println(r.Method, r.URL.String())
		h(w, r)
		return
	}
	w.WriteHeader(404)
}

func HandlerRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("<center><h1>200 Service Available</h1></center>"+
		"<hr></hr>"+
		"<center>FaceDetectServer/%s</center>", Version)))
}

func HandlerDetectionUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" && r.Method != "PUT" {
		w.WriteHeader(405)
		return
	}
	defer r.Body.Close()

	f, err := ioutil.TempFile("", "fdserver_")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_, err = io.Copy(f, r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		f.Close()
		return
	}
	f.Close()
	defer os.Remove(f.Name())

	fileinfo, err := os.Stat(f.Name())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if fileinfo.Size() > 1024*1024*8 {
		http.Error(w, "File size limit exceeded [8MB]", 400)
		return
	}

	dr := FaceDetect(f.Name())
	if dr.Size[0] == 0 && dr.Size[1] == 0 {
		w.WriteHeader(400)
		w.Write([]byte("Invalid file operand"))
		return
	}
	drs, _ := json.Marshal(dr)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(drs))
	return
}

func main() {
	Handler := NewHandler()
	Handler.Bind("/", HandlerRoot)
	Handler.Bind("/detection/upload", HandlerDetectionUpload)
	server := http.Server{
		Addr:    ":" + Port,
		Handler: Handler,
	}
	server.ListenAndServe()
}
