package fetcher

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var rateLimiter = time.Tick(10 * time.Millisecond) //设置成100ms一个tick，也就意味着每秒10个request。保持这个速度去爬，可能就不会被对方的网站限制

func Fetch(url string) ([]byte, error) {
	//resp, err := http.Get(url)
	//if err != nil {
	//	return nil, err
	//}
	//defer resp.Body.Close()

	<-rateLimiter //防止爬取过快，被网站给限制了
	client := &http.Client{}
	newUrl := strings.Replace(url, "http://", "https://", 1)
	req, err := http.NewRequest("GET", newUrl, nil)
	if err != nil {
		panic(err)
		return nil, err
	}
	//https://blog.csdn.net/kemuxiaozi000/article/details/107166072/ 模拟浏览器访问，避免403和202问题
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36")
	cookie1 := "sid=8cdccae7-2c7c-49d6-b61c-a3ee8dd6d0ae; ec=NrYIwnN6-1603065539092-a1e797ce64a69-462919474; Hm_lvt_2c8ad67df9e787ad29dbd54ee608f5d2=1603065549; FSSBBIl1UgzbN7NO=58iMcUCzMnS43eM_Fychxw47UTTDGlYpTKhHYABdJgsGKKlPESUDmP3NqeZdoddupV3H69MY_TJNKQ5S7lBIMFA; Hm_lpvt_2c8ad67df9e787ad29dbd54ee608f5d2=1603510594; _efmdata=qQ1RQJKdrOvMWcN6t7VPDn9G19Hs8rzfZvHAhl5lz7BLd4Mpx0SxLGf5v7k7jr1yEibXZXWaI%2BpfLOkONBKt1J0fa3TYsadVzIQNQc7AonE%3D; _exid=nCxk29IEvhVK2uHJVG5zgg6VPCSL%2Fjsg90CUdNwKIy0I8c%2BbMnPuq%2Bjjh7O3MfUSNs7IYDitc4KRl%2Fy3rAJUBA%3D%3D; FSSBBIl1UgzbN7NP=5UyifA25sLVlqqqmTKWAAZaV.L5.sn3JH3N7bciIZ5bTR4KSIUqhtCy1nzeGAfjKSwffXWI_mu5ZKFzQK6MeYObSztMK4C1AtfeXEwynmxvYucI.W3Lm8DWj3NK3Yz.DIvTowXkjRim_9uDt_yXzzwYED7bWxCyP2D0l0G4vXP.6OF2dqaLJ457fT9s7Oju4AS4rzsbupZ73cwzaQ7XII1RzCyK5YyG4rxHk31XfnNKQr.yA6QXFIc8a9uOiwGesm6hkdUv9UnH3SYEPK1_waM3"
	req.Header.Add("cookie", cookie1)

	resp,err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)
	}

	//处理网页编码的问题
	bodyReader := bufio.NewReader(resp.Body)
	e := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())

	return ioutil.ReadAll(utf8Reader)
}

func determineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Printf("Fetcher error: %v", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")

	return e
}
