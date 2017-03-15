package handlefunc

import (
	"net/http"
	"time"
	"strconv"
	"encoding/json"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"io/ioutil"
	. "qiniupkg.com/api.v7/conf"
	"qiniupkg.com/api.v7/kodo"
)

type MyQiniuResponse struct {
	Status string `json:"status"`
	Message string `json:"message"`
	AccessKey string `json:"access_key,omitempty"`
	Token string `json:"token,omitempty"`
}

func initQiniu() {
	urlFuncMap["qiniu/list"] = getListToken
	urlFuncMap["qiniu/upload"] = getUploadToken
	urlFuncMap["qiniu/download"] = getDownloadToken
	urlFuncMap["qiniu/stat"] = statFile
	urlFuncMap["qiniu/delete"] = deleteFile

  ACCESS_KEY = "4oK4Tp4L4zVVwZl6Vk_d2C5O1wC08hfXRi9bAu-Q"
  SECRET_KEY = "sEhV0aPeFD57hNc4MzJIMmE39VEtxTL2K87TTOOB"
	kodo.SetMac(ACCESS_KEY, SECRET_KEY)
}

func writeQiniuResult(w http.ResponseWriter, r *http.Request, token string, err error) {
	var response MyQiniuResponse
  if err == nil {
  	response.Status = "200"
  	response.Message = "success"
  	response.Token = token
  	response.AccessKey = ACCESS_KEY
  } else {
    response.Status = "500"
  	response.Message = err.Error()
  }
  respBytes, _ := json.Marshal(response)
	w.Write(respBytes)
}

func getListToken(w http.ResponseWriter, r *http.Request) {
	bucket, ok := parseKeyword(r, "bucket")
	if !ok {
		return
	}
	host := "http://rsf.qbox.me"
	path := "/list?"
	body := "bucket=" + bucket
	accessToken := getManageToken(path + body + "\n")

	req, err := http.NewRequest("POST", host+path+body, nil)
	if err != nil {
		err.Error()
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "QBox "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
    writeQiniuResult(w, r, "", err)
		return
	}
	respB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
    writeQiniuResult(w, r, "", err)
		return
	}
	respStr := string(respB)

	writeQiniuResult(w, r, respStr, nil)
}

func getManageToken(signingStr string) string {
	mac := hmac.New(sha1.New, []byte(SECRET_KEY))
	mac.Write([]byte(signingStr))
	encodedSign := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	accessToken := ACCESS_KEY + ":" + encodedSign
	return accessToken
}

func getUploadToken(w http.ResponseWriter, r *http.Request) {
	bucket, ok := parseKeyword(r, "bucket")
	if !ok {
		return
	}
	key, ok := parseKeyword(r, "key")
	if !ok {
		return
	}
	policy := &kodo.PutPolicy{
		Scope:      bucket + ":" + key,
		Expires:    uint32(3600 + time.Now().Unix()),
		ReturnBody: "{\"w\":$(imageInfo.width),\"h\":$(imageInfo.height)}",
	}
	b, err := json.Marshal(policy)
	if err != nil {
		err.Error()
		return
	}
	encodedPutPolicy := base64.URLEncoding.EncodeToString(b)
	uptoken := ACCESS_KEY + ":" + getEncodedSign([]byte(encodedPutPolicy)) + ":" + encodedPutPolicy
	writeQiniuResult(w, r, uptoken, nil)
}

func getDownloadToken(w http.ResponseWriter, r *http.Request) {
	encodedURL, ok := parseKeyword(r, "url")
	if !ok {
		return
	}
	b, err := base64.URLEncoding.DecodeString(encodedURL)
	if err != nil {
		err.Error()
		return
	}
	realURL := string(b) + "?e=" + strconv.Itoa(int(time.Now().Unix())+3600)
	downloadtoken := ACCESS_KEY + ":" + getEncodedSign([]byte(realURL))
	downloadUrl := realURL + "&token=" + downloadtoken
	writeQiniuResult(w, r, base64.URLEncoding.EncodeToString([]byte(downloadUrl)), nil)
}

func statFile(w http.ResponseWriter, r *http.Request) {
	bucket, ok := parseKeyword(r, "bucket")
	if !ok {
		return
	}
	key, ok := parseKeyword(r, "key")
	if !ok {
		return
	}
	entry := bucket + ":" + key
	encodedEntryURI := base64.URLEncoding.EncodeToString([]byte(entry))

	actionFile(w, r, "/stat/", encodedEntryURI)
}

func deleteFile(w http.ResponseWriter, r *http.Request) {
	bucket, ok := parseKeyword(r, "bucket")
	if !ok {
		return
	}
	key, ok := parseKeyword(r, "key")
	if !ok {
		return
	}
	entry := bucket + ":" + key
	encodedEntryURI := base64.URLEncoding.EncodeToString([]byte(entry))

	actionFile(w, r, "/delete/", encodedEntryURI)
}

func actionFile(w http.ResponseWriter, r *http.Request, action, encodedEntryURI string) {
	host := "http://rs.qiniu.com"
	path := action + encodedEntryURI
	encodedSign := signString(path + "\n")
	accessToken := ACCESS_KEY + ":" + encodedSign

	req, err := http.NewRequest("GET", host+path, nil)
	if err != nil {
		err.Error()
		return
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authorization", "QBox "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
    writeQiniuResult(w, r, "", err)
		return
	}
	respB, _ := ioutil.ReadAll(resp.Body)
	writeQiniuResult(w, r, string(respB), nil)
}

func signString(data string) string {
	return getEncodedSign([]byte(data))
}

func getEncodedSign(b []byte) string {
	mac := hmac.New(sha1.New, []byte(SECRET_KEY))
	mac.Write([]byte(b))
	sign := mac.Sum(nil)
	encodedSign := base64.URLEncoding.EncodeToString(sign)
	return encodedSign
}
