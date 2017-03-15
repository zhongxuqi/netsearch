package spider

import (
	"strings"
  "github.com/PuerkitoBio/goquery"
	. "DesertEagleSite/bean"
)

func GetCSDNData(keyword string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument("http://so.csdn.net/so/search/s.do?t=blog&o=&s=&q=" + keyword)
	if err != nil {
		return nil, "", err
	}
	return ParseCSDNHTML(resp)
}

func ParseCSDNUrl(url string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument(url)
	if err != nil {
		return nil, "", err
	}
	return ParseCSDNHTML(resp)
}

func ParseCSDNHTML(resp *goquery.Document) ([]DataItem, string, error) {
	resItems := make([]DataItem, 0)
  resp.Find(".add-tag-con, .search-list-con dl.search-list").
	Each(func(i int, s *goquery.Selection) {
    resItem := DataItem{}
		if len(s.Find("h3 a").Nodes) > 0 {
			resItem.Title = strings.Replace(strings.Trim(
				s.Find("h3 a").First().Text(), " \n"), "\n", " ", -1)
	    resItem.Link = s.Find("h3 a").First().AttrOr("href", "")
		} else if len(s.Find("dt a").Nodes) > 0 {
	    resItem.Title = strings.Replace(strings.Trim(
				s.Find("dt a").First().Text(), " \n"), "\n", " ", -1)
	    resItem.Link = s.Find("dt a").First().AttrOr("href", "")
		} else {
			return
		}
		if len(s.Find("dd.search-detail").Nodes) > 0 {
    	resItem.Abstract = strings.Replace(strings.Trim(
				s.Find("dd.search-detail").Text(), " \n"), "\n", " ", -1)
		}
    resItem.ImageUrl = s.Find("img").AttrOr("src", "")
    resItems = append(resItems, resItem)
  })
	nextPage := ""
	nextHtml := resp.Find("a.btn-next")
	if len(nextHtml.Nodes) > 0 {
    nextPage = "http://so.csdn.net/so/search/s.do" + nextHtml.Last().AttrOr("href", "")
	}
	return resItems, nextPage, nil
}
