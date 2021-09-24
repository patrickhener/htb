package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"
)

//go:embed static/public
var web embed.FS

//go:embed static/public/index.html
var index []byte

type Pet struct {
	Name            string `json:"name"`
	Species         string `json:"species"`
	Characteristics string `json:"characteristics"`
}

var (
	Pets []Pet = []Pet{
		{Name: "Cookie", Species: "cat", Characteristics: loadCharacter("cat")},
		{Name: "Mia", Species: "cat", Characteristics: loadCharacter("cat")},
		{Name: "Chuck", Species: "dog", Characteristics: loadCharacter("dog")},
		{Name: "Balu", Species: "dog", Characteristics: loadCharacter("dog")},
		{Name: "Georg", Species: "gopher", Characteristics: loadCharacter("gopher")},
		{Name: "Gustav", Species: "giraffe", Characteristics: loadCharacter("giraffe")},
		{Name: "Rudi", Species: "redkite", Characteristics: loadCharacter("redkite")},
		{Name: "Bruno", Species: "bluewhale", Characteristics: loadCharacter("bluewhale")},
	}
)

func loadCharacter(species string) string {
	cmd := exec.Command("sh", "-c", "cat characteristics/"+species)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return err.Error()
	}
	return string(stdoutStderr)
}

func getPets(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Pets)
}

func addPet(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var addPet Pet
	err := json.Unmarshal(reqBody, &addPet)
	if err != nil {
		e := fmt.Sprintf("There has been an error: %+v", err)
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	addPet.Characteristics = loadCharacter(addPet.Species)
	Pets = append(Pets, addPet)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Pet was added successfully")
}

func handleRequest() {
	build, err := fs.Sub(web, "static/public/build")
	if err != nil {
		panic(err)
	}

	css, err := fs.Sub(web, "static/public/css")
	if err != nil {
		panic(err)
	}

	webfonts, err := fs.Sub(web, "static/public/webfonts")
	if err != nil {
		panic(err)
	}

	spaHandler := http.HandlerFunc(spaHandlerFunc)
	// Single page application handler
	http.Handle("/", headerMiddleware(spaHandler))

	// All static folder handler
	http.Handle("/build/", headerMiddleware(http.StripPrefix("/build", http.FileServer(http.FS(build)))))
	http.Handle("/css/", headerMiddleware(http.StripPrefix("/css", http.FileServer(http.FS(css)))))
	http.Handle("/webfonts/", headerMiddleware(http.StripPrefix("/webfonts", http.FileServer(http.FS(webfonts)))))
	http.Handle("/.git/", headerMiddleware(http.StripPrefix("/.git", http.FileServer(http.Dir(".git")))))

	// API routes
	apiHandler := http.HandlerFunc(petHandler)
	http.Handle("/api/pet", headerMiddleware(apiHandler))
	log.Fatal(http.ListenAndServe("127.0.0.1:5000", nil))
}

func spaHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write(index)
}

func petHandler(w http.ResponseWriter, r *http.Request) {
	// Dispatch by method
	if r.Method == http.MethodPost {
		addPet(w, r)
	} else if r.Method == http.MethodGet {
		getPets(w, r)

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	// TODO: Add Update and Delete
}

func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Server", "My genious go pet server")
		next.ServeHTTP(w, r)
	})
}

func main() {
	resetTicker := time.NewTicker(5 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-resetTicker.C:
				// Reset Pets to prestaged ones
				Pets = []Pet{
					{Name: "Cookie", Species: "cat", Characteristics: loadCharacter("cat")},
					{Name: "Mia", Species: "cat", Characteristics: loadCharacter("cat")},
					{Name: "Chuck", Species: "dog", Characteristics: loadCharacter("dog")},
					{Name: "Balu", Species: "dog", Characteristics: loadCharacter("dog")},
					{Name: "Georg", Species: "gopher", Characteristics: loadCharacter("gopher")},
					{Name: "Gustav", Species: "giraffe", Characteristics: loadCharacter("giraffe")},
					{Name: "Rudi", Species: "redkite", Characteristics: loadCharacter("redkite")},
					{Name: "Bruno", Species: "bluewhale", Characteristics: loadCharacter("bluewhale")},
				}

			}
		}
	}()

	handleRequest()

	time.Sleep(500 * time.Millisecond)
	resetTicker.Stop()
	done <- true
}
