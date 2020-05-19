package main

import (
"net/http"
"fmt"
"log"
"github.com/julienschmidt/httprouter"
"path/filepath"
"io/ioutil"
"strings"
"os"
"encoding/json"
"strconv"
)

var ChapterData = map[string]map[string]interface{}{}

/*
 * Cache-Control: no-cache, no-store, must-revalidate
 * Pragma: no-cache
 * Expires: 0
*/
func setNoCacheHeader(w http.ResponseWriter) {
	w.Header().Set("Cache-Control","no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma","no-cache")
	w.Header().Set("Expires","0")
}

func parseChapter(path string, info os.FileInfo, err error) error {
	var retErr error
	if info.IsDir() {
		if info.Name() != "static" && info.Name() != "build"{
			chapter := 1
			prevChapterPath := ""
			for true {
				chapterString := strconv.Itoa(chapter)
				filePath := filepath.Join(path, chapterString+".html")
				_, err := os.Stat(filePath)
				if err != nil {
					retErr = err
					if chapter > 1 {
						ChapterData[prevChapterPath]["hasNext"] = false
					}
					fmt.Println(err)
					break
				} else {
					fileData, err := ioutil.ReadFile(filePath)

					if err != nil {
						retErr = err
						if chapter > 1 {
							ChapterData[prevChapterPath]["hasNext"] = false
						}
						fmt.Println(err)
						break
					}

					ChapterData[filePath]= map[string]interface{}{
						"data": string(fileData),
						"hasPrev": (chapter != 1),
						"hasNext": true,
					}
					prevChapterPath = filePath
				}
				chapter = chapter + 1
			}
		}
	}
	return retErr
}

func parseChapters() {
	basePath := os.Getenv("BOOK_BASE_PATH")
	filepath.Walk(basePath,parseChapter)
}

func Chapter(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	basePath := os.Getenv("BOOK_BASE_PATH")
	name :=  ps.ByName("name")
	chapter :=  ps.ByName("chapter")

	filePath := filepath.Join(basePath, name, chapter+".html")
	if strings.Contains(filePath,"..") {
		http.Error(w, "Invalid path", http.StatusInternalServerError)
		return
	}
	fmt.Println(filePath)
	if data,ok := ChapterData[filePath]; ok {
		jsonData,err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
    }
		setNoCacheHeader(w)
		setContentHeader(w, "data.json")
		w.Write(jsonData)
	} else {
		http.Error(w, "Invalid Path", http.StatusInternalServerError)
		return
	}
}

func setContentHeader(w http.ResponseWriter, path string) {
	if strings.HasSuffix(path,".jpg") {
		w.Header().Set("Content-Type", "image/jpeg")
	} else if strings.HasSuffix(path,".jpeg") {
		w.Header().Set("Content-Type", "image/jpeg")
	} else if strings.HasSuffix(path, ".png") {
		w.Header().Set("Content-Type", "image/png")
	} else if strings.HasSuffix(path, ".html") {
		w.Header().Set("Content-Type", "text/html")
		setNoCacheHeader(w)
	} else if strings.HasSuffix(path, ".js") {
		w.Header().Set("Content-Type", "text/javascript")
		setNoCacheHeader(w)
	} else if strings.HasSuffix(path, ".css") {
		w.Header().Set("Content-Type", "text/css")
	} else if strings.HasSuffix(path, ".json") {
		w.Header().Set("Content-Type", "application/json")
	}
}

func StaticPath(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	basePath := os.Getenv("BOOK_BASE_PATH")
	path :=  r.URL.Path

	if strings.Contains(path,"..") {
		http.Error(w, "Invalid path", http.StatusInternalServerError)
		return
	}
	filePath := filepath.Join(basePath, path)
	fmt.Println(filePath)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if fileInfo.IsDir() {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileData, err := ioutil.ReadFile(filePath)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	setContentHeader(w, filePath)
	w.Write(fileData)
}


func Static(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	basePath := os.Getenv("BOOK_BASE_PATH")
	path :=  ps.ByName("path")

	if strings.Contains(path,"..") {
		http.Error(w, "Invalid path", http.StatusInternalServerError)
		return
	}
	filePath := filepath.Join(basePath, "/static/", path)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if fileInfo.IsDir() {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileData, err := ioutil.ReadFile(filePath)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	setContentHeader(w, filePath)
	w.Write(fileData)
}
func StoryStatic(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	basePath := os.Getenv("BOOK_BASE_PATH")
	story :=  ps.ByName("name")
	path :=  ps.ByName("path")

	filePath := filepath.Join(basePath, story, path)
	if strings.Contains(filePath,"..") {
		http.Error(w, "Invalid path", http.StatusInternalServerError)
		return
	}
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if fileInfo.IsDir() {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileData, err := ioutil.ReadFile(filePath)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	setContentHeader(w, filePath)
	w.Write(fileData)
}


func Start(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	basePath := os.Getenv("BOOK_BASE_PATH")
	filePath:= filepath.Join(basePath, "index.html")
	_, err := os.Stat(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileData, err := ioutil.ReadFile(filePath)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(fileData)
}

func main() {
	port := os.Getenv("PORT")
	parseChapters()
	router := httprouter.New()
	router.GET("/read/:name/", Start)
	router.GET("/read/:name/chapter/:chapter/", Chapter)
	router.GET("/read/:name/static/*path", StoryStatic)
	router.GET("/static/*path", Static)
	router.GET("/service-worker.js", StaticPath)
	router.GET("/index.html", StaticPath)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
