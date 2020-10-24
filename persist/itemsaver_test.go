package persist

import (
	"context"
	"crawler/engine"
	"crawler/model"
	"encoding/json"
	"gopkg.in/olivere/elastic.v5"
	"testing"
)

func TestSave(t *testing.T) {
	 expected := engine.Item{
	 	Url: "http://album.zhenai.com/u/108906739",
	 	Type: "zhenai",
	 	Id: "108906739",
	 	Payload: model.Profile{
			Name:       "安静的雪",
			Gender:     "女",
			Age:        "34岁",
			Height:     "162CM",
			Weight:     "57KG",
			Income:     "3001-5000元",
			Marriage:   "离异",
			Education:  "大学本科",
			Occupation: "人事/行政",
			Hokou:      "山东菏泽",
			Xinzuo:     "牡羊座",
			House:      "已购房",
			Car:        "未购车",
		},
	 }

	err := save(expected)
	if err != nil {
		panic(err)
	}

	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	resp, err := client.Get().
		Index("dating_profile").
		Type(expected.Type).
		Id(expected.Id).
		Do(context.Background())
	if err != nil {
		panic(err)
	}

	t.Logf("%s", resp.Source)
	var actual engine.Item
	err = json.Unmarshal(*resp.Source, &actual)
	if err != nil {
		panic(err)
	}

	actualProfile, _ := model.FromJsonObj(actual.Payload)
	actual.Payload = actualProfile
	if actual != expected {
		t.Errorf("got %v; expected %v", actual, expected)
	}
}
