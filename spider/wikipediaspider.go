package spider

import (
	"strings"
  "github.com/PuerkitoBio/goquery"
	. "DesertEagleSite/bean"
)

func GetWikipediaData(keyword string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument("https://wuu.wikipedia.org/w/index.php?title=Special:搜索&profile=default&fulltext=Search&search=" + keyword)
	if err != nil {
		return nil, "", err
	}
	return ParseWikipediaHTML(resp)
}

func ParseWikipediaUrl(url string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument(url)
	if err != nil {
		return nil, "", err
	}
	return ParseWikipediaHTML(resp)
}

func ParseWikipediaHTML(resp *goquery.Document) ([]DataItem, string, error) {
	resItems := make([]DataItem, 0)
  resp.Find("ul.mw-search-results li").Each(func(i int, s *goquery.Selection) {
    resItem := DataItem{}
    resItem.Title = strings.Replace(strings.Trim(
			s.Find("div.mw-search-result-heading a").First().Text(), " \n"), "\n", " ", -1)
    resItem.Link = "https://wuu.wikipedia.org" +
			s.Find("div.mw-search-result-heading a").First().AttrOr("href", "")
		if len(s.Find("div.searchresult").Nodes) > 0 {
    	resItem.Abstract = strings.Replace(strings.Trim(
				s.Find("div.searchresult").Text(), " \n"), "\n", " ", -1)
		}
    resItem.ImageUrl = s.Find("img").AttrOr("src", "")
    resItems = append(resItems, resItem)
  })
	nextPage := ""
	nextHtml := resp.Find("p.mw-search-pager-bottom a.mw-nextlink")
	if len(nextHtml.Nodes) > 0 {
    nextPage = "https://wuu.wikipedia.org" + nextHtml.Last().AttrOr("href", "")
	}
	return resItems, nextPage, nil
}
