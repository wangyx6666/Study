package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

var Page chan int

func HttpGet(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1
		return
	}

	defer resp.Body.Close() //别忘了关闭Body流

	buf := make([]byte, 1024*4)
	for {
		n, err2 := resp.Body.Read(buf)
		if n == 0 {
			fmt.Println("resp.Body.Read err2 = ", err2)
			return
		}
		result += string(buf[:n])

	}

}

func PaPage(i int) {
	//1.明确爬取范围和网址
	url := "http://tieba.baidu.com/f?kw=%E7%BB%9D%E5%9C%B0%E6%B1%82%E7%94%9F&ie=utf-8&pn=" + strconv.Itoa((i-1)*50)

	//2.开始爬取,调用http.Get()
	result, err := HttpGet(url)
	if err != nil {
		fmt.Println("HttpGet err = ", err)
		return
	}

	//按照页码写入文件
	fileName := strconv.Itoa(i) + ".html"
	f, err3 := os.Create(fileName)
	if err3 != nil {
		fmt.Println("os.Create err3 = ", err3)
		return
	}

	f.WriteString(result)

	Page <- i

	defer f.Close()
}

func DoWork(start, end int) {
	fmt.Printf("我们要爬取的页面从%d到%d\n", start, end)

	//爬取百度贴吧，分页，我们肯定要一页一页的爬，此时用一个for循环
	for i := start; i <= end; i++ {
		//建立end-start+1个协程，分别去爬
		go PaPage(i)
	}

	for i := start; i <= end; i++ {
		<-Page
	}
}

func main() {
	//首先要知道我爬哪几页
	var start, end int
	fmt.Printf("请输入我们要爬取的起始页(>=1):")
	fmt.Scan(&start)
	fmt.Printf("请输入我们要爬取的终止页(>start):")
	fmt.Scan(&end)

	DoWork(start, end)
}
