package main


import "net/http"
import "io"
// import "log"


func initControl() {
	http.HandleFunc("/search_hash_info", SearchHashInfo)
}


func SearchHashInfo(w http.ResponseWriter, req *http.Request) {
	info_hash := req.FormValue("info_hash")

	infoHashQueryChan <- info_hash
	io.WriteString(w, "Done\n");
}