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

const CMD_LINE_ARG_NAME_IP string = "ip"
const CMD_LINE_ARG_NAME_PORT string = "port"
const CMD_LINE_ARG_NAME_SERVING_DIRECTORY string = "dir"

const CMD_LINE_ARG_NAME_PREFIX string = "-"

type ComandLineArg struct {
    Name string
    Value  string
}

func parseComandLineArgs(args []ComandLineArg, argNamePrefix string) []ComandLineArg {
	for argIndex,arg := range args {
		argFinded := false
		for osArgIndex,osArg := range os.Args {
			if(osArgIndex == 0) { continue }
			if (argFinded){
				args[argIndex].Value = osArg
				break
			}
			if(argNamePrefix+arg.Name == osArg){
				argFinded = true
			}
		}
	}
	return args
}

func findComandLineArg(args []ComandLineArg, name string, defValue string) string {
	resValue := defValue
	for _,arg := range args {
		if(arg.Name == name){
			resValue = arg.Value
			break
		}
	}
	return resValue
}

func printSupportedComandLineArgs(args []ComandLineArg, argNamePrefix string, header string, footer string) {
	if(header != ""){
		_,err := fmt.Println(header)
		if(err != nil){
			log.Println(fmt.Sprintf("Error printing to standard otput %s: ", err))
		}
	}
	for _,arg := range args {
		_,err := fmt.Println(fmt.Sprintf("%s%s: %s",argNamePrefix,arg.Name,arg.Value))
		if(err != nil){
			log.Println(fmt.Sprintf("Error printing to standard otput %s: ", err))
		}
	
		//log.Println(fmt.Sprintf("%s%s: %s",argNamePrefix,arg.Name,arg.Value))
	}
	if(footer != ""){
		_,err := fmt.Println(footer)
		if(err != nil){
			log.Println(fmt.Sprintf("Error printing to standard otput %s: ", err))
		}
		//log.Println(footer)
	}
}

/*func absPath(relPath string) string{
	
		wd, err := filepath.Abs(relPath)
		if err != nil {
			log.Fatal(fmt.Sprintf("filepath.Abs failed: %s", err))
		}
		return wd
	
}*/

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

/*func getPortIp(defPort string, defIp string) (string, string) {
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
}*/


func readParameters() (string, string, string) {
	args := parseComandLineArgs([]ComandLineArg{
			ComandLineArg{CMD_LINE_ARG_NAME_SERVING_DIRECTORY,DEFAULT_SERVING_DIRECTORY},
			ComandLineArg{CMD_LINE_ARG_NAME_PORT,DEFAULT_PORT},
			ComandLineArg{CMD_LINE_ARG_NAME_IP,DEFAULT_IP}},
		CMD_LINE_ARG_NAME_PREFIX)
	
	servingDir := absAppPath(findComandLineArg(args,"dir",DEFAULT_SERVING_DIRECTORY))
	port := findComandLineArg(args,"port",DEFAULT_PORT)
	ip := findComandLineArg(args,"ip",DEFAULT_PORT)
	
	return servingDir, port, ip
} 

func printParametersHelp(){
	defIpForPrint := DEFAULT_IP
	if(defIpForPrint == ""){
		defIpForPrint = "blank (use all available IPs)"
	}
	printSupportedComandLineArgs(
		[]ComandLineArg{
			ComandLineArg{
				CMD_LINE_ARG_NAME_SERVING_DIRECTORY, 
				fmt.Sprintf("Web server root directory. Default is \"%s\"", DEFAULT_SERVING_DIRECTORY)},
			ComandLineArg{
				CMD_LINE_ARG_NAME_PORT, 
				fmt.Sprintf("Listening port. Default is %s", DEFAULT_PORT )},
			ComandLineArg{
				CMD_LINE_ARG_NAME_IP, 
				fmt.Sprintf("Listening IP. Default is %s", defIpForPrint )}},
		CMD_LINE_ARG_NAME_PREFIX,
		"Supported command line arguments:","**********************************")
}

func logStartCondition(servingDir, port, ip string){
	log.Println(fmt.Sprintf("Application path: %s", os.Args[0]))
	log.Println(fmt.Sprintf("Web server root directory: %s", servingDir))
	if(ip != ""){
		log.Println(fmt.Sprintf("Listening IP address \"%s\" on port \"%s\"",ip, port))
	} else {
		log.Println(fmt.Sprintf("Listening all IP addresses on port \"%s\"", port))
	}
	log.Println("Starting server...")
}

func main() {
	printParametersHelp()
	servingDir, port, ip := readParameters()
	fs := http.FileServer(http.Dir(servingDir))
	logStartCondition(servingDir, port, ip)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s",ip,port), fs))
}