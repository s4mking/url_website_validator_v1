package main

import (
	"fmt"
	"github.com/bradfitz/slice"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type linkWeb struct{
	Link string
	Code string
}
var errorArray []linkWeb
var errorCodeHttp string

func parseSite(sitename string,racine string)(err error){
	/*Récupération de tous les a href puis traitement afin de n'avoir plus que link pareil pour le js
	obligé de prendre les scripts sinon images aussi */
		var baseSite = racine
		site,_ := http.Get(sitename)
		html,_ := ioutil.ReadAll(site.Body)
		errorCodeHttp=site.Status
		all_js :=regexp.MustCompile(`src="[^"]*"+`)
		all_a := regexp.MustCompile(`(.)*href="[^"]*"+`)
		var resultsJs = all_js.FindAllStringSubmatch(string(html),-1)
		//fmt.Println(resultsJs)
		var resultUrls = all_a.FindAllStringSubmatch(string(html),-1)

		for _,linkJs := range resultsJs{

			re := regexp.MustCompile(`src="[^"]*"`)
			var execReg = re.FindAllStringSubmatch(linkJs[0],-1)
			resultrmvbefore := strings.TrimPrefix(execReg[0][0], "src=\"")
			resultrmvafter := strings.TrimSuffix(resultrmvbefore, "\"")
			if resultrmvafter==""{
				resultrmvafter=racine
			}
			firstChar := resultrmvafter[:1]
			twoChars := resultrmvafter[:2]
			if firstChar == "#" {
				continue
			}
			if twoChars == "//"{
			break
			}
			if firstChar == "/" {
				baseSite = strings.TrimSuffix(baseSite, "/")
				resultrmvafter = baseSite+resultrmvafter
			}
			//Test si on est bien sur le site original
			siteName := regexp.MustCompile(`(.)*`+baseSite+`(.)*`)
			siteName_presence := siteName.MatchString(resultrmvafter)
			//Insertion code et erreur dans le tableau global et test présence avant
			var dejavue=0
			for _,item_url := range errorArray{
				if resultrmvafter == item_url.Link{
					dejavue=1
				}
			}
			if(dejavue==0 && siteName_presence == false){
				var actualState = linkWeb{resultrmvafter,"-> "+errorCodeHttp+"\r\n"}
				errorArray = append(errorArray,actualState)
				continue
			}else{
				if dejavue==0 {
					var actualState = linkWeb{resultrmvafter,"-> "+errorCodeHttp+"\r\n"}
					errorArray = append(errorArray,actualState)
					parseSite(string(resultrmvafter),baseSite)
				}
			}
		}
		for _,link := range resultUrls{
			re := regexp.MustCompile(`href="[^"]*"`)
			var testsam = re.FindAllStringSubmatch(link[0],-1)
			resultrmvbefore := strings.TrimPrefix(testsam[0][0], "href=\"")
			resultrmvafter := strings.TrimSuffix(resultrmvbefore, "\"")
			if resultrmvafter==""{
				resultrmvafter=racine
			}
			firstChar := resultrmvafter[:1]
			if firstChar == "#" {
				continue
			}
			if firstChar == "/" {
				baseSite = strings.TrimSuffix(baseSite, "/")
				resultrmvafter = baseSite+resultrmvafter
			}
			//Test si on est bien sur le site original
			siteName := regexp.MustCompile(`(.)*`+baseSite+`(.)*`)
			siteName_presence := siteName.MatchString(resultrmvafter)
			//Insertion code et erreur dans le tableau global et test présence avant
			var dejavue=0
			for _,item_url := range errorArray{
				if resultrmvafter == item_url.Link{
					dejavue=1
				}
			}
			if(dejavue==0 && siteName_presence == false){
				var actualState = linkWeb{resultrmvafter,"-> "+errorCodeHttp+"\r\n"}
				errorArray = append(errorArray,actualState)
				continue
			}else{
				if dejavue==0 {
					var actualState = linkWeb{resultrmvafter,"-> "+errorCodeHttp+"\r\n"}
					errorArray = append(errorArray,actualState)
					parseSite(string(resultrmvafter),baseSite)
				}
			}
		}
	return
}

func main() {
	//Appel function et traitement des données puis insertiond dans fichier texte
	parseSite(string(os.Args[1]),string(os.Args[1]))
	regexFileName := regexp.MustCompile("/(.)*")
	file_name := (regexFileName.FindString(string(os.Args[1])))
	file_name = strings.TrimPrefix(file_name, "//")
	newfile,_:= os.Create((file_name+".txt"))
	fmt.Println(len(errorArray))
	slice.Sort(errorArray[:], func(i, j int) bool {
		return errorArray[i].Link < errorArray[j].Link
	})
	var array_url [][]byte
		for _,item_array := range errorArray {
			var linkByte = []byte(item_array.Link)
			var urlByte = []byte(item_array.Code)
			array_url = append(array_url,linkByte)
			array_url = append(array_url,urlByte)
		}

		for _,row_array := range array_url {
			newfile.Write(row_array)
		}
	}