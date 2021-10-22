package main

import (
	"github.com/alexstan12/LearnGobyTests/dependencyInj"
	"log"
	"net/http"
)

func main(){
	//dependencyInj.Greet(os.Stdout, "Elodie")

	log.Fatal(http.ListenAndServe(":5000", http.HandlerFunc(dependencyInj.MyGreeterHandler)))
}
