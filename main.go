package main

import (
	"crawler/engine"
	"crawler/zhenai/parser"
)

func main() {
	//body, err := fetcher.Fetch("https://album.zhenai.com/u/1765758535")
	//if err != nil {
	//	panic(err)
	//	//log.Printf("Fetcher: error " + "fetching url %s: %v", r.Url, err)
	//	return
	//}
	//
	//var Re = regexp.MustCompile(`<div class="m-btn purple" data-v-8b1eac0c>([^<]+)</div>`)
	//parser.ExtractString1(body, Re)

	engine.Run(engine.Request{
		Url: "http://www.zhenai.com/zhenghun",
		ParserFunc: parser.ParseCityList,
	})
}
