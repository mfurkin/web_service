package main 

import (
	"io"
	"os"
    "net/http"
    "fmt"
    "fileChecker"
)
var dirname ="/home/mikhailf/web_work/workdir1"
var fc = fileChecker.NewFileChecker(&dirname)
// Default Request Handler
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	
	method := r.Method
//	fileChecker := fileChecker.NewFileChecker(&dirname)
	fmt.Println("defaultHandler method="+method)
	switch method {
		case "GET" : 
			len1 := len(r.URL.Path[1:]) 
			fmt.Printf("GET method len1=%d\n", len1)
			if len1 == 0 {		
				files, err := fc.GetFiles()
				if err != nil {
					fmt.Println("Some error during file list getting: "+err.Error())
					w.WriteHeader(500)
					fmt.Fprintf(w, "Произошла внутренняя ошибка\n")
				} else {
					w.WriteHeader(200)
					fmt.Printf("File list have been got successfully len=%d\n",len(files))
					for _, file := range files {
						fmt.Fprintf(w,"<a href=\"localhost:8080/%s\">%s</a>\n", file,file)
					}
				}		
			} else {
						fname  := r.URL.Path[1:]
						err:= fc.GetFile(&fname,w)
						if err != nil {
							if err == io.EOF {
								w.WriteHeader(404)
								fmt.Fprintf(w, "File %s have not been found", fname)
							} else {
								w.WriteHeader(500)
								fmt.Fprintf(w, "File %s could not be got", fname)
							}
						}/* else {
							w.WriteHeader(200)
						}	*/		
				// ТODO кинуть конкретный файл
			}
			break
		case "POST":
			len1 := len(r.URL.Path[1:])
			if len1 == 0 {
				w.WriteHeader(400)
				fmt.Fprintf(w, "Incorrect request -filename have not been specified\n")
			} else {
				fname := r.URL.Path[1:]
				if err := fc.CreateFile(&fname,r.Body); err != nil {
					w.WriteHeader(500)
					fmt.Fprintf(w,"File %s could not be created\n",fname)
				} else {
					w.WriteHeader(201)
					fmt.Fprintf(w,"File %s have been creted successfully\n",fname)
				}
			}
			break
		case "DELETE":
			len1 := len(r.URL.Path[1:])
			if len1 == 0 {
				w.WriteHeader(400)
				fmt.Fprintf(w, "Incorrect request - file name not \n")
			} else {
				fname := r.URL.Path[1:]
				err := fc.RemoveFile(&fname)
				if  err != nil {
					if os.IsNotExist(err) {
						w.WriteHeader(404)
						fmt.Fprintf(w, "File %s have not been found\n",fname)
					} else {
						w.WriteHeader(500)
						fmt.Fprintf(w, "FIle %s could not be deleted\n",fname)
					}
				} else {
					w.WriteHeader(200)
					fmt.Fprintf(w, "File %s have been deleted successfully\n", fname)
				}
			}
			break	
	}
}

func main() {
    http.HandleFunc("/", defaultHandler)
    http.ListenAndServe(":8080", nil)
}

