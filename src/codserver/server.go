package main

import (
	"codcli/launcher"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"text/template"
	"util"
)

type pageData map[string]interface{}

func details(w http.ResponseWriter, r *http.Request) {
	config := util.ConfigStruct()
	statuses := launcher.ComponentsStatuses()
	fmt.Println(statuses)
	t, err := template.ParseFiles("../src/codserver/static/details.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, pageData{"config": config, "statuses": statuses})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func deploy(w http.ResponseWriter, r *http.Request) {
	logfile, err := os.Create("../src/codserver/static/log.out")
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logfile.Close()
	log.SetOutput(logfile)
	if r.ContentLength == 0 {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	log.Println("Starting deployment ...")
	config := util.ConfigStruct()
	s := reflect.ValueOf(&config).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		fieldType := field.Kind()
		fieldName := typeOfT.Field(i).Name
		if fieldType == reflect.Int {
			val, _ := strconv.Atoi(r.FormValue(fieldName))
			field.SetInt(int64(val))
		}
		if fieldType == reflect.Bool {
			val, _ := strconv.ParseBool(r.FormValue(fieldName))
			field.SetBool(val)
		}
		if fieldType == reflect.String {
			field.SetString(r.FormValue(fieldName))
		}
		fmt.Println(s.Field(i))
	}
	launcher.LaunchComponents(false, false, config)
	http.Redirect(w, r, "/details", http.StatusFound)
	return
}

func index(w http.ResponseWriter, r *http.Request) {
	config := util.InitConfigStruct()
	t, err := template.ParseFiles("../src/codserver/static/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	util.SetDefaultConfig()
	fs := http.FileServer(http.Dir("../src/codserver/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", index)
	http.HandleFunc("/deploy", deploy)
	http.HandleFunc("/details", details)

	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}
