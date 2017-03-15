package handlefunc

import (
	"strings"
	"errors"
	"net/http"
	"encoding/json"
	"encoding/base64"
	"net/url"
	"DesertEagleSite/spider"
	. "DesertEagleSite/bean"
	"DesertEagleSite/config"
)

type ListResponse struct {
  BaseResponse
  Webs []SpiderObject
}

func initSpider() {
  urlFuncMap["app/list"] = ListSpiders

  for _, spider := range config.SpiderMap {
    urlFuncMap[spider.Url] = SearchData
  }
	urlFuncMap["app/submit_union_task"] = SearchUnion
	urlFuncMap["app/get_union_result"] = GetUnionResult
	urlFuncMap["app/submit_monitor_task"] = SubmitMonitorTask
	urlFuncMap["app/get_monitor_result"] = GetMonitorResult
	urlFuncMap["app/delete_monitor_task"] = DeleteMonitorTask
	urlFuncMap["app/custom_search"] = CustomSearch
	urlFuncMap["app/monitor_search/tasks"] = GetMonitorTasks
}

func writeSpiderResult(w http.ResponseWriter, r *http.Request, resItems []DataItem, nextPage string, err error) {
	var response SpiderResponse
	if err != nil {
		response.Status = "500"
		response.Message = err.Error()
	} else {
		response.Status = "200"
		response.Message = "search success"
		response.ResultData = resItems
		response.NextPage = nextPage
	}
	respBytes, _ := json.Marshal(response)
	w.Write(respBytes)
}

func ListSpiders(w http.ResponseWriter, r *http.Request) {
  webs := make([]SpiderObject, 0)
  for _, spider := range config.SpiderMap {
    webs = append(webs, spider)
  }
  var response ListResponse
  response.Status = "200"
  response.Message = "success"
  response.Webs = webs
  respBytes, _ := json.Marshal(response)
  w.Write(respBytes)
}

func SearchData(w http.ResponseWriter, r *http.Request) {
	if keyword, ok := parseKeyword(r, "keyword"); ok {
		for _, spider := range config.SpiderMap {
			if spider.Url == r.URL.Path[1:] {
				resItems, nextPage, err := spider.GetDataFunc(keyword)
				writeSpiderResult(w, r, resItems, base64.URLEncoding.EncodeToString([]byte(nextPage)), err)
				break
			}
		}
	}
}

func parseSpiders(parser_names string) (spiderList []SpiderObject) {
  for _, parser_name := range strings.Split(parser_names, ",") {
    for _, spider := range config.SpiderMap {
      if parser_name == spider.ParserName {
        spiderList = append(spiderList, spider)
        break
      }
    }
  }
  return
}

func SearchUnion(w http.ResponseWriter, r *http.Request) {
	paramsMap, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		writeResult(w, r, "", err)
		return
	}
	keyword := paramsMap.Get("keyword");
	parser_names := paramsMap.Get("parser_names");
	registration_id := paramsMap.Get("registration_id");
	if len(keyword) == 0 || len(parser_names) == 0 || len(registration_id) == 0 {
		writeResult(w, r, "", errors.New("argument is error"))
		return
	}
  spiderList := parseSpiders(parser_names)
  if len(spiderList) == 0 {
    writeResult(w, r, "", errors.New("spider list is none"))
		return
  }
	go spider.GetUnionData(keyword, parser_names, registration_id, spiderList)
	writeResult(w, r, "task has submitted.", nil)
}

func GetUnionResult(w http.ResponseWriter, r *http.Request) {
	mapkey, ok := parseKeyword(r, "map_key")
	if !ok {
		writeResult(w, r, "", errors.New("argument is error"))
		return
	}
	respBytes, _ := json.Marshal(spider.GetResultByKey(mapkey))
	w.Write(respBytes)
}

func SubmitMonitorTask(w http.ResponseWriter, r *http.Request) {
	paramsMap, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		writeResult(w, r, "", err)
		return
	}
	keyword := paramsMap.Get("keyword");
	registration_id := paramsMap.Get("registration_id");
	target_url := paramsMap.Get("target_url");
	if len(keyword) == 0 || len(registration_id) == 0 || len(target_url) == 0 {
		writeResult(w, r, "", errors.New("argument is error"))
		return
	}
	decode_url, err := base64.URLEncoding.DecodeString(target_url)
	if err != nil {
		writeResult(w, r, "", err)
		return
	}
	go spider.SubmitRawMonitorTask(keyword, registration_id, string(decode_url))
	writeResult(w, r, "task has submitted.", nil)
}

func GetMonitorResult(w http.ResponseWriter, r *http.Request) {
	mapkey, ok := parseKeyword(r, "map_key")
	if !ok {
		writeResult(w, r, "", errors.New("argument is error"))
		return
	}
	respBytes, _ := json.Marshal(spider.GetMonitorResultByKey(mapkey))
	w.Write(respBytes)
}

func GetMonitorTasks(w http.ResponseWriter, r *http.Request) {
	var response MonitorListResponse
	response.Status = "200"
	response.Message = "success"
	response.Tasks = spider.GetMonitorTaskList()
	respBytes, _ := json.Marshal(response)
	w.Write(respBytes)
}

func DeleteMonitorTask(w http.ResponseWriter, r *http.Request) {
	paramsMap, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		writeResult(w, r, "", err)
		return
	}
	keyword := paramsMap.Get("keyword");
	registration_id := paramsMap.Get("registration_id");
	target_url := paramsMap.Get("target_url");
	if len(keyword) == 0 || len(registration_id) == 0 || len(target_url) == 0 {
		writeResult(w, r, "", errors.New("argument is error"))
		return
	}
	decode_url, err := base64.URLEncoding.DecodeString(target_url)
	task := &MonitorTask{
		Url: string(decode_url),
		Keyword: keyword,
		RegistrationId: registration_id,
	}
	go spider.DeleteMonitorTask(task)
	writeResult(w, r, "task has deleted.", nil)
}

func CustomSearch(w http.ResponseWriter, r *http.Request) {
	PaserName, _ := parseKeyword(r, "parser_name")
	url, _ := parseKeyword(r, "url")
	if len(PaserName) == 0 || len(url) == 0 {
		return
	}
	decodeUrl, err := base64.URLEncoding.DecodeString(url)
	if err != nil {
		return
	}
	url = string(decodeUrl)
	for _, spider := range config.SpiderMap {
		if spider.ParserName == PaserName {
			resItems, nextPage, err := spider.ParseFunc(url)
			writeSpiderResult(w, r, resItems, base64.URLEncoding.EncodeToString([]byte(nextPage)), err)
			break
		}
	}
}
