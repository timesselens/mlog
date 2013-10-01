package webserver

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"mlog/core"
	"net/http"
	"os"
	"strconv"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fname := "index.html"

	t := template.New(fname)
	if _, err := os.Stat("tmpl/" + fname); err == nil {
		t, err = t.ParseFiles("tmpl/index.html")
		t.ParseGlob("tmpl/*")
	} else {
		for k, v := range templates {
			fmt.Println(k, v)
			if t.Name() == k {
				t, _ = t.Parse(v)
			} else {
				_, _ = t.New(k).Parse(v)
			}
		}
	}

	err := t.Execute(w, core.GetProcessorList())
	if err != nil {
		log.Print(err)
	}
}

func processorHandler(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	fmt.Println(r.Method, r.URL.Path)
	switch r.Method {
	case "GET":
		enc.Encode(core.GetProcessorList())
	case "POST":
		body, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		x, err := core.CreateProcessor(string(body))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
        log.Printf("new proc %T %s\n", x, x)
        core.AppendProcessor(x)
	}
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	fmt.Println(r.Method, r.URL.Path)
	switch r.Method {
	case "GET":
		enc.Encode(core.GetProcessorList())
	case "POST":
		body, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		x, err := core.CreateProcessor(string(body))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
        log.Printf("new proc %T %s\n", x, x)
        core.AppendProcessor(x)
	}
}

func processorDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
    // body, _ := ioutil.ReadAll(r.Body)
    // r.Body.Close()
    // fmt.Println("body: ", string(body))
	i, err := strconv.Atoi(r.Form.Get("i"))
    // fmt.Println(i, err)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if i < 0 {
		http.Error(w, "index i out of bound", 500)
		return
	}
	fmt.Println("deleting processor at", i)
	_, err = core.DeleteProcessor(i)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func pipeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Printf("%#v\n", r.Form)

	i, _ := strconv.Atoi(r.Form.Get("i"))
	j, _ := strconv.Atoi(r.Form.Get("j"))
	if i < 0 || j < 0 {
		http.Error(w, "indexes i and j not provided or out of bound", 500)
		return
	}
	fmt.Println(">>>", i, j)
	err := core.PipeProcessor(i, j)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func HTTPServer(hostport string) {
	fmt.Println("starting http")
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/api/v1/processor/", processorHandler)
	http.HandleFunc("/api/v1/processor/delete", processorDeleteHandler)
	http.HandleFunc("/api/v1/processor/pipe", pipeHandler)
	http.HandleFunc("/api/v1/config", configHandler)
	http.HandleFunc("/app/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	err := http.ListenAndServe(hostport, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

//TODO: inject file here from makefile
var templates = map[string]string{
	"header.html": `<html>`,
	"index.html": `{{template "header"}}
                <h1>hello world</h1>
                {{range .}}
                <div>{{.Name}} <code>{{.Re}}</code></div>
                {{end}}
            {{template "footer"}}
              `,
	"footer.html": `</html>`,
}
