package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Randomized User Agent
var userAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 7_0_1 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) Version/7.0 Mobile/11A470a Safari/9537.53"

// Path and Files
var pathF = "files"
var extensions = "extensions.txt"
var files = "files.txt"
var folders = "folders.txt"
var foundedFolders []string

// httpGet isminde bir fonksiyon oluşturuyoruz. Bu fonksiyon içerisine bir url parametresi alıyor. Bu fonksiyon http get sorgusu yaparak dönen response'u döndürüyor.
func httpGet(url string) *http.Response {
	
	// http.Get fonksiyonu ile url parametresi ile gelen url'e bir http get sorgusu yapıyoruz. Bu fonksiyon bize bir response ve error döndürüyor.
	//response, err := http.Get(url)

	// http get request without no follow redirect
	 client := &http.Client{
	 	CheckRedirect: func(req *http.Request, via []*http.Request) error {
	 		return http.ErrUseLastResponse
	 	},
	 }
	 response, err := client.Get(url)

	// Eğer error varsa log.Fatal ile hata mesajını ekrana yazdırıyoruz.
	if err != nil {
		log.Fatal(err)
	}

	// response'u döndürüyoruz.
	return response
}

// getStatusCode isminde bir fonksiyon oluşturuyoruz. Bu fonksiyon içerisine bir response parametresi alıyor. Bu fonksiyon response'un status kodunu döndürüyor.
func getStatusCode(resp *http.Response) string {
	
	// response'un status kodunu döndürüyoruz.
	return strconv.Itoa(resp.StatusCode)
}

func getContentType(resp *http.Response) string {
	
	var contentType = resp.Header.Get("Content-type")
	
	if contentType == "" {
		return "application/octet-stream"
	}
	
	return contentType
}

// checkSuitableContentType isimli bir fonksiyon olsun. Bu fonksiyon içerisine bir contentType parametresi alsın. Bu fonksiyon contentType parametresi ile gelen değerin text/html olup olmadığını kontrol edip true veya false döndürsün.
func checkSuitableContentType(contentType string) bool {

	// Eğer contentType içeriğinde text/html değeri içeriyorsa true döndürüyoruz.
    if strings.Contains(contentType, "text/html") {
        return false
    }

	return true
}

func checkSuitableStatusCode(statusCode string) bool {

	if statusCode == "200" || statusCode == "301" || statusCode == "302" || statusCode == "304" || statusCode == "307" || statusCode == "403" {
		return true
	} else {
		return false
	}
}

func scanFiles() {

	for _, fndPTHS := range foundedFolders {
		fmt.Println("\n************* Starting Scan Backups Files. / PATH : " + fndPTHS + " *************\n")

		fileslist, err := os.Open(pathF + "/files.txt")
		
		if err != nil {
			log.Fatal(err)
		}
		
		defer fileslist.Close()

		fileScan := bufio.NewScanner(fileslist)
		
		for fileScan.Scan() {

			extensionss, err := os.Open(pathF + "/extensions.txt")
			
			if err != nil {
				log.Fatal(err)
			}
			
			defer extensionss.Close()

			scanner := bufio.NewScanner(extensionss)
			
			for scanner.Scan() {

				var urlE = fndPTHS + fileScan.Text() + scanner.Text()
				
				var httpResponse = httpGet(urlE)

				var lastStatusCode = getStatusCode(httpResponse)
				var lastContentType = getContentType(httpResponse)

				var chckDrm = urlE + " | Response Code : " + lastStatusCode

				var suitableStatusCode = checkSuitableStatusCode(lastStatusCode)
				var suitableContentType = checkSuitableContentType(lastContentType)
				
				if suitableStatusCode == true && suitableContentType == true {
					fmt.Printf("\033[2K\r%s\n", "* Founded File Path : " + chckDrm)
					foundedFolders = append(foundedFolders, urlE)
				} else {
					fmt.Printf("\033[2K\r%s", "Checking File Path : " + chckDrm)
				}

			}

		}

		fmt.Printf("\033[2K\r%s", "")
		fmt.Println("\n************* Scan Ended. / PATH : " + fndPTHS + " *************\n")
	}

}

func scanPath(filename string, hostname string) string {

	file, err := os.Open(pathF + "/" + filename)
	
	if err != nil {
		log.Fatal(err)
	}
	
	defer file.Close()

	scanner := bufio.NewScanner(file)

	fmt.Println("\n************* Starting Scan Backups Dir PATHS *************\n")

	var lastStatusCode = ""

	for scanner.Scan() {
		var urlE = hostname + "/" + scanner.Text() + "/"

		var httpResponse = httpGet(urlE)

		lastStatusCode = getStatusCode(httpResponse)
		var chckDrm = "" + urlE + " | Response Code : " + lastStatusCode

		if lastStatusCode == "200" || lastStatusCode == "301" || lastStatusCode == "302" || lastStatusCode == "304" || lastStatusCode == "307" || lastStatusCode == "403" {
		// if lastStatusCode == "200" || lastStatusCode == "301" || lastStatusCode == "302" || lastStatusCode == "304" || lastStatusCode == "307" || lastStatusCode == "403" {
			fmt.Printf("\033[2K\r%s\n", "* Founded Dir Path : "+chckDrm)
			foundedFolders = append(foundedFolders, urlE)
		} else {
			fmt.Printf("\033[2K\r%s", "Checking Dir Path : " + chckDrm)
		}

	}

	fmt.Printf("\033[2K\r%s", "\nPath Scaning Ended.\n")
	scanFiles()
	
	return lastStatusCode
}

func main() {

	hostname := flag.String("hostname", "", "Please input hostname")
	flag.Parse()

	fmt.Println(`
	____ ____ ____ ____ ____ ____ ____ ____ ____ ____ 
	||O |||h |||M |||y |||B |||a |||c |||k |||U |||p ||
	||__|||__|||__|||__|||__|||__|||__|||__|||__|||__||
	|/__\|/__\|/__\|/__\|/__\|/__\|/__\|/__\|/__\|/__\|
						
	Backup Directories & Backup Files Scanner.
	Github : github.com/muminkoykiran
	Host : ` + *hostname + `
	`)

	scanPath(folders, *hostname)

}
