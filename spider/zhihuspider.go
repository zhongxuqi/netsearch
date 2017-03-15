package spider

import (
	"bytes"
	"strings"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"net/url"
  "github.com/PuerkitoBio/goquery"
	. "DesertEagleSite/bean"
)

type Page struct {
	Next string `json:"next"`
}

type ZhiHuData struct {
	Paging Page `json:"paging"`
	Htmls []string `json:"htmls"`
}

func GetZhihuData(keyword string) ([]DataItem, string, error) {
	resp, err := http.Get("http://www.zhihu.com/r/search?range=&type=question&offset=0&q=" + keyword)
	if err != nil {
		return nil, "", err
	}
	return ParseZhihuHTML(resp)
}

func ParseZhihuUrl(url string) ([]DataItem, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	return ParseZhihuHTML(resp)
}

func ParseZhihuHTML(resp *http.Response) ([]DataItem, string, error) {
	resJson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	var zhihudata ZhiHuData
	err = json.Unmarshal(resJson, &zhihudata)
	if err != nil {
		return nil, "", err
	}

	resItems := make([]DataItem, 0)
	for _, htmlnode := range zhihudata.Htmls {
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(htmlnode)))
		if err != nil {
			continue
		}
		resItem := DataItem{}
    resItem.Title = strings.Replace(strings.Trim(
			doc.Find("div.title a").First().Text(), " \n"), "\n", " ", -1)
		if tmp, err := url.QueryUnescape(resItem.Title); err == nil {
			resItem.Title = tmp
		}
    resItem.Link = "http://www.zhihu.com" +
			doc.Find("div.title a").First().AttrOr("href", "")
		if len(doc.Find(".summary").Nodes) > 0 {
    	resItem.Abstract = strings.Replace(strings.Trim(
				doc.Find(".summary").Text(), " \n"), "\n", " ", -1)
		}
		if tmp, err := url.QueryUnescape(resItem.Abstract); err == nil {
			resItem.Abstract = tmp
		}
    resItem.ImageUrl = doc.Find("img").AttrOr("src", "")
    resItems = append(resItems, resItem)
	}
	nextPage := "http://www.zhihu.com" + zhihudata.Paging.Next
	return resItems, nextPage, nil
}
