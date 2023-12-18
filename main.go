package main

import (
	"Pixel/database"
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)


func init() {
	const appinfo string = `

    ______     ________      ______       __     __      __          
   /_____/\   /_______/\    /_____/\     /__/\ /__/\    /_/\         
   \:::_ \ \  \__.::._\/    \::::_\/_    \ \::\\:.\ \   \:\ \        
    \:(_) \ \    \::\ \      \:\/___/\    \_\::_\:_\/    \:\ \       
     \: ___\/    _\::\ \__    \::___\/_     _\/__\_\_/\   \:\ \____  
      \ \ \     /__\::\__/\    \:\____/\    \ \ \ \::\ \   \:\/___/\ 
       \_\/     \________\/     \_____\/     \_\/  \__\/    \_____\/ 

																	 
   `
	fmt.Println(appinfo)

	database.Initdb() //初始化数据库
	dirPath := "./data/img"

	// 使用 os.Stat 检查目录是否存在
	_, err := os.Stat(dirPath)

	if os.IsNotExist(err) {
		// 目录不存在，可以调用 os.Mkdir 创建
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			fmt.Println("无法创建目录:", err)
			return
		}
		fmt.Println("目录创建成功:", dirPath)
	} else if err == nil {
		// 目录已存在
		fmt.Println("目录已存在:", dirPath)
	} else {
		// 发生其他错误
		fmt.Println("发生错误:", err)
	}
}


func main() {

	http.HandleFunc("/upload", upload)        //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}


// 处理/upload 逻辑
func upload(w http.ResponseWriter, r *http.Request) {
    fmt.Println("method:", r.Method) // 获取请求的方法

    if r.Method == "GET" { // 前端页面渲染
        crutime := time.Now().Unix()
        h := md5.New()
        io.WriteString(h, strconv.FormatInt(crutime, 10))
        token := fmt.Sprintf("%x", h.Sum(nil))

        t, _ := template.ParseFiles("Web/upload.html")
        t.Execute(w, token)
    } else { // 后端POST接收逻辑
        r.ParseMultipartForm(32 << 20)
        file, handler, err := r.FormFile("file")
        if err != nil {
            fmt.Println(err)
            return
        }
        defer file.Close()

        // 生成文件的MD5哈希
        h := md5.New()
        if _, err := io.Copy(h, file); err != nil {
            fmt.Println(err)
            return
        }
        md5sum := fmt.Sprintf("%x", h.Sum(nil))

        // 获取文件扩展名
		fname := handler.Filename
        ext := path.Ext(fname)

        // 创建新文件，使用MD5哈希和原始扩展名
        newFilename := md5sum + ext
        f, err := os.OpenFile("./data/img/"+newFilename, os.O_WRONLY|os.O_CREATE, 0666)
        if err != nil {
            fmt.Println(err)
            return
        }
        defer f.Close()

        // 将文件内容拷贝到新文件
        _, err = file.Seek(0, 0)
        if err != nil {
            fmt.Println(err)
            return
        }
        io.Copy(f, file)

        // 存入数据库
		database.NewFile(RandomString(10),md5sum)
        
    }
}


func RandomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		bytes[i] = letters[rand.Intn(len(letters))]
	}
	return string(bytes)
}
