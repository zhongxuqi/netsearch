package config

import (
  . "DesertEagleSite/bean"
)

const (
  BaseURL = "http://192.168.1.100:8089"
)

var SpiderMap = make(map[string]SpiderObject)
