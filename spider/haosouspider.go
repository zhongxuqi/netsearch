package spider

import (
  "regexp"
	"strings"
  "github.com/PuerkitoBio/goquery"
	. "DesertEagleSite/bean"
)

func GetHaosouData(keyword string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument("https://www.so.com/s?ie=utf-8&q=" + keyword)
	if err != nil {
		return nil, "", err
	}
	return ParseHaosouHTML(resp)
}

func ParseHaosouUrl(url string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument(url)
	if err != nil {
		return nil, "", err
	}
	return ParseHaosouHTML(resp)
}

func ParseHaosouHTML(resp *goquery.Document) ([]DataItem, string, error) {
	resItems := make([]DataItem, 0)
  resp.Find("li.res-list").Each(func(i int, s *goquery.Selection) {
    resItem := DataItem{}
		if len(s.Find("h3.res-title a").Nodes) > 0 {
	    resItem.Title = strings.Replace(strings.Trim(
				s.Find("h3.res-title a").First().Text(), " \n"), "\n", " ", -1)
	    resItem.Link = s.Find("h3.res-title a").First().AttrOr("href", "")
		} else if len(s.Find("h3.title a").Nodes) > 0 {
	    resItem.Title = strings.Replace(strings.Trim(
				s.Find("h3.title a").First().Text(), " \n"), "\n", " ", -1)
	    resItem.Link = s.Find("h3.title a").First().AttrOr("href", "")
		} else {
			return
		}
		resItem.Link = GetHaosouRealUrl(resItem.Link)
		if len(s.Find("p.res-desc").Nodes) > 0 {
    	resItem.Abstract = strings.Replace(strings.Trim(
				s.Find("p.res-desc").Text(), " \n"), "\n", " ", -1)
		} else if len(s.Find(".res-rich p").Nodes) > 0 {
    	resItem.Abstract = strings.Replace(strings.Trim(
				s.Find(".res-rich p").Text(), " \n"), "\n", " ", -1)
		} else if len(s.Find(".mh-wrap").Nodes) > 0 {
    	resItem.Abstract = strings.Replace(strings.Trim(
				s.Find(".mh-wrap").Text(), " \n"), "\n", " ", -1)
		}
    resItem.ImageUrl = s.Find("img").AttrOr("src", "")
    resItems = append(resItems, resItem)
  })
	nextPage := ""
	nextHtml := resp.Find("a#snext")
	if len(nextHtml.Nodes) > 0 {
    nextPage = "https://www.so.com" + nextHtml.Last().AttrOr("href", "")
	}
	return resItems, nextPage, nil
}

func GetHaosouRealUrl(in_url string) (string) {
	if strings.Index(in_url, "www.so.com/link?url=") < 0 {
		return in_url
	}
	resp, err := goquery.NewDocument(in_url)
	if err != nil {
		return in_url
	}
	re := regexp.MustCompile("http://[^'\" ]+")
  url := re.FindString(resp.Find("noscript").Text())
	if len(url) > 0 {
		return url
	} else {
		return in_url
	}
}
