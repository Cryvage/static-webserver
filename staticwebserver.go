package main

import (
	"log"
	"fmt"
	"os"
	"net/http"
	"path/filepath"
)

const DEFAULT_IP string = ""
const DEFAULT_PORT string = "80"
const DEFAULT_SERVING_DIRECTORY string = "www"

func absPath(relPath string) string{
	
		wd, err := filepath.Abs(relPath)
		if err != nil {
			log.Fatal(fmt.Sprintf("filepath.Abs failed: %s", err))
		}
		return wd
	
}

func absAppPath(relPath string) string{
	if !filepath.IsAbs(relPath) {
		ad, an := filepath.Split(os.Args[0]);
		if(len(ad) > 0){
			return filepath.Clean(filepath.Join(ad, relPath))
		}
		
		log.Fatal(fmt.Sprintf("filepath.Split(os.Args[0]) return an empty directory. Output was: %s", an))
		
		return relPath
	}
	return relPath
}

func getPortIp(defPort string, defIp string) (string, string) {
	port := defPort
	ip := defIp
	if(len(os.Args) > 1){
		port = os.Args[1]
		if(len(os.Args) > 2){
			ip = os.Args[2]
		}
	}
	return port, ip
}

func getServingDir(defServingDir string) string {
	servingDir := defServingDir
	if(len(os.Args) > 3){
		servingDir = os.Args[3]
	}
	return servingDir
}

func main() {
	servingDir := absAppPath(getServingDir(DEFAULT_SERVING_DIRECTORY))
	
	fs := http.FileServer(http.Dir(servingDir))
	log.Println(fmt.Sprintf("App contains in folder: %s", os.Args[0]))
	log.Println(fmt.Sprintf("Serving static files in folder: %s", servingDir))

	port, ip := getPortIp(DEFAULT_PORT, DEFAULT_IP)
	log.Println(fmt.Sprintf("Listening port: %s",port))
	if(ip != ""){
		log.Println(fmt.Sprintf("Listening ip address: %s",ip))
	}
	
	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s",ip,port), fs))
}