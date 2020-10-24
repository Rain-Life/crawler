package parser

import (
	"crawler/engine"
	"crawler/model"
	"regexp"
)

var ageRe = regexp.MustCompile(`<td><span class="label">年龄：</span>(\d+)岁</td>`)
var heightRe = regexp.MustCompile(`<td><span class="label">身高：</span>(\d+)CM</td>`)
var incomeRe = regexp.MustCompile(`<td><span class="label">月收入：</span>([^<]+)</td>`)
var weightRe = regexp.MustCompile(`<td><span class="label">体重：</span><span field="">(\d+)KG</span></td>`)
var genderRe = regexp.MustCompile(`<td><span class="label">性别：</span><span field="">([^<]+)</span></td>`)
var xinzuoRe = regexp.MustCompile(`<td><span class="label">星座：</span><span field="">([^<]+)</span></td>`)
var marriageRe = regexp.MustCompile(`<td><span class="label">婚况：</span>([^<]+)</td>`)
var educationRe = regexp.MustCompile(`<td><span class="label">学历：</span>([^<]+)</td>`)
var occupationRe = regexp.MustCompile(`<td><span class="label">职业：</span><span field="">([^<]+)</span></td>`)
var hokouRe = regexp.MustCompile(`<td><span class="label">籍贯：</span>([^<]+)</td>`)
var houseRe = regexp.MustCompile(`<td><span class="label">住房条件：</span><span field="">([^<]+)</span></td>`)
var carRe = regexp.MustCompile(`<td><span class="label">是否购车：</span><span field="">([^<]+)</span></td>`)
var idUrlRe = regexp.MustCompile(`http://album.zhenai.com/u/([\d]+)`)

func ParseProfile(contents []byte, url string, name string) engine.ParseResult {
	profile := model.Profile{}

	profile.Name = name

	profile.Weight = extractString(contents, weightRe)
	profile.Age = extractString(contents, ageRe)
	profile.Height = extractString(contents, heightRe)
	profile.Income = extractString(contents, incomeRe)
	profile.Gender = extractString(contents, genderRe)
	profile.Car = extractString(contents, carRe)
	profile.Education = extractString(contents, educationRe)
	profile.Hokou = extractString(contents, hokouRe)
	profile.House = extractString(contents, houseRe)
	profile.Marriage = extractString(contents, marriageRe)
	profile.Occupation = extractString(contents, occupationRe)
	profile.Xinzuo = extractString(contents, xinzuoRe)

	result := engine.ParseResult{
		Items: []engine.Item {
			{
				Url: url,
				Type: "zhenai",
				Id: extractString([]byte(url), idUrlRe),
				Payload: profile,
			},
		},
	}

	return result
}

func extractString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)

	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}

var Re = regexp.MustCompile(`<div class="m-btn purple" data-v-8b1eac0c>([^<]+)</div>`)

func ParseProfile1(contents []byte,url string, name string) engine.ParseResult {
	profile := model.Profile{}

	profile.Name = name

	matches := ExtractString1(contents, Re)
	a := make(map[int]string)
	for i, m := range matches {
		a[i] = string(m[1])
	}

	profile.Gender = "男"
	profile.Age = a[1]
	profile.Height = a[3]
	profile.Weight = a[4]
	profile.Income = a[6]
	profile.Marriage = a[0]
	profile.Education = a[8]
	profile.Occupation = a[7]
	profile.Hokou = "未知"
	profile.Xinzuo = a[2]
	profile.House = "有房"
	profile.Car = "劳斯莱斯幻影"

	result := engine.ParseResult{
		Items: []engine.Item {
			{
				Url: url,
				Type: "zhenai",
				Id: extractString([]byte(url), idUrlRe),
				Payload: profile,
			},
		},
	}

	return result
}

func ExtractString1(contents []byte, re *regexp.Regexp) [][][]byte {
	matches := re.FindAllSubmatch(contents, -1)

	return matches
}
