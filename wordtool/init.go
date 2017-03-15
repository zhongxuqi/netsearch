package wordtool

import (
  "log"
  "os"
  "encoding/csv"
  "strconv"
)

type WordObject struct {
  word string
  attr string
  times int
}

var (
  csvfilename = "data/CorpusWordPOSlist.csv"
  wordDict = make(map[string] *WordObject)
)

func init() {
  wordfile, err := os.Open(csvfilename)
  if err != nil {
    log.Fatal(err.Error())
  }
  defer wordfile.Close()

  reader := csv.NewReader(wordfile)
  reader.Read()
  for {
    record, err := reader.Read()
    if err != nil {
      break
    }
    word_times, _ := strconv.Atoi(record[3])
    wordObject := &WordObject{
      word: record[1],
      attr: record[2],
      times: word_times,
    }
    wordDict[record[1]] = wordObject
  }
}
