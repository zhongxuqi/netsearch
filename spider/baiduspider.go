package spider

import (
	"regexp"
	"strings"
  "github.com/PuerkitoBio/goquery"
	"DesertEagleSite/bean"
)

func GetBaiduData(keyword string) ([]bean.DataItem, string, error) {
	resp, err := goquery.NewDocument("http://www.baidu.com/s?ie=uft-8&word=" + keyword)
	if err != nil {
		return nil, "", err
	}
	return ParseBaiduHTML(resp)
}

func ParseBaiduUrl(url string) ([]bean.DataItem, string, error) {
	resp, err := goquery.NewDocument(url)
	if err != nil {
		return nil, "", err
	}
	return ParseBaiduHTML(resp)
}

func ParseBaiduHTML(resp *goquery.Document) ([]bean.DataItem, string, error) {
	resItems := make([]bean.DataItem, 0)
  resp.Find(".result, .result-op").Each(func(i int, s *goquery.Selection) {
    resItem := bean.DataItem{}
    resItem.Title = strings.Replace(strings.Trim(
			s.Find("h3").Text(), " \n"), "\n", " ", -1)
    resItem.Link = s.Find("h3 a").First().AttrOr("href", "")
		if len(resItem.Link) == 0 {
			return
		}
		resItem.Link = GetBaiduRealUrl(resItem.Link)
		if len(s.Find(".c-abstract").Nodes) > 0 {
    	resItem.Abstract = strings.Replace(strings.Trim(
				s.Find(".c-abstract").Text(), " \n"), "\n", " ", -1)
		} else if len(s.Find(".c-row div p").Nodes) > 0 {
			resItem.Abstract = strings.Replace(strings.Trim(
				s.Find(".c-row div p").Text(), " \n"), "\n", " ", -1)
		} else if len(s.Find(".op-tieba-general-firsttd").Nodes) > 0 {
			resItem.Abstract = strings.Replace(strings.Trim(
				s.Find(".op-tieba-general-firsttd").Text(), " \n"), "\n", " ", -1)
		}
    resItem.ImageUrl = s.Find("img").AttrOr("src", "")
    resItems = append(resItems, resItem)
  })
	nextPage := ""
	nextHtml := resp.Find("a.n")
	if len(nextHtml.Nodes) > 0 {
    nextPage = "http://www.baidu.com" + nextHtml.Last().AttrOr("href", "")
	}
	return resItems, nextPage, nil
}

func GetBaiduRealUrl(in_url string) (string) {
	if strings.Index(in_url, "www.baidu.com/link?url=") < 0 {
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

func GetBaiduXueShuData(keyword string) ([]bean.DataItem, string, error) {
	resp, err := goquery.NewDocument("http://xueshu.baidu.com/s?ie=uft-8&wd="+ keyword)
	if err != nil {
		return nil, "", err
	}
	return ParseBaiduXueShuHTML(resp)
}

func ParseBaiduXueShuUrl(url string) ([]bean.DataItem, string, error) {
	resp, err := goquery.NewDocument(url)
	if err != nil {
		return nil, "", err
	}
	return ParseBaiduXueShuHTML(resp)
}

func ParseBaiduXueShuHTML(resp *goquery.Document) ([]bean.DataItem, string, error) {
	resItems := make([]bean.DataItem, 0)
  resp.Find(".result").Each(func(i int, s *goquery.Selection) {
    resItem := bean.DataItem{}
    resItem.Title = s.Find("div.sc_content h3.t a").First().Text()
    resItem.Link = "http://xueshu.baidu.com" + s.Find("div.sc_content h3.t a").First().AttrOr("href", "")
		if len(s.Find(".c_abstract").Nodes) > 0 {
			resItem.Abstract = strings.Replace(strings.Trim(
				s.Find(".c_abstract").Text(), " \n"), "\n", " ", -1)
		}
    resItem.ImageUrl = s.Find("img").AttrOr("src", "")
    resItems = append(resItems, resItem)
  })
	nextPage := ""
	nextHtml := resp.Find("a.n")
	if len(nextHtml.Nodes) > 0 {
    nextPage = "http://xueshu.baidu.com" + nextHtml.Last().AttrOr("href", "")
	}
	return resItems, nextPage, nil
}
