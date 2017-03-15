package bean

import (
  "errors"
  "strconv"
  "time"
)

type UrlTask struct {
  Url string
  Keywords []string
  Info DataItem
}

type UrlResult struct {
  Task UrlTask
  Eval []int
}

func (result *UrlResult) Compare(another *UrlResult) (int, error) {
  if len(result.Eval) != len(another.Eval) {
    return 0, errors.New("length is not match")
  }
  for i := 0; i < len(result.Eval); i++ {
    if result.Eval[i] > another.Eval[i] {
      return 1, nil
    } else if result.Eval[i] < another.Eval[i] {
      return -1, nil
    }
  }
  return 0, nil
}

func (result *UrlResult) ToString() (str string) {
  str = result.Task.Info.Title + ", " +result.Task.Url + ":"
  str += "["
  for _, item := range result.Eval {
    str += strconv.Itoa(item) + ","
  }
  str += "]"
  return
}

const (
  MONITOR_TASK_LIFE_DEFAULT = 25
)

type MonitorTask struct {
  Url string
  Keyword string
  Keywords []string
  RegistrationId string
  RespMapKey string
  Time time.Time `json:"-"`
  Life int
}

func (task *MonitorTask) IsEqual(another *MonitorTask) bool {
  if task.Url == another.Url &&
  task.Keyword == another.Keyword &&
  task.RegistrationId == another.RegistrationId {
    return true
  } else {
    return false
  }
}

type MonitorResponse struct {
  BaseResponse
  Url string
  Title string
  Keyword string
  Task MonitorTask
  ResultData []UrlResult
}

type MonitorListResponse struct {
  BaseResponse
  Tasks []MonitorTask
}
