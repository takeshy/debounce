package debounce

import "testing"
import "strconv"
import "time"
import "os"

type Result struct {
	data string
	isFirst bool
}

func TestExecute(t *testing.T) {
  i := 0
  var result [3]Result
  resultCnt := 0

  for i < 10 {
	  go func(num int){
		  Execute(100, "./buffer.txt", strconv.Itoa(num), func(data string, isFirst bool){
			result[resultCnt]= Result{ data, isFirst }
			resultCnt+=1
		  })
	  }(i)
	  i += 1
  }
  defer os.Remove("./buffer.txt")
  time.Sleep(110 * time.Millisecond)
  if len(result[0].data) != 1 {
    t.Fatalf("first data len shuld be 1 but %v", len(result[0].data))
  }
  if result[0].isFirst != true {
    t.Fatalf("first data's isFirst shuld be true but %v", result[0].isFirst)
  }
  if len(result[1].data) != 9 {
    t.Fatalf("second data len shuld be 9 but %v", len(result[1].data))
  }
  if result[1].isFirst != false {
    t.Fatalf("second data's isFirst shuld be false but %v", result[1].isFirst)
  }
  go func(){
      Execute(100, "./buffer.txt", "10", func(data string, isFirst bool){
		  result[resultCnt]= Result{ data, isFirst }
		  resultCnt+=1
      })
  }()
  if result[2].data != "" {
    t.Fatalf("third data's shuld not be received until pass waittime. but received %v", result[2])
  }
  time.Sleep(110 * time.Millisecond)
  if result[2].data != "10" {
    t.Fatalf("third data's shuld  be 10. but received %v", result[2].data)
  }
  if result[2].isFirst != false {
    t.Fatalf("third data's isFirst shuld be false but %v", result[2].isFirst)
  }
}
