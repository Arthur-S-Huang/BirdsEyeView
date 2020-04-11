package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func handleRequests() {
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/labels/{name}", getImageCaption).Methods("GET")
    myRouter.HandleFunc("/caption/{name}", getImageSpeech).Methods("GET")
    myRouter.HandleFunc("/video/{name}", createVideoCaption).Methods("GET")
    myRouter.HandleFunc("/presigned/{file}", getPresignedUrl).Methods("GET")
    log.Fatal(http.ListenAndServe(":19603", myRouter))
}
