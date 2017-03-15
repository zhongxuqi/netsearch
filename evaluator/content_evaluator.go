package evaluator

import (
  "strings"
  "DesertEagleSite/wordtool"
  "github.com/PuerkitoBio/goquery"
)

func EvaluateUrlByKeyWords(url string, keywords []string) (eval []int, err error) {
  resp, err := goquery.NewDocument(url)
	if err != nil {
		return []int{}, err
	}
  eval = EvaluateContentByKeyWords(resp.Find("body").Text(), keywords)
  return
}

func EvaluateContentByKeyWords(content string, keywords []string) (eval []int) {
  eval = make([]int, len(keywords))
  for i := 0; i < len(eval); i++ {
    eval[i] = 0
  }
  for _, sentence := range wordtool.SplitArticle2Sentence(content) {
    score := EvaluateSentenceByKeyWords(sentence, keywords)
    if score > 0 {
      eval[len(eval) - score]++
    }
  }
  return
}

func EvaluateSentenceByKeyWords(sentence string, keywords []string) (score int) {
  score = 0
  for _, keyword := range keywords {
    if strings.Index(strings.ToLower(sentence), strings.ToLower(keyword)) >= 0 {
      score++
    }
  }
  return
}
