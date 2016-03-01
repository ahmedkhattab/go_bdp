package main

import (
	"bdp/launcher"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"text/template"
	"util"
)

func details(w http.ResponseWriter, r *http.Request) {
	config := util.ConfigStruct()

	t, err := template.ParseFiles("../src/web/static/details.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func deploy(w http.ResponseWriter, r *http.Request) {

	if r.ContentLength == 0 {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	logfile, err := os.Create("../src/web/static/log.out")
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logfile.Close()
	log.SetOutput(logfile)

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
	_, deployed := launcher.LaunchComponents(false, false, config)
	for key, value := range deployed {
		config.Set(key, value)
	}
	http.Redirect(w, r, "/details", http.StatusOK)
	return
}

func index(w http.ResponseWriter, r *http.Request) {
	config := util.InitConfigStruct()
	t, err := template.ParseFiles("../src/web/static/index.html")
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
	fs := http.FileServer(http.Dir("../src/web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", index)
	http.HandleFunc("/deploy", deploy)
	http.HandleFunc("/details", details)

	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}
