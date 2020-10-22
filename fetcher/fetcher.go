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
	cookie1 := "FSSBBIl1UgzbN7NO=5Iz4GkebY3qhBBgOsu3TFm7uOlt5hhh_ZO0ajkEqNAA08RMyrgHRzlAe4ife4ALAYyz36P84jzhNvfWzCxoBLWa; sid=r3gVGp1Oyg1qxQg4h478; Hm_lvt_2c8ad67df9e787ad29dbd54ee608f5d2=1603385391; Hm_lpvt_2c8ad67df9e787ad29dbd54ee608f5d2=1603385391; ec=05TnWrxT-1603385391907-17e8693248176-306326971; FSSBBIl1UgzbN7NP=5UyDgWm5pXcgqqqmTVs8csq7dytWwhH4y25pNG5H19ME0u5QpVp5yGW1sWbXIFwi1YQXUyXiJf0PqnVpZW4l0re346Qfd7P2wxKoCFtA0jLi310502_zYXKDt3iuTMb_OW9UBIwwinDF153d49AcZ3oIrUKZ597lG9Z8WI24ESgtVXTAvCI1TVO0SamCgSrlfg8ERL5TyzsEIsjZEkDhFrv7D3FoAnPlnePw6cq2gHwCkpWJAgBaFmfoOWnQ40zdIjeOhAad0HXUs_PFG_zzLGahcatkfH0qQXRPV6I3NSScUH6TCvKqqJsd0tcz9vgJq3; _efmdata=NJBqJAOKVsOQoSMdxBQbIqOuP77KBjHyrjeW4Tvnr9fkNdMBXHEGbm3haBzsEqWgrscgIYWBXAX4QhOAnXOeAJq%2FjTS30%2BfwIkwcmEpl3RA%3D; _exid=9x0rrwRW2f2DlvX2eIW4QWBcpjiObn0ZudL0lvS2SchSDvIaqXTVHWJccnQKAl0NFBagFHT0jD%2FEa%2BwlxYMNgw%3D%3D"
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
