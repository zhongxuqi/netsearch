package spider

import (
	"strings"
  "github.com/PuerkitoBio/goquery"
	. "DesertEagleSite/bean"
)

func GetDoubanData(keyword string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument("http://www.douban.com/search?q=" + keyword)
	if err != nil {
		return nil, "", err
	}

	resItems := make([]DataItem, 0)
  resp.Find("div.result").Each(func(i int, s *goquery.Selection) {
    resItem := DataItem{}
    resItem.Title = strings.Replace(strings.Trim(
			s.Find("div.title h3 a").First().Text(), " \n"), "\n", " ", -1)
    resItem.Link = s.Find("div.title h3 a").First().AttrOr("href", "")
		if len(s.Find("div.content p").Nodes) > 0 {
    	resItem.Abstract = strings.Replace(strings.Trim(
				s.Find("div.content p").Text(), " \n"), "\n", " ", -1)
		} else if len(s.Find(".c-row div p").Nodes) > 0 {
			resItem.Abstract = strings.Replace(strings.Trim(
				s.Find(".c-row div p").Text(), " \n"), "\n", " ", -1)
		}
    resItem.ImageUrl = s.Find("img").AttrOr("src", "")
    resItems = append(resItems, resItem)
  })
	nextPage := ""
	return resItems, nextPage, nil
}
