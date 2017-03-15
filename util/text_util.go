package util

import (
  "os"
  "fmt"
  "bytes"
  "time"
  "strconv"
  "io/ioutil"
  "unicode/utf8"
  "net/http"
  "encoding/json"
  "golang.org/x/text/encoding/simplifiedchinese"
  "golang.org/x/text/transform"
  "github.com/saintfish/chardet"
)

func DecodeUtf8String(encodedString string) (decodedString string) {
  decodedString = ""
  if utf8.ValidString(encodedString) {
    for len(encodedString) > 0 {
      r, size := utf8.DecodeRuneInString(encodedString)
      decodedString = decodedString + string(r)
      encodedString = encodedString[size:]
    }
  } else {
    decodedString = encodedString
  }
  return
}

func Write2File(filename string, data []byte) {
  f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
  if err != nil {
    fmt.Println(err.Error())
    return
  }
  defer f.Close()
  f.Write(data)
}

func GetResp2Bytes(resp *http.Response) ([]byte) {
  buf := bytes.NewBuffer([]byte(""))
  resp.Write(buf)
  return buf.Bytes()
}

func WriteResp2File(filename string, resp *http.Response) {
  Write2File(filename, GetResp2Bytes(resp))
}

func GbkToUtf8(s []byte) ([]byte, error) {
  reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
  d, e := ioutil.ReadAll(reader)
  if e != nil {
      return nil, e
  }
  return d, nil
}

func DecodeResponse2Utf8Bytes(resp *http.Response) ([]byte, error) {
  buf := bytes.NewBuffer([]byte(""))
	resp.Write(buf)
  decoder := simplifiedchinese.GBK.NewDecoder()
  result, _ := chardet.NewHtmlDetector().DetectBest(buf.Bytes())
  if result.Charset == "GB-18030" {
    decoder = simplifiedchinese.GBK.NewDecoder()
  }
  reader := transform.NewReader(bytes.NewReader(buf.Bytes()), decoder)
  b, err := ioutil.ReadAll(reader)
  if err != nil {
			return []byte{}, err
  }
	return b, nil
}

func ConvBytes2Reader(data []byte) (*bytes.Reader) {
  return  bytes.NewReader(data)
}

func GetFormatTimeNow() string {
  t := time.Now()
  return strconv.Itoa(t.Year()) + "-" + t.Month().String() + "-" + strconv.Itoa(t.Day()) +
    "-" + strconv.Itoa(t.Hour()) + "-" + strconv.Itoa(t.Minute()) + "-" + strconv.Itoa(t.Second())
}

func ConvObject2Json(o interface{}) string {
  respBytes, _ := json.Marshal(o)
  return string(respBytes)
}
