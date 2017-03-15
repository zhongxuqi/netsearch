package spider

import (
	"strings"
  "regexp"
  "net/url"
  "github.com/PuerkitoBio/goquery"
	. "DesertEagleSite/bean"
)

func GetSougouData(keyword string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument("https://www.sogou.com/web?query=" + keyword)
	if err != nil {
		return nil, "", err
	}
	return ParseSougouHTML(resp)
}

func ParseSougouUrl(url string) ([]DataItem, string, error) {
	resp, err := goquery.NewDocument(url)
	if err != nil {
		return nil, "", err
	}
	return ParseSougouHTML(resp)
}

func ParseSougouHTML(resp *goquery.Document) ([]DataItem, string, error) {
	resItems := make([]DataItem, 0)
  resp.Find(".results .vrwrap, .results .rb, .results .vrPic").Each(func(i int, s *goquery.Selection) {
    resItem := DataItem{}
		if len(s.Find(".vrTitle a").Nodes) > 0 {
	    resItem.Title = strings.Replace(strings.Trim(
				s.Find(".vrTitle a").First().Text(), " \n"), "\n", " ", -1)
      if len(resItem.Title) <= 0 {
        resItem.Title = strings.Replace(strings.Trim(
  				s.Find(".vrTitle").First().Text(), " \n"), "\n", " ", -1)
      }
	    resItem.Link = s.Find(".vrTitle a").First().AttrOr("href", "")
		} else if len(s.Find("h3 a").Nodes) > 0 {
      resItem.Title = strings.Replace(strings.Trim(
				s.Find("h3 a").First().Text(), " \n"), "\n", " ", -1)
	    resItem.Link = s.Find("h3 a").First().AttrOr("href", "")
    } else {
			return
		}
		if tmp, err := url.QueryUnescape(resItem.Title); err == nil {
			resItem.Title = tmp
		}
    resItem.Link = GetSougouRealUrl(resItem.Link)
		if len(s.Find(".str_info").Nodes) > 0 {
    	resItem.Abstract = strings.Replace(strings.Trim(
				s.Find(".str_info").Text(), " \n"), "\n", " ", -1)
		} else if len(s.Find(".ft").Nodes) > 0 {
      resItem.Abstract = strings.Replace(strings.Trim(
				s.Find(".ft").Text(), " \n"), "\n", " ", -1)
    } else if len(s.Find(".div-p2 p").Nodes) > 0 {
      resItem.Abstract = strings.Replace(strings.Trim(
				s.Find(".div-p2 p").Text(), " \n"), "\n", " ", -1)
      if len(resItem.Abstract) <= 0 {
        resItem.Abstract = strings.Replace(strings.Trim(
  				s.Find(".div-p2 div").First().Text(), " \n"), "\n", " ", -1)
      }
    } else if len(s.Find(".tribal-introduction").Nodes) > 0 {
      resItem.Abstract = strings.Replace(strings.Trim(
				s.Find(".tribal-introduction").Text(), " \n"), "\n", " ", -1)
    } else if len(s.Find("tbody").Nodes) > 0 {
      resItem.Abstract = strings.Replace(strings.Trim(
				s.Find("tbody").Text(), " \n"), "\n", " ", -1)
    } else if len(s.Find(".biz_ft").Nodes) > 0 {
      resItem.Abstract = strings.Replace(strings.Trim(
				s.Find(".biz_ft").Text(), " \n"), "\n", " ", -1)
    } else if len(s.Find("p").Nodes) > 0 {
      resItem.Abstract = strings.Replace(strings.Trim(
				s.Find("p").Text(), " \n"), "\n", " ", -1)
    } else if len(s.Find(".str-box-v4").Nodes) > 0 {
      resItem.Abstract = strings.Replace(strings.Trim(
				s.Find(".str-box-v4").Text(), " \n"), "\n", " ", -1)
    } else if len(s.Find(".wx-box-new").Nodes) > 0 {
      resItem.Abstract = strings.Replace(strings.Trim(
				s.Find(".wx-box-new").Text(), " \n"), "\n", " ", -1)
    }
		if tmp, err := url.QueryUnescape(resItem.Abstract); err == nil {
			resItem.Abstract = tmp
		}
    resItem.ImageUrl = s.Find("img").AttrOr("src", "")
    resItems = append(resItems, resItem)
  })
	nextPage := ""
	nextHtml := resp.Find("a#sogou_next")
	if len(nextHtml.Nodes) > 0 {
    nextPage = "https://www.sogou.com/web" + nextHtml.Last().AttrOr("href", "")
	}
	return resItems, nextPage, nil
}

func GetSougouRealUrl(in_url string) (string) {
	if strings.Index(in_url, "www.sogou.com/link?url=") < 0 {
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
