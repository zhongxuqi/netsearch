package spider

import (
	"strings"
  "github.com/PuerkitoBio/goquery"
	. "DesertEagleSite/bean"
)

func GetGithubData(keyword string) ([]DataItem, string, error) {
  resp, err := goquery.NewDocument("https://github.com/search?q=" + keyword)
	if err != nil {
		return nil, "", err
	}
	return ParseGithubHTML(resp)
}

func ParseGithubUrl(url string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument(url)
	if err != nil {
		return nil, "", err
	}
	return ParseGithubHTML(resp)
}

func ParseGithubHTML(resp *goquery.Document) ([]DataItem, string, error) {
	resItems := make([]DataItem, 0)
  resp.Find(".repo-list .repo-list-item").Each(func(i int, s *goquery.Selection) {
    resItem := DataItem{}
    resItem.Title = strings.Replace(strings.Trim(
			s.Find("h3 a").First().Text(), " \n"), "\n", " ", -1)
    resItem.Link = "https://github.com" +
      s.Find("h3 a").First().AttrOr("href", "")
		if len(s.Find(".repo-list-description").Nodes) > 0 {
    	resItem.Abstract = strings.Replace(strings.Trim(
				s.Find(".repo-list-description").Text(), " \n"), "\n", " ", -1)
		}
    resItem.ImageUrl = s.Find("img").AttrOr("src", "")
    resItems = append(resItems, resItem)
  })
	nextPage := ""
	nextHtml := resp.Find(".pagination a.next_page")
	if len(nextHtml.Nodes) > 0 {
    nextPage = "https://github.com" + nextHtml.Last().AttrOr("href", "")
	}
	return resItems, nextPage, nil
}
