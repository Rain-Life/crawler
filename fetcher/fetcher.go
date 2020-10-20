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
)

func Fetch(url string) ([]byte, error) {
	//resp, err := http.Get(url)
	//if err != nil {
	//	return nil, err
	//}
	//defer resp.Body.Close()

	client := &http.Client{}
	newUrl := strings.Replace(url, "http://", "https://", 1)
	req, err := http.NewRequest("GET", newUrl, nil)
	if err != nil {
		panic(err)
		return nil, err
	}
	//https://blog.csdn.net/kemuxiaozi000/article/details/107166072/ 模拟浏览器访问，避免403和202问题
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36")
	cookie1 := "sid=d2c9053b-b940-4906-ad36-6f21527d837f; Hm_lvt_2c8ad67df9e787ad29dbd54ee608f5d2=1603192106; FSSBBIl1UgzbN7NO=5LZ_1Lwx7erLryncImCKEYqO5BkSxbSWM0K4VdqzpWnbm4mZh3mDjvEvaSrlT7RwAqHJyW0pAKASRG4ucrnwTgq; ec=2l3cKCkA-1603192115645-0ea021cf190f51274025277; Hm_lpvt_2c8ad67df9e787ad29dbd54ee608f5d2=1603206501; _efmdata=NJBqJAOKVsOQoSMdxBQbIqOuP77KBjHyrjeW4Tvnr9fkNdMBXHEGbm3haBzsEqWgRUcHeuw67Zp%2BeX%2BSxDRojxzyn3x3aEkIgq6d5Hpm51w%3D; _exid=gtIARSAhoPZTl%2F3omhRCynF1W9MR9Rqs3VtmzonqsqUcR0xaWd9xl4Cftzd7qrg1ndleWAptOfYN9PaTzmji1A%3D%3D; FSSBBIl1UgzbN7NP=5Ug4eR25idKAqqqmTpN.pEaaa_uev7g8Q60l41UQtbMdzltpQ8GobLacQX8CZ0_LxYEF7GH6zEav9fwL27NBJqn.XlUc1TJXKLYZ_Ux2UuU0Myftkf8csSfVmteaKPVJQm_XhXlzHDtvd8c.WpqMjkmFvL.RLEyuGyF7nkGfU2_6yOjyZKYCxAafnngKBKPzun1tqKnktD_swxrYYnvhUIsO3mKRuuZiL7pA4c50uOl_9IJEcj9sp5h_VO_85IVnLjUItg3WusISTDXjQkEIxN_6lK..ggYLbyRgIpRbXwCzo.Mwbs8uwZzXutjc_mB.ga"
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
