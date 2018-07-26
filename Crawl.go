package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"os"
	"regexp"
	"net/url"
	"runtime"
)

var urlChanne = make(chan string, 200)// url channe distribute task
var atagRegExp = regexp.MustCompile(`<a[^>]+[(href)|(HREF)]\s*\t*\n*=\s*\t*\n*[(".+")|('.+')][^>]*>[^<]*</a>`)

type Sconfig struct {
	XMLName xml.Name `xml:"note"`   //point to outer label
	SeedUrl string `xml:"url"`   // read url and save url into SeedUrl
}
var Seedurl string    //SeedUrl
func   getSeedUrl(){       //get seedUrl from default.xml
	file,err:=os.Open("Default.xml")
	if err!=nil{
		fmt.Printf("error:%v",err)
		return
	}
	defer file.Close()
	data, err :=ioutil.ReadAll(file)
	if err !=nil{
		fmt.Printf("error:%v",err)
		return
	}
	v :=Sconfig{}
	err=xml.Unmarshal(data,&v)
	if err !=nil{
		fmt.Printf("error:%v",err)
		return
	}
	Seedurl=v.SeedUrl
}
func main() {
	getSeedUrl()   //get SeedUrl
//	proxy_addr :="http://183.221.250.137:80/"

//	html:=fetchHtml(&url,&proxy_addr)
//	fmt.Println(html)

go Spy(Seedurl)
	for url := range urlChanne {
		fmt.Println("rountines num=", runtime.NumGoroutine(), "chan len=", len(urlChanne))
	go Spy(url)
	}
	fmt.Println("a")
}
func fetchHtml(url,proxy_addr *string)(html string){
	transport :=getTransportFieldURL(proxy_addr)
	client :=&http.Client{Transport:transport}
	req,err :=http.NewRequest("GET",*url,nil)

	if err!=nil{
		log.Fatal(err.Error())
	}
	fmt.Println(req)
	resq,err:=client.Do(req)
	fmt.Println("sdsd")
	fmt.Println((ioutil.ReadAll(req.Body)))
	if err!=nil{
		log.Fatal(err.Error())
	}



	fmt.Println(resq.StatusCode)
	if resq.StatusCode==200{
		fmt.Println(8)
		robots,err:=ioutil.ReadAll(req.Body)
		resq.Body.Close()
		if err!=nil{
			log.Fatal(err.Error())
		}
		html=string(robots)
		fmt.Println("html:="+html)
	}else{
		html=""
	}
	return
}
func getTransportFieldURL(proxy_addr *string)(transport *http.Transport){
	url_i :=url.URL{}
	url_proxy ,_ :=url_i.Parse(*proxy_addr)
	transport =&http.Transport{Proxy: http.ProxyURL(url_proxy)}
	return
}
func getTransport()(transport *http.Transport){
	transport =&http.Transport{Proxy:http.ProxyFromEnvironment}
	return
}

func Spy(url string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("[E]", r)
		}
	}()

	req, _ := http.NewRequest("GET", url, nil)
	client := http.DefaultClient
	res, e := client.Do(req)



	if e != nil {
		fmt.Printf("get请求%s返回错误：%s", url, e)
		return
	}

	if res.StatusCode == 200 {
		body := res.Body
		defer body.Close()
		bodyByte, _ := ioutil.ReadAll(body)
		resStr := string(bodyByte)
		fmt.Println()
		atag := atagRegExp.FindAllString(resStr, -1)

		for _, a := range atag {
			href, _ := GetHref(a)
			fmt.Println("^", href)

			urlChanne <- href
		}
	}
}


var r = rand.New(rand.NewSource(time.Now().UnixNano()))



func GetHref(atag string) (href, content string) {
	inputReader := strings.NewReader(atag)
	decoder := xml.NewDecoder(inputReader)
	fmt.Println(decoder)
	for t, err := decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.StartElement:
			for _, attr := range token.Attr {
				attrName := attr.Name.Local
				attrValue := attr.Value
				if strings.EqualFold(attrName, "href") || strings.EqualFold(attrName, "HREF") {
					href = attrValue
				}
			}
		case xml.EndElement:
		case xml.CharData:
			content = string([]byte(token))
		default:
			href = ""
			content = ""
		}
	}
	return href, content
}
