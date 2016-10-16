# debounce
it execute immediately  first call, but in second time, it will postpone its execution until after wait milliseconds have elapsed since the last time it was invoked.
it designed for multi process.

## Usage

debounce.Execute(waitTime int, bufferFilePath, data, func(data string, isFirst bool))

waitTime is elapsed time(milliseconds)
bufferFilePath is buffer file for data.
func is callback. if first call then isFirst is true.

```go
package main
import "github.com/takeshy/debounce"
import "fmt"
import "io"
import "os"
import  "os/exec"

func main(){
    debounce.Execute(60000, "/tmp/test", os.Args[1], func(data string, isFirst bool) {
        subject := "error"
        if ! isFirst {
          subject = "continuous error"
        }
        subProcess := exec.Command("mail", "-s", subject, "takeshy")
        stdin, err := subProcess.StdinPipe()
        if err != nil {
            fmt.Println(err)
        }
        if err = subProcess.Start(); err != nil {
            panic(err)
        }
        io.WriteString(stdin, data)
        stdin.Close()
        subProcess.Wait()
    })
}

```

## Installation

```
go get github.com/takeshy/debounce
```


## License

MIT

## Author

Takeshi Morita
