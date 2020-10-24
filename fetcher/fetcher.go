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
	cookie1 := "sid=74026cbc-17c9-4a13-965e-03904c504003; ec=L8PwNAXN-1603458306973-e29e32137bb33-1975754967; Hm_lvt_2c8ad67df9e787ad29dbd54ee608f5d2=1603458312; __channelId=900122%2C0; _efmdata=NJBqJAOKVsOQoSMdxBQbIqOuP77KBjHyrjeW4Tvnr9fkNdMBXHEGbm3haBzsEqWgzeGXzWBOqG2cdLZSaJbqxFYJSjIHjHIREyKrp%2Fx2NLs%3D; _exid=oR3Ug6vruo8bzqXeIZ9fjW6D6SOp4c0OQyOhY%2F4sxNzHsVTNyvk8IkihU90ZE12zRCOgQ%2BNdzuFL6qfFzQvQAw%3D%3D; Hm_lpvt_2c8ad67df9e787ad29dbd54ee608f5d2=1603464649"
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
