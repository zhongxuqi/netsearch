package spider

import (
  "fmt"
  "time"
  "sync"
  "strings"
  . "DesertEagleSite/bean"
  "DesertEagleSite/util"
  "DesertEagleSite/wordtool"
  // "DesertEagleSite/evaluator"
  "github.com/PuerkitoBio/goquery"
  "DesertEagleSite/push_manager"
)

var mutex = &sync.Mutex{}
var mTaskList = make([]*MonitorTask, 0)
var mMonitorResultMap = make(map[string] MonitorResponse)

func init() {
  go loop()
}

func loop() {
  for {
    fmt.Println("go to sleep.")
    time.Sleep(time.Minute)
    fmt.Println("go to wake.")
    mutex.Lock()

    for i, task := range mTaskList {
      fmt.Println("begin: ", task.Url)
      if time.Now().Before(task.Time) {
        continue
      }
      task.Life--
      for time.Now().After(task.Time) {
        task.Time = task.Time.Add(2 * time.Hour)
      }
      mTaskList[i] = task

      // get response struct
      response, ok := mMonitorResultMap[task.RespMapKey]
      if !ok {
        DeleteMonitorTaskByIndex(i)
        continue
      }
      response.ResultData = make([]UrlResult, 0)

      // eval main url
      subtask := UrlResult{}
      subtask.Task.Url = task.Url
      subtask.Task.Keywords = task.Keywords
      subtask.Task.Info.Title = "Main Page"
      subtask.Task.Info.Link = task.Url
      doc, err := goquery.NewDocument(subtask.Task.Url)
    	if err != nil {
        fmt.Println(err.Error())
        mutex.Unlock()
        DeleteMonitorTask(task)
        mutex.Lock()
    		continue
    	}
      //subtask.Eval = evaluator.EvaluateContentByKeyWords(doc.Find("body").Text(), subtask.Task.Keywords)
      //response.ResultData = append(response.ResultData, subtask)

      // parser a tag from url
      UrlList := make([]DataItem, 0)
      doc.Find("a").Each(func(i int, s *goquery.Selection) {
        subtask := DataItem{}
        subtask.Title = strings.Replace(strings.Trim(
    			s.Text(), " \n"), "\n", " ", -1)
        subtask.Link = strings.Replace(strings.Trim(
    			s.First().AttrOr("href", ""), " \n"), "\n", " ", -1)
        if len(subtask.Link) == 0 {
          return
        }
        UrlList = append(UrlList, subtask)
      })
      fmt.Println(subtask.Task.Info.Title, " size: ", len(UrlList))

      // get result list
      resultList := execTasks(UrlList, task.Keyword, true)
      fmt.Println("reciver size: ", len(resultList))

      // store the response to map
      mapKey := task.RegistrationId + "-" + util.GetFormatTimeNow()
      task.RespMapKey = mapKey
      response.ResultData = make([]UrlResult, 0)
      for _, item := range resultList {
        fmt.Println(item.ToString())
        response.ResultData = append(response.ResultData, *item)
      }
    	response.Status = "200"
    	response.Message = "search success"
      response.Task = *task
      response.Title = strings.Replace(strings.Trim(
  			doc.Find("title").Text(), " \n"), "\n", " ", -1)
      response.Url = task.Url
      response.Keyword = task.Keyword
      mMonitorResultMap[mapKey] = response

      // notice the client
    	var message PushMessage
    	message.MapKey = mapKey
      message.Keyword = task.Keyword
      message.Type = MONITOR_TYPE
    	ret := push_manager.PushJPushMessage(task.RegistrationId, util.ConvObject2Json(message))
      if !ret || task.Life <= 0 {
        DeleteMonitorTaskByIndex(i)
      }
    }
    mutex.Unlock()
  }
}

func SubmitMonitorTask(task *MonitorTask) bool {
  for _, item := range mTaskList {
    if task.IsEqual(item) {
      item.Life = MONITOR_TASK_LIFE_DEFAULT
      return false
    }
  }

  // add task to list
  mutex.Lock()
  task.Time = time.Now()
  task.Keywords = wordtool.SplitContent2Words(task.Keyword)
  task.Life = MONITOR_TASK_LIFE_DEFAULT
  mTaskList = append(mTaskList, task)
  mutex.Unlock()
  fmt.Println(len(mTaskList))

  // add response to list
  mapKey := task.RegistrationId + "-" + util.GetFormatTimeNow()
  task.RespMapKey = mapKey
  var response MonitorResponse
  response.Status = "200"
	response.Message = "search success"
  response.Task = *task
  mMonitorResultMap[mapKey] = response

  return true
}

func SubmitRawMonitorTask(keyword, registration_id, target_url string) {
  task := &MonitorTask{}
  task.Url = target_url
  task.Keyword = keyword
  task.RegistrationId = registration_id
  SubmitMonitorTask(task)
}

func GetMonitorResultByKey(mapkey string) MonitorResponse {
  resp, ok := mMonitorResultMap[mapkey]
  if ok {
    // delete(mMonitorResultMap, mapkey)
    return resp
  } else {
    var response MonitorResponse
    response.Status = "400"
  	response.Message = "has not the map key"
    return response
  }
}

func GetMonitorTaskList() []MonitorTask {
  taskList := make([]MonitorTask, 0)
  for _, item := range mTaskList {
    taskList = append(taskList, *item)
  }
  return taskList
}

func DeleteMonitorTaskByIndex(index int) {
  tmpList := make([]*MonitorTask, 0)
  for i, item := range mTaskList {
    if i == index {
      continue
    }
    tmpList = append(tmpList, item)
  }
  mTaskList = tmpList
}

func DeleteMonitorTask(task *MonitorTask) {
  mutex.Lock()
  newList := make([]*MonitorTask, 0)
  for _, item := range mTaskList {
    if task.IsEqual(item) {
      continue
    }
    newList = append(newList, item)
  }
  mTaskList = newList
  mutex.Unlock()
}
