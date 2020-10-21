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
	cookie1 := "sid=d2c9053b-b940-4906-ad36-6f21527d837f; Hm_lvt_2c8ad67df9e787ad29dbd54ee608f5d2=1603192106; FSSBBIl1UgzbN7NO=5LZ_1Lwx7erLryncImCKEYqO5BkSxbSWM0K4VdqzpWnbm4mZh3mDjvEvaSrlT7RwAqHJyW0pAKASRG4ucrnwTgq; ec=2l3cKCkA-1603192115645-0ea021cf190f51274025277; Hm_lpvt_2c8ad67df9e787ad29dbd54ee608f5d2=1603207359; _efmdata=NJBqJAOKVsOQoSMdxBQbIqOuP77KBjHyrjeW4Tvnr9fkNdMBXHEGbm3haBzsEqWgvnL2%2BNvWaW9y45eZTm057Nv9Lco%2BqEw7JdY4bwRM%2B5o%3D; _exid=g%2F659p8J4bB%2Fp%2F1xUygO0nyb1lBN5plx6fCZEQswlebTT2C%2BqFjb8B7MhNq%2BIWX6aHEES4mAcDH9hHdSBZRgwA%3D%3D; FSSBBIl1UgzbN7NP=5UgZqam5HqhgqqqmT1oZfcG4OPkB0FX1A36BJZ58u1ZlZjPLou9.60F1vbBlpcsd_u1K5D9a1.4tLEYDZdQShb9LyqV1eJwepjgp964rFkoovuq6h6QXN6g9mTUJX1sCjW939vKB81m6ncSNfOFJ3tPwXvQgCH4KSNnTphc9Eze6HUn54D4eRTBBpfOIaIOvOMZJWdLNrUf.9IUI1hEwZocuAqS_mG8D6mXrdqfUTqmKqtAZNY_S86Emzp5xg4924KXrM_pEx741K0OUMaOLB36Ho7LpzBkjf2_zOxZWQL128OAqcyXUfu.i.nKDqKZstW"
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
