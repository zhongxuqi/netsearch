package bean

type DataItem struct {
	Title string
	Abstract string
	Link string
	ImageUrl string
}

type SpiderObject struct {
  Url string
  Name string
  ParserName string
  Logo string
	GetDataFunc func(string) ([]DataItem, string, error) `json:"-"`
	ParseFunc func(string) ([]DataItem, string, error) `json:"-"`
}

type SpiderResponse struct {
	BaseResponse
	ResultData []DataItem `json:",omitempty"`
	NextPage string `json:",omitempty"`
}

type UnionResponse struct {
	BaseResponse
  Keyword string `json:",omitempty"`
	ResultData []UrlResult `json:",omitempty"`
  ParserNames string `json:",omitempty"`
}

type PushMessage struct {
	MapKey string
	Keyword string
	Type int
}
