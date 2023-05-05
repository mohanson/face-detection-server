package main

/*
#cgo CFLAGS: -Ilibfaced -ISeetaFaceEngine/FaceDetection/include
#cgo LDFLAGS: -Llibfaced -LSeetaFaceEngine/FaceDetection/build -lfaced -lseeta_facedet_lib
#include "faced.h"
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"unsafe"
)

const (
	Version = "0.0.2"
	Port    = "8080"
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

var mutex sync.Mutex

func FaceDetect(path string) *DetectionResult {
	mutex.Lock()
	defer mutex.Unlock()

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
func HandlerRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("<center><h1>200 Service Available</h1></center>"+
		"<hr></hr>"+
		"<center>FaceDetectServer/%s</center>", Version)))
}

func HandlerImageBinDetection(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" && r.Method != "PUT" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	f, err := os.CreateTemp("", "fd_")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	if _, err := io.Copy(f, r.Body); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	f.Close()

	fileinfo, err := os.Stat(f.Name())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if fileinfo.Size() > 1024*1024*8 {
		http.Error(w, "File size limit exceeded [8MB]", http.StatusRequestEntityTooLarge)
		return
	}

	dr := FaceDetect(f.Name())
	if dr.Size[0] == 0 && dr.Size[1] == 0 {
		http.Error(w, "Invalid file operand", 400)
		return
	}
	drctx, _ := json.Marshal(dr)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(drctx))
}

func main() {
	http.HandleFunc("/", HandlerRoot)
	http.HandleFunc("/image/bin/detection", HandlerImageBinDetection)
	log.Printf("main: listen and server on :%s", Port)
	http.ListenAndServe(":"+Port, nil)
}
