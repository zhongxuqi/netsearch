package spider

import (
	"strings"
  "github.com/PuerkitoBio/goquery"
	. "DesertEagleSite/bean"
)

func GetJianshuData(keyword string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument("http://www.jianshu.com/search?utf8=%E2%9C%93&q=" + keyword)
	if err != nil {
		return nil, "", err
	}
	return ParseJianShuHTML(resp)
}

func ParseJianShuUrl(url string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument(url)
	if err != nil {
		return nil, "", err
	}
	return ParseJianShuHTML(resp)
}

func ParseJianShuHTML(resp *goquery.Document) ([]DataItem, string, error) {
	resItems := make([]DataItem, 0)
  resp.Find("ul.list li").Each(func(i int, s *goquery.Selection) {
    resItem := DataItem{}
    resItem.Title = strings.Replace(strings.Trim(
			s.Find("h4.title a").First().Text(), " \n"), "\n", " ", -1)
    resItem.Link = s.Find("h4.title a").First().AttrOr("href", "")
		if len(s.Find("p").Nodes) > 0 {
    	resItem.Abstract = strings.Replace(strings.Trim(
				s.Find("p").Text(), " \n"), "\n", " ", -1)
		}
    resItem.ImageUrl = s.Find("img").AttrOr("src", "")
    resItems = append(resItems, resItem)
  })
	nextPage := ""
	nextHtml := resp.Find("div.pagination ul li.next a")
	if len(nextHtml.Nodes) > 0 {
    nextPage = "http://www.jianshu.com" + nextHtml.Last().AttrOr("href", "")
	}
	return resItems, nextPage, nil
}
