package spider

import (
	"sync"
	"net/http"
	"time"
	"github.com/Pallinder/go-randomdata"
	"github.com/PuerkitoBio/goquery"
	"errors"
	"strings"
	"fmt"
	"net/http/cookiejar"
	"log"
)

const (
	NoLogin = iota
	AlreadyLogin
)


type TF31Proxy struct {
	sync.RWMutex
	loginState int
	httpClient *http.Client
	xsrf string
	cookies []*http.Cookie
} 

func NewF31Proxy() *TF31Proxy {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	return &TF31Proxy{
		loginState: NoLogin,
		httpClient:&http.Client{
			Jar:jar,
			Transport: &http.Transport{
				IdleConnTimeout: 10 * time.Second,
			},
		},
	}
}

func (obj *TF31Proxy) Pull() {
	obj.Lock()
	defer obj.Unlock()

	if obj.loginState == NoLogin {
		err := obj.login()
		if err != nil {
			Logger.Printf("login f31 failed %v", err)
			return
		}
		obj.loginState = AlreadyLogin
	}

	request, err := http.NewRequest("GET", GTSpiderConfig.TargetURL, nil)
	if err != nil {
		Logger.Print(err)
		return
	}

	request.Header.Set("Referer", "http://31f.cn/")
	request.Header.Set("User-Agent", randomdata.UserAgentString())
	response, err := obj.httpClient.Do(request)
	if err != nil {
		Logger.Print(err)
		return
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		Logger.Print(err)
		return
	}

	document.Find("table").Eq(0).Find("tr").Each(func(i int, selection *goquery.Selection) {
		if i == 0 {
			return
		}

		matchSelection := selection.Find("td")
		if matchSelection.Eq(0).Length() == 0 || matchSelection.Eq(1).Length() == 0  || matchSelection.Eq(2).Length() == 0 {
			Logger.Print("not found match selection")
			return
		}

		address := matchSelection.Eq(1).Text()
		port := matchSelection.Eq(2).Text()
		country := matchSelection.Eq(3).Text()
		category := matchSelection.Eq(4).Text()

		SpiderGlobalVars.TQueue.Push(TQueueItem{IsSock5, &TProxyItem{category, country, address, port, time.Now().Unix()}})
	})
}

func (obj *TF31Proxy) getRsrf() error  {
	request, err := http.NewRequest("GET", "http://31f.cn/login", nil)
	if err != nil {
		Logger.Print(err)
		return err
	}

	request.Header.Set("Referer", "http://31f.cn/")
	request.Header.Set("user-agent", randomdata.UserAgentString())
	response, err := obj.httpClient.Do(request)
	if err != nil {
		Logger.Print(err)
		return err
	}

	defer response.Body.Close()
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		Logger.Print(err)
		return err
	}

	value, exists := document.Find("input[name=_xsrf]").Eq(0).Attr("value")
	if !exists {
		return errors.New("found login _xsrf failed")
	}

	obj.xsrf = value
	Logger.Printf("found login _xsrf %s", value)

	return nil
}

func (obj *TF31Proxy) login() error  {
	err := obj.getRsrf()
	if  err != nil {
		return err
	}

	body := fmt.Sprintf("_xsrf=%s&next=/&username=%s&password=%s", obj.xsrf, GTSpiderConfig.F31User.Name, GTSpiderConfig.F31User.Password)
	request, err := http.NewRequest("POST", "http://31f.cn/login", strings.NewReader(body))
	if err != nil {
		Logger.Print(err)
		return err
	}

	request.Header.Set("Referer", "http://31f.cn/login")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Set("X-Requested-With", "XMLHttpRequest")
	request.Header.Set("User-Agent", randomdata.UserAgentString())

	response, err := obj.httpClient.Do(request)
	if err != nil {
		Logger.Print(err)
		return err
	}

	Logger.Printf("get header %v", response.Header)
	Logger.Printf("get cookie %v", response.Cookies())

	return nil
}