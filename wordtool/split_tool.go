package wordtool

import (
  "regexp"
  "strings"
)

var (
  WordBlock = "[\u4e00-\u9fa5a-zA-Z0-9]+"
  ChineseBlock = "[\u4e00-\u9fa5]+"
  RawSubItem = []string{"[^\u4e00-\u9fa5a-zA-Z0-9]+", "[\u4e00-\u9fa5]+", "[a-zA-Z0-9]+"}
  ForceBlock = []string{"《.*?》", "\".*?\"", "“.*?”", "'.*?'"}
  BANSIGN = []string{",", ".", "，", "。", "+", " ", "\n"}
  SECTION_END = []string{",", ".", ";", "。", "，"}
)

func SplitBlock(block string) (wordList []string) {
  re := regexp.MustCompile("(" + strings.Join(RawSubItem, "|") + ")")
  for _, subBlock := range re.FindAllString(block, -1) {
    subRunes := []rune(subBlock)
    prevWord := []int{0, 0}
    indexCurr := 0
    clipSecondWordLen := 0
    for indexCurr < len(subRunes) {
      index := indexCurr
      // match word
      for i := 0; (i < 4) && (indexCurr + i < len(subRunes)); i++ {
        if _, ok := wordDict[string(subRunes[indexCurr:indexCurr + i + 1])]; ok {
          index = indexCurr + i + 1
        }
      }

      // if not match, match reversely
      if (index == indexCurr) && (indexCurr > 0) && (prevWord[1] == indexCurr) {
        backOff := 1
FindClip:
        for indexCurr - backOff > prevWord[0] {
          indexClip := indexCurr - backOff
          for clipSecondWordLen = backOff; clipSecondWordLen < 4; clipSecondWordLen++ {
            if _, ok_prev := wordDict[string(subRunes[indexClip:(indexClip + 1 + clipSecondWordLen)])]; ok_prev {
              if _, ok_next := wordDict[string(subRunes[prevWord[0]:indexClip])]; ok_next {
                index = indexClip
                break FindClip
              }
            }
          }
          backOff++
        }
      }

      // check the match result
      if index > indexCurr {
        if indexCurr > prevWord[1] {
          wordList = append(wordList, string(subRunes[prevWord[1]:indexCurr]))
        }
        prevWord = []int{indexCurr, index}

        // append matched words
        wordList = append(wordList, string(subRunes[prevWord[0]:prevWord[1]]))
        indexCurr = index
      } else if index < indexCurr {

        // append matched words
        wordList = append(wordList[:len(wordList) - 1],
        string(subRunes[prevWord[0]:index]),
        string(subRunes[index:(index + clipSecondWordLen + 1)]))
        prevWord = []int{index, index + clipSecondWordLen + 1}
        indexCurr = index + clipSecondWordLen + 1
      } else {
        indexCurr++

        // append lost matched words
        if indexCurr >= len(subRunes) {
          wordList = append(wordList, string(subRunes[prevWord[1]:indexCurr]))
        }
      }
    }
  }
  return wordList
}

func SplitSentence(sentence string) (wordList []string) {
  re := regexp.MustCompile("(" + strings.Join(ForceBlock, "|") + ")")
  lastEnd := 0
  for _, indexs := range re.FindAllStringIndex(sentence, -1) {
    if indexs[0] > lastEnd {
      for _, word := range SplitBlock(sentence[lastEnd : indexs[0]]) {
        wordList = append(wordList, word)
      }
    }
    if indexs[0] < indexs[1] {
      rs := []rune(sentence[indexs[0]:indexs[1]])
      if len(rs) > 2 {
        wordList = append(wordList, string(rs[1:len(rs) - 1]))
      }
    }
    lastEnd = indexs[1]
  }
  if lastEnd < len(sentence) {
    for _, word := range SplitBlock(sentence[lastEnd:]) {
      wordList = append(wordList, word)
    }
  }
  return
}

func SplitContent2Words(keywords string) (wordList []string) {
  re := regexp.MustCompile("[^" + strings.Join(BANSIGN, "") + "]+")
  for _, block := range re.FindAllString(keywords, -1) {
    for _, word := range SplitSentence(block) {
      wordList = append(wordList, word)
    }
  }
  return
}

func SplitArticle2Sentence(article string) (sentenceList []string) {
  re := regexp.MustCompile("[^" + strings.Join(SECTION_END, "") + "]+")
  for _, sentence := range re.FindAllString(article, -1) {
    sentenceList = append(sentenceList, sentence)
  }
  return
}
