package debounce

import (
	"fmt"
	"syscall"
	"os"
	"strconv"
	"strings"
	"time"
)

const layout = "2006-01-02 15:04:05.000000000 MST"
const byteLen = 16
const metaDataLen = len(layout) + byteLen + 1

type Debounce struct {
	fd			int
	bufferPath	string
	waitTime	int
	callback	func(data string, isFirst bool)
}

func (d *Debounce) callCallback(isFirst bool, text string){
	d.clearData(isFirst)
	d.callback(text, isFirst)
}

func (d *Debounce) clearData(isTimer bool) {
	if _, err := syscall.Seek(d.fd, 0, 0); err != nil {
		panic(err)
	}
	byteNum := 0
	if isTimer{
		byteNum = -1
	}
	if _, err := syscall.Write(d.fd, []byte(time.Now().Format(layout) + fmt.Sprintf("%" + strconv.Itoa(byteLen) + "d\n", byteNum))); err != nil {
		panic(err)
	}
	if err := syscall.Ftruncate(d.fd, int64(metaDataLen)); err != nil {
		panic(err)
	}
}

func (d *Debounce) updateMetaByteCnt(byteCnt int) {
	if _, err := syscall.Seek(d.fd, int64(len(layout)), 0); err != nil {
		panic(err)
	}
	if _, err := syscall.Write(d.fd, []byte(fmt.Sprintf("%" + strconv.Itoa(byteLen) + "d\n", byteCnt))); err != nil {
		panic(err)
	}
}

func (d *Debounce) writeData(byteCnt int, text string) {
	d.updateMetaByteCnt(byteCnt + len(text))
	if _, err := syscall.Seek(d.fd, int64(metaDataLen + byteCnt), 0); err != nil {
		panic(err)
	}
	if _, err := syscall.Write(d.fd, []byte(text)); err != nil {
		panic(err)
	}
}

func (d *Debounce) openWithLock() {
	fd, err := syscall.Open(d.bufferPath, syscall.O_RDWR | syscall.O_CREAT, 0600);
	if err != nil {
		panic(err)
	}
	d.fd = fd
	if err = syscall.Flock(d.fd, syscall.LOCK_EX); err != nil {
		panic(err)
	}
}


func (d *Debounce) controlWithMeta(text string, useTimer bool) {
	d.openWithLock()

	metaData := make([]byte, metaDataLen)
	count, err := syscall.Read(d.fd, metaData)
	if err != nil {
		panic(err)
	}

	if count > 0{
		t, err := time.Parse(layout, string(metaData[0:len(layout)]))
		if err != nil {
			panic(err)
		}
		byteCnt, _ := strconv.Atoi(strings.TrimSpace(string(metaData[len(layout):])))
		millSeconds := int(time.Now().Sub(t).Seconds() * 1000)
		if(byteCnt <= 0){
			if millSeconds < d.waitTime {
				d.writeData(0, text)
				if byteCnt < 0 {
					syscall.Close(d.fd)
					return
				}
			} else {
				if(text != ""){
					d.callCallback(true, text)
				}
			}
		} else {
			if millSeconds < d.waitTime && text != ""{
			  d.writeData(byteCnt, text)
			  syscall.Close(d.fd)
			  return
			}
			data := make([]byte, byteCnt)
			_, err := syscall.Read(d.fd, data)
			if err != nil {
				panic(err)
			}
			d.callCallback(false, string(data) + text)
			syscall.Close(d.fd)
			return
		}
	} else {
		d.callCallback(true, text)
	}
	syscall.Close(d.fd)
	if !useTimer {
		return
	}
	time.Sleep(time.Duration(d.waitTime) * time.Millisecond)
	d.controlWithMeta("", false)
}

func Execute(waitTime int, bufferPath string, text string, callback func(data string, isFirst bool)) {
	if text == "" {
		os.Exit(-1)
	}
	d := new(Debounce)
	d.waitTime = waitTime
	d.bufferPath = bufferPath
	d.callback = callback
	d.controlWithMeta(text, true)
}
