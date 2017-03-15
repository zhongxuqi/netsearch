package spider

import (
	"strings"
  "github.com/PuerkitoBio/goquery"
	. "DesertEagleSite/bean"
)

func GetStackOverflowData(keyword string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument("http://stackoverflow.com/questions/tagged/" + keyword)
	if err != nil {
		return nil, "", err
	}
	return ParseStackOverflowHTML(resp)
}

func ParseStackOverflowUrl(url string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument(url)
	if err != nil {
		return nil, "", err
	}
	return ParseStackOverflowHTML(resp)
}

func ParseStackOverflowHTML(resp *goquery.Document) ([]DataItem, string, error) {
	resItems := make([]DataItem, 0)
  resp.Find("#questions .question-summary").Each(func(i int, s *goquery.Selection) {
    resItem := DataItem{}
    resItem.Title = strings.Replace(strings.Trim(
			s.Find(".summary h3 a").First().Text(), " \n"), "\n", " ", -1)
    resItem.Link = "http://stackoverflow.com" +
    	s.Find(".summary h3 a").First().AttrOr("href", "")
		if len(s.Find(".excerpt").Nodes) > 0 {
    	resItem.Abstract = strings.Replace(strings.Trim(
				s.Find(".excerpt").Text(), " \n"), "\n", " ", -1)
		}
    resItem.ImageUrl = s.Find("img").AttrOr("src", "")
    resItems = append(resItems, resItem)
  })
	nextPage := ""
	nextHtml := resp.Find(".pager a")
	if len(nextHtml.Nodes) > 0 {
    nextPage = "http://stackoverflow.com" + nextHtml.Last().AttrOr("href", "")
	}
	return resItems, nextPage, nil
}
