package main

import (
	"encoding/json"
	"fmt"

	// "io"
	"log"
	"net"
	"os"

	"github.com/ipipdotnet/ipdb-go"
	"github.com/kataras/iris/v12"
)

var db *ipdb.City

func init() {
	fmt.Println("init")
	db1, err := ipdb.NewCity("qqwry.ipdb")
	if err != nil {
		log.Fatal(err)
	} else {
		db = db1
	}
}

func main() {
	app := iris.New()
	// tmpl注册html页面,并重载所有方法
	tmpl := iris.HTML("./template", ".html")
	//
	tmpl.Reload(true)
	// app注册tmpl
	app.RegisterView(tmpl)

	// 访问静态文件
	app.HandleDir("/static", "./static")
	app.HandleDir("/assets", ".")
	app.Get("/", home)
	app.Get("/{ip:string}", query)
	app.Listen(":5200")
}

func query(ctx iris.Context) {
	var temp []interface{}
	var sult map[string]interface{} = make(map[string]interface{}, 0)
	ip := ctx.Params().Get("ip")
	remote_ip, _, err := net.SplitHostPort(ctx.Request().RemoteAddr)
	if err != nil {
		log.Fatal(err)
	}
	if len(ip) == 0 {
		ip = remote_ip
	}
	// fmt.Println(len(ip), ip)
	fmt.Printf("query[%s]\n", ip)
	city, err := db.FindMap(ip, "CN")
	if err != nil {
		// ctx.JSON(map[string]string{"status": "0", "msg": "请检查ip地址是否正确"})
		sult["code"] = 1
		sult["msg"] = "处理错误"
		sult["count"] = 1
		sult["data"] = temp
	} else {
		city["status"] = "1"
		city["ip"] = ip
		sult["code"] = 0
		sult["msg"] = ""
		sult["count"] = 1
		temp = append(temp, city)
		sult["data"] = temp
		// ctx.JSON(city)
	}
	// 保存成json数据
	err = saveJSONToFile(sult, "city.json")
	if err != nil {
		fmt.Println(err)
	}
	// ctx.JSON(sult)
	fmt.Println(sult)
	ctx.ViewData("city", sult)
	ctx.View("index.html")
}

func home(ctx iris.Context) {
	// 空接口列表用来存储列表型字典
	var temp []interface{}
	var sult map[string]interface{} = make(map[string]interface{}, 0)
	ip, _, err := net.SplitHostPort(ctx.Request().RemoteAddr)
	if err != nil {
		log.Fatal(err)
	}
	// 指定ip
	// ip = "101.227.131.220"
	fmt.Printf("home[%s]\n", ip)
	city, err := db.FindMap(ip, "CN")
	if err != nil {
		// ctx.JSON(map[string]string{"status": "0", "msg": "请检查ip地址是否正确"})
		sult["code"] = 1
		sult["msg"] = "处理错误"
		sult["count"] = 1
		sult["data"] = temp
	} else {
		city["status"] = "1"
		city["ip"] = ip
		// layui返回的json数据需要满足code/msg/count/data,data需要返回列表格式
		sult["code"] = 0
		sult["msg"] = ""
		sult["count"] = 1
		temp = append(temp, city)
		sult["data"] = temp

	}
	// 保存成json数据
	err = saveJSONToFile(sult, "city.json")
	if err != nil {
		fmt.Println(err)
	}
	// ctx.JSON(sult)
	fmt.Println(sult)
	ctx.ViewData("city", sult)
	ctx.View("index.html")
}

// 用来保存json文件到当前目录下
func saveJSONToFile(jsonObject interface{}, filename string) error {
	// 创建文件名
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	// 函数完成时关闭
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(jsonObject)
	if err != nil {
		return err
	}
	return nil
}
