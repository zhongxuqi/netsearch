package wordtool

import (
  "fmt"
  "testing"
)

func Test_SplitSentence(t *testing.T) {
  tests := []string{"我们;fdsf,去玩《haha》《游戏王》\"方法\"“fdf”服务看两个",
    "我们;fds23f,去玩《ha范德萨ha》《游戏王》\"方法\"“fdf”"}
  for _, item := range tests {
    for _, word := range SplitSentence(item) {
      fmt.Println(word)
    }
  }
}
