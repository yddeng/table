package main

import (
	"fmt"
	"github.com/yddeng/dutil/dhttp"
	"github.com/yddeng/dutil/io"
	"io/ioutil"
)

func main() {
	resp, err := dhttp.Get("https://docs.qq.com/sheet/DQWNwRkl1dnZpT3Np?tab=BB08J2&c=A1A0A0", 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	io.WriteByte("./", "e.html", data)

}
