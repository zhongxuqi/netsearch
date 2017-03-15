package spider

import (
  "strings"
  "net/http"
  "github.com/PuerkitoBio/goquery"
  "DesertEagleSite/util"
  . "DesertEagleSite/bean"
)

func GetGoogleData(keyword string) ([]DataItem, string, error) {
  resp, err := http.Get("http://i.xsou.co/search?q=" + keyword)
	if err != nil {
		return nil, "", err
	}
  b, err := util.DecodeResponse2Utf8Bytes(resp)
  if err != nil {
		return nil, "", err
	}
  doc, err := goquery.NewDocumentFromReader(util.ConvBytes2Reader(b))
  if err != nil {
		return nil, "", err
	}
	return ParseGoogleHTML(doc)
}

func ParseGoogleUrl(url string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument(url)
	if err != nil {
		return nil, "", err
	}
	return ParseGoogleHTML(resp)
}

func ParseGoogleHTML(resp *goquery.Document) ([]DataItem, string, error) {
  resItems := make([]DataItem, 0)
  resp.Find(".g").Each(func(i int, s *goquery.Selection) {
    resItem := DataItem{}
    resItem.Title = strings.Replace(strings.Trim(
			s.Find(".r a").First().Text(), " \n"), "\n", " ", -1)
    resItem.Link = s.Find(".r a").First().AttrOr("href", "")
    startIndex := strings.Index(resItem.Link, "http://")
    if startIndex < 0 {
      startIndex = strings.Index(resItem.Link, "https://")
    }
    if startIndex >= 0 {
      resItem.Link = resItem.Link[startIndex:
        strings.Index(resItem.Link[startIndex:], "&") + startIndex]
    } else {
      return
    }
    resItem.Abstract = s.Find(".s .st").Text()
    resItem.ImageUrl = s.Find("img").AttrOr("src", "")
    resItems = append(resItems, resItem)
  })
	nextPage := ""
	nextHtml := resp.Find(".fl")
	if len(nextHtml.Nodes) > 0 {
    nextPage = "http://i.xsou.co" + nextHtml.Last().AttrOr("href", "")
	}
	return resItems, nextPage, nil
}
