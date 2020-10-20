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
	cookie1 := "sid=d2c9053b-b940-4906-ad36-6f21527d837f; Hm_lvt_2c8ad67df9e787ad29dbd54ee608f5d2=1603192106; FSSBBIl1UgzbN7NO=5LZ_1Lwx7erLryncImCKEYqO5BkSxbSWM0K4VdqzpWnbm4mZh3mDjvEvaSrlT7RwAqHJyW0pAKASRG4ucrnwTgq; ec=2l3cKCkA-1603192115645-0ea021cf190f51274025277; Hm_lpvt_2c8ad67df9e787ad29dbd54ee608f5d2=1603196084; _efmdata=NJBqJAOKVsOQoSMdxBQbIqOuP77KBjHyrjeW4Tvnr9fkNdMBXHEGbm3haBzsEqWg5iIjuHle5V1Hh0vasiwCLtpyQj4F1yV3it9Q8dxY0Ik%3D; _exid=vw7EGigLeOaInktYQgDstfqjEPgVQULcD345kwbEXYTJp3W6TXZsUSpPIjBzuXH6yQgtS72vom5GgMqnBETliA%3D%3D; FSSBBIl1UgzbN7NP=5Ug4Vzm5iThWqqqmTpIOADAA1qyRJBIm9RiIw5LK44rbSCzmwotKwfe9vTqN1fmiNq6XissxY9eLDIECMU5RpieoeH0SFweDC1GgStxVjrQD.4RszhqUYwU2XaSPfQ5qeBs3cDtd2oAByY52RmP2SGsqso7ltdSQjTnVQBz9RY0AvnXxDW2KGoI7_dVZuiXPViooREsdK.Offi.1hsQuDeXDlK41Cn66LM_KyBnX.yFOXcshqBPmNKh4RVUdJ_cBoSIS32Fqt.fCiuA5grbR4CTmAcOVAG3kadhQqaFNZEr8TTfkoQjIRV4waCw.EFy.9l"
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
