package spider

import (
  "fmt"
  . "DesertEagleSite/bean"
  "DesertEagleSite/util"
  "DesertEagleSite/push_manager"
)

var mSearchResultMap = make(map[string] UnionResponse)

func GetUnionData(keyword, parser_names, registration_id string, spiderList []SpiderObject) {
  UrlList := make([]DataItem, 0)
  for _, spider := range spiderList {
    oldSize := len(UrlList)
    itemList, nextPage, err := spider.GetDataFunc(keyword)
    if err == nil && itemList != nil {
      for _, item := range itemList {
        UrlList = append(UrlList, item)
      }
    }

    // get next page data
    itemList, _, err = spider.ParseFunc(nextPage)
    if err == nil && itemList != nil {
      for _, item := range itemList {
        UrlList = append(UrlList, item)
      }
    }
    fmt.Println("search in: ", spider.Name, ", size: ", len(UrlList) - oldSize)
  }

  resultList := execTasks(UrlList, keyword, false)

  mapKey := registration_id + "-" + util.GetFormatTimeNow()
  var response UnionResponse
  response.ResultData = make([]UrlResult, 0)
  for _, item := range resultList {
    fmt.Println(item.ToString())
    response.ResultData = append(response.ResultData, *item)
  }
	response.Status = "200"
	response.Message = "search success"
  response.ParserNames = parser_names
  response.Keyword = keyword
  mSearchResultMap[mapKey] = response
	var message PushMessage
	message.MapKey = mapKey
  message.Keyword = keyword
  message.Type = UNION_TYPE
	push_manager.PushJPushMessage(registration_id, util.ConvObject2Json(message))
}

func GetResultByKey(mapkey string) UnionResponse {
  resp, ok := mSearchResultMap[mapkey]
  if ok {
    delete(mSearchResultMap, mapkey)
    return resp
  } else {
    var response UnionResponse
    response.Status = "400"
  	response.Message = "has not the map key"
    return response
  }
}
