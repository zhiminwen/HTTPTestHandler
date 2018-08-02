package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"text/template"
	"time"
)

func greet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "Hello World! %s", time.Now().Format("2006-01-02 15:04:05"))
}

func genLoad(load float64, period int, wg *sync.WaitGroup) {
	defer wg.Done()
	now := time.Now()
	endTime := now.Add(time.Duration(period) * time.Second)
	for {
		if time.Now().After(endTime) {
			break
		}
		ms := time.Now().UnixNano() / 1e+6 //ms
		if ms%100 == 0 {
			//every 100 ms, sleep (1-load) %
			time.Sleep(time.Duration((1-load)*100) * time.Millisecond)
		}
	}
}

func cpuHogging(w http.ResponseWriter, r *http.Request) {
	load, err := strconv.ParseFloat(r.FormValue("load"), 64)
	if err != nil {
		http.Error(w, "Could not parse load form value", http.StatusBadRequest)
		return
	}
	period, err := strconv.Atoi(r.FormValue("period"))
	if err != nil {
		http.Error(w, "Could not parse lasting period form value", http.StatusBadRequest)
		return
	}

	log.Printf("start: load=%0.2f, period=%d sec, CPU=%d", load, period, runtime.NumCPU())
	var wg sync.WaitGroup
	for index := 0; index < runtime.NumCPU(); index++ {
		wg.Add(1)
		go genLoad(load, period, &wg)
	}

	wg.Wait()
	log.Printf("done: load=%0.2f, period=%d sec, CPU=%d", load, period, runtime.NumCPU())
	fmt.Fprintf(w, "load=%0.2f, period=%d sec, CPU=%d", load, period, runtime.NumCPU())
}

func timeHogging(w http.ResponseWriter, r *http.Request) {
	period, err := strconv.Atoi(r.FormValue("period"))
	if err != nil {
		http.Error(w, "Could not parse lasting period form value", http.StatusBadRequest)
		return
	}

	time.Sleep(time.Duration(period) * time.Second)

	fmt.Fprintf(w, "Time is up. Hogging for %d seconds", period)
}

func memoryHogging(w http.ResponseWriter, r *http.Request) {
	period, err := strconv.Atoi(r.FormValue("period"))
	if err != nil {
		http.Error(w, "Could not parse lasting period form value", http.StatusBadRequest)
		return
	}

	size, err := strconv.Atoi(r.FormValue("size"))
	if err != nil {
		http.Error(w, "Could not parse memory size form value", http.StatusBadRequest)
		return
	}

	//in MB
	b := make([]byte, size*1024*1024)
	b[0] = 1
	b[len(b)-1] = 1

	time.Sleep(time.Duration(period) * time.Second)
	fmt.Fprintf(w, "Time is up. Hogging memory %d MB for %d seconds", size, period)
}

func bgHandler(w http.ResponseWriter, r *http.Request) {
	htmlTemplate := `
	<!doctype html>
	<html>
	<head>
		<meta charset="utf-8">
		<title>A Blue Greem Deployment Sample</title>
		<style>
			body {
				background-color: {{ .BGColor }};
			}
			h1 {
				color: white;
			}
		</style>
	</head>
	<body>
		<h1>Hello Knative!</h1>
	</body>
	</html>		
	`
	t, err := template.New("bg").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, "Failed to create html template", http.StatusInternalServerError)
		return
	}

	type HTMLData struct {
		BGColor string
	}

	bg := os.Getenv("BGROUND")
	t.Execute(w, HTMLData{BGColor: bg})

}

func main() {
	http.HandleFunc("/", greet)
	http.HandleFunc("/cpuHog", cpuHogging)
	http.HandleFunc("/memHog", memoryHogging)
	http.HandleFunc("/timeHog", timeHogging)
	http.HandleFunc("/blue-green", bgHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listen on port:%s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
