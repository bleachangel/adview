package main

import (
	"adview/delivery"
	"github.com/golang/glog"
	"net/http"
)

func main() {
	//delivery.Run()
	glog.Fatal(http.ListenAndServe(":80", nil))
}
