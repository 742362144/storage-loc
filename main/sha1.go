package main

import (
	b64 "encoding/base64"
	"fmt"
	"github.com/742362144/storage-loc/util"
	"sync"
)

func main() {
	//这是将要编解码的字符串。
	var wg sync.WaitGroup
	wg.Add(10000)
	//count := make(chan bool, 100)
	for i := 0; i < 1000; i++ {
		go func(w sync.WaitGroup) {
			data := util.GENERATEVALUE(1024 * 64)
			//Go 同时支持标准的和 URL 兼容的 base64 格式。编码需要使用 []byte 类型的参数，所以要将字符串转成此类型。
			sEnc := b64.StdEncoding.EncodeToString([]byte(data))
			fmt.Println(sEnc)
			//解码可能会返回错误，如果不确定输入信息格式是否正确，那么，你就需要进行错误检查了。
			sDec, _ := b64.StdEncoding.DecodeString(sEnc)
			fmt.Println(string(sDec))
			fmt.Println()
			//使用 URL 兼容的 base64 格式进行编解码。
			uEnc := b64.URLEncoding.EncodeToString([]byte(data))
			fmt.Println(uEnc)
			uDec, _ := b64.URLEncoding.DecodeString(uEnc)
			fmt.Println(string(uDec))
			//标准 base64 编码和 URL 兼容 base64 编码的编码字符串存在稍许不同（后缀为 + 和 -），但是两者都可以正确解码为原始字符串。
			w.Done()
		}(wg)
	}
	//go func(c chan bool) {
	//	total := 0
	//	for {
	//		select {
	//		case <-c:
	//			total++;
	//			if total {
	//
	//			}
	//		}
	//	}
	//}(count)
	wg.Wait()
}