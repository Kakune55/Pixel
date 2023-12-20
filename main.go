package main

import (
	"Pixel/database"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
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

	dirPath := "./data/img"

	// 使用 os.Stat 检查目录是否存在
	_, err := os.Stat(dirPath)

	if os.IsNotExist(err) {
		// 目录不存在，可以调用 os.Mkdir 创建
		err := os.MkdirAll(dirPath, 0755)
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

	database.Initdb() //初始化数据库
	fmt.Println("数据库初始化完成")
}


func main() {
	http.HandleFunc("/info", showimg)
	http.HandleFunc("/info/list", showlist)//
	http.HandleFunc("/upload", upload)//上传图片
    http.HandleFunc("/img/",downloadHandler)//图片接口
	http.HandleFunc("/img/mini",displayThumbnailHandler)//缩略图接口
	http.HandleFunc("/idlist",arrayHandler)//获取现有图片id
	http.HandleFunc("/img/del",deleteImagesHandler)//删除相应图片
	http.HandleFunc("/login",login)//登录页
	fmt.Println("Web服务器已启动")      
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func showimg(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("Web/info.html")
	t.Execute(w, "Hello")
}

func showlist(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("login")
	if cookie == nil{ //未授权禁止访问
		w.WriteHeader(401)
		w.Write([]byte(`<html><a href="/login">验证失败 点此登录</a><html>`))
		return
	}

	t, _ := template.ParseFiles("Web/list.html")
	t.Execute(w, "Hello")
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
		var linkid = RandomString(10)
		database.NewFile(linkid,md5sum,ext)
		w.Write([]byte(linkid))
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

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// 获取请求参数，例如文件名
	filename := r.FormValue("id")
	if filename == "" {
		http.Error(w, "未提供文件名", http.StatusBadRequest)
		return
	}

	// 拼接文件路径，确保路径安全性
	filePath := filepath.Join("./data/img", database.GetFileName(filename))

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "文件未找到", http.StatusNotFound)
		return
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "无法获取文件信息", http.StatusInternalServerError)
		return
	}

	// 设置响应头，告诉浏览器直接显示文件
	w.Header().Set("Content-Disposition", "inline; filename="+filename)
	w.Header().Set("Content-Type", "image/jpeg") // 适用于 JPEG 图片，根据实际文件类型设置

	// 将文件内容拷贝到响应体中
	http.ServeContent(w, r, filename, fileInfo.ModTime(), file)
}

// 辅助函数，获取文件大小
func fileSize(file *os.File) int64 {
	stat, err := file.Stat()
	if err != nil {
		return 0
	}
	return stat.Size()
}

func generateThumbnail(filePath string, width, height int) (*image.NRGBA, error) {
	// 打开图像文件
	img, err := imaging.Open(filePath)
	if err != nil {
		return nil, err
	}

	// 调整图像大小以创建缩略图
	thumbnail := imaging.Resize(img, width, height, imaging.Lanczos)

	return thumbnail, nil
}

func displayThumbnailHandler(w http.ResponseWriter, r *http.Request) {
	// 获取请求参数，例如文件名
	filename := r.FormValue("id")
	if filename == "" {
		http.Error(w, "未提供文件名", http.StatusBadRequest)
		return
	}

	// 拼接文件路径，确保路径安全性
	filePath := filepath.Join("./data/img", database.GetFileName(filename))

	// 生成缩略图
	thumbnail, err := generateThumbnail(filePath, 128, 0) // 设置缩略图的宽度和高度
	if err != nil {
		http.Error(w, "无法生成缩略图", http.StatusInternalServerError)
		return
	}

	// 设置响应头，指示这是一张图像
	w.Header().Set("Content-Type", "image/jpeg") // 使用 JPEG 格式或根据需要调整

	// 将缩略图写入响应体
	if err := imaging.Encode(w, thumbnail, imaging.JPEG); err != nil {
		http.Error(w, "无法编码图像", http.StatusInternalServerError)
		return
	}
}

func arrayHandler(w http.ResponseWriter, r *http.Request) {  //获取全部图片ID
	cookie, _ := r.Cookie("login")
	if cookie == nil{ //未授权禁止访问
		w.WriteHeader(401)
		w.Write([]byte(`<html><a href="/login">验证失败 点此登录</a><html>`))
		return
	}
	// 获取数组数据
	data, err := database.QueryId()
	if err != nil {
		http.Error(w, "数组数据获取失败", http.StatusInternalServerError)
		return
	}
	// 将数组数据转换为 JSON 格式
	responseData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "无法处理数组数据", http.StatusInternalServerError)
		return
	}

	// 设置响应头，指定内容类型为 JSON
	w.Header().Set("Content-Type", "application/json")

	// 将 JSON 数据写入响应体
	w.Write(responseData)
}

func deleteImagesHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("login")
	if cookie == nil{ //未授权禁止访问
		w.WriteHeader(401)
		w.Write([]byte(`<html><a href="/login">验证失败 点此登录</a><html>`))
		return
	}
	// 从请求参数中获取目录名
	id := r.FormValue("id")

	if id == "" {
		http.Error(w, "未提供id", http.StatusBadRequest)
		return
	}

	// 拼接目录路径，确保路径安全性
	dirPath := filepath.Join("./data/img/", database.GetFileName(id))

	// 删除目录及其所有内容
	if err := os.Remove(dirPath); err != nil {
		http.Error(w, fmt.Sprintf("无法删除 %s: %s", database.GetFileName(id), err), http.StatusInternalServerError)
		return
	}

	database.DelFile(id) //删除数据库相关记录

	// 返回成功的响应
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("成功删除"))
}

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("Web/login.html")
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w,"")
	} else {
		userlist,_:= database.QueryUser()
		fmt.Println(userlist)
		if len(userlist) == 0 {
			database.NewUser("admin", r.FormValue("passwd"))
		} else {
			if !database.CheckUserPasswd("admin", r.FormValue("passwd")) {
				http.Redirect(w, r, "/login",http.StatusFound)
				fmt.Println("密码错误")
				return
			}
		}
		cookie := http.Cookie{Name: "login", Value: "yes"}
		http.SetCookie(w, &cookie)
		fmt.Println("密码正确")
		http.Redirect(w, r, "/info/list",http.StatusFound)
	}
}