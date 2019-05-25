package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gomoltp/pkg/moltp"
)

type (
	htmlData struct {
		Static string
	}

	infomessage struct {
		Info string `json:"info"`
	}
)

var (
	indexTemplate   *template.Template
	staticFolder    string
	templatesFolder string
	debugOn         bool
)

func init() {
	flag.StringVar(&staticFolder, "static", "/var/www/html/static", "Path to folder holding static files.")
	flag.StringVar(&templatesFolder, "templates", "/var/www/html/templates", "Path to folder holding html pages templates files.")
	flag.BoolVar(&debugOn, "v", false, "Swith for log printing")
}

func fixFolderPath(p string) string {
	p = strings.TrimSuffix(p, "/")

	info, err := os.Stat(p)
	if err != nil {
		log.Fatal(err)
	}

	if !info.IsDir() {
		log.Fatal(fmt.Sprintf("%s is not a directry", p))
	}

	return p
}

func doInit() {
	var err error
	flag.Parse()

	staticFolder = fixFolderPath(staticFolder)
	templatesFolder = fixFolderPath(templatesFolder)

	indexTemplate, err = template.ParseFiles(
		fmt.Sprintf("%s/index.tmpl", templatesFolder),
		fmt.Sprintf("%s/base.tmpl", templatesFolder),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	err := indexTemplate.ExecuteTemplate(w, "base", htmlData{Static: staticFolder})
	if err != nil {
		log.Print(err)
		http.NotFoundHandler()
	}
}

func proofHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Println("bad body", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(infomessage{Info: "Bad body"})
		return
	}

	var formulas moltp.RawFormula
	err = json.Unmarshal(body, &formulas)
	if err != nil {
		log.Println("bad formula", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(infomessage{Info: "Bad formulas"})
		return
	}

	solution, err := moltp.Prove(formulas, debugOn)

	// err = json.NewEncoder(w).Encode(reports)
	if err != nil {
		log.Println("error solving", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(infomessage{Info: fmt.Sprintf("error solving: %s", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(solution)
}

func main() {
	doInit()

	http.HandleFunc("/", index)
	http.HandleFunc("/prover", proofHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticFolder))))

	log.Fatal(http.ListenAndServe("127.0.0.1:4000", nil))
}
