package spider

import (
  "DesertEagleSite/config"
  . "DesertEagleSite/bean"
)

const (
  BAIDU = "baidu"
  ZHIHU = "zhihu"
  HAOSOU = "haosou"
  WIKIPEDIA = "wikipedia"
  BAIDUXUESHU = "baiduxuehsu"
  DOUBAN = "douban"
  JIANSHU = "jianshu"
  CSDN = "csdn"
  BING = "bing"
	GOOGLE = "google"
	STACKOVERFLOW = "Stack Overflow"
	GITHUB = "github"
  SOUGOU = "sougou"

  UNION_TYPE = 1
  MONITOR_TYPE = 2
)

var SpiderMap = map[string] SpiderObject{
  BAIDU: SpiderObject{
    Url: "app/search_baidu",
    Name: "百度",
    ParserName: "Baidu",
    Logo: "/data/baidu_logo.png",
		GetDataFunc: GetBaiduData,
		ParseFunc: ParseBaiduUrl,
  },
  ZHIHU: SpiderObject{
    Url: "app/search_zhihu",
    Name: "知乎",
    ParserName: "Zhihu",
    Logo: "/data/zhihu_logo.png",
		GetDataFunc: GetZhihuData,
		ParseFunc: ParseZhihuUrl,
  },
  HAOSOU: SpiderObject{
    Url: "app/search_haosou",
    Name: "好搜",
    ParserName: "Haosou",
    Logo: "/data/haosou_logo.png",
		GetDataFunc: GetHaosouData,
		ParseFunc: ParseHaosouUrl,
  },
  WIKIPEDIA: SpiderObject{
    Url: "app/search_wikipedia",
    Name: "维基百科",
    ParserName: "Wikipedia",
    Logo: "/data/wikipedia_logo.png",
		GetDataFunc: GetWikipediaData,
		ParseFunc: ParseWikipediaUrl,
  },
  BAIDUXUESHU: SpiderObject{
    Url: "app/search_baiduxueshu",
    Name: "百度学术",
    ParserName: "BaiduXueShu",
    Logo: "/data/baidu_logo.png",
		GetDataFunc: GetBaiduXueShuData,
		ParseFunc: ParseBaiduXueShuUrl,
  },
  // DOUBAN: {
  //   Url: "app/search_douban",
  //   Name: "豆瓣",
  //   ParserName: "",
  // },
  // JIANSHU: spider.SpiderObject{
  //   Url: "app/search_jianshu",
  //   Name: "简书",
  //   ParserName: "JianShuData",
  //   Logo: "/data/jianshu_logo.png",
	// 	GetDataFunc: spider.GetJianshuData,
	// 	ParseFunc: spider.ParseJianShuUrl,
  // },
  CSDN: SpiderObject{
    Url: "app/search_csdn",
    Name: "CSDN",
    ParserName: "CSDN",
    Logo: "/data/csdn_logo.png",
		GetDataFunc: GetCSDNData,
		ParseFunc: ParseCSDNUrl,
  },
	BING: SpiderObject{
    Url: "app/search_bing",
    Name: "Bing",
    ParserName: "Bing",
    Logo: "/data/bing_logo.png",
		GetDataFunc: GetBingData,
		ParseFunc: ParseBingUrl,
  },
	GOOGLE: SpiderObject{
		Url: "app/search_google",
    Name: "Google",
    ParserName: "Google",
    Logo: "/data/google_logo.png",
		GetDataFunc: GetGoogleData,
		ParseFunc: ParseGoogleUrl,
	},
	STACKOVERFLOW: SpiderObject{
		Url: "app/search_stackoverflow",
    Name: "StackOverflow",
    ParserName: "StackOverflow",
    Logo: "/data/stackoverflow_logo.png",
		GetDataFunc: GetStackOverflowData,
		ParseFunc: ParseStackOverflowUrl,
	},
	GITHUB: SpiderObject{
		Url: "app/search_github",
    Name: "Github",
    ParserName: "Github",
    Logo: "/data/github_logo.png",
		GetDataFunc: GetGithubData,
		ParseFunc: ParseGithubUrl,
	},
  SOUGOU: SpiderObject{
		Url: "app/search_sougou",
    Name: "搜狗",
    ParserName: "Sougou",
    Logo: "/data/sougou_logo.png",
		GetDataFunc: GetSougouData,
		ParseFunc: ParseSougouUrl,
	},
}

func init() {
  for k, v := range SpiderMap {
    config.SpiderMap[k] = v
  }
}
