package main

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/valyala/fasthttp"
)

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/getParams":
		getParamsHandler(ctx)
	case "/postJson":
		postJsonHandler(ctx)
	case "/postForm":
		postFormHandler(ctx)
	case "/upload":
		uploadHandler(ctx)
	case "/uploadMulti":
		uploadMultiHandler(ctx)
	case "/bar":
		barHandler(ctx)
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}

func getParamsHandler(ctx *fasthttp.RequestCtx) {
	// 获取Get请求中的查询字符串参数
	// http://localhost:8081/foo?name=zhangsan
	values := ctx.QueryArgs()
	name := values.Peek("name")
	// http://localhost:8081/foo?name=zhangsan&name=lisi
	names := values.PeekMulti("name")
	fmt.Println(string(names[0]), " ", string(names[1]))

	// 获取Get请求中的请求头中的参数
	name2 := ctx.Request.Header.Peek("name")
	fmt.Println(string(name2))
	fmt.Fprintf(ctx, "Hello, foo! Method: %s, Name: %s", ctx.Request.Header.Method(), name)
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func postJsonHandler(ctx *fasthttp.RequestCtx) {
	// 获取Post请求中的Body-Raw-Json参数
	p := Person{}
	err := json.Unmarshal(ctx.PostBody(), &p)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(p.Name, " ", p.Age)
}

func postFormHandler(ctx *fasthttp.RequestCtx) {
	// 获取Post的表单数据
	name := ctx.FormValue("name")
	age := ctx.FormValue("age")
	fmt.Println(string(name), " ", string(age))

	// 获取多个表单数据
	form, err := ctx.MultipartForm()
	if err != nil {
		fmt.Println(err)
		return
	}
	name2 := form.Value["name"]
	age2 := form.Value["age"]
	fmt.Println(name2, " ", age2)
}

func uploadHandler(ctx *fasthttp.RequestCtx) {
	header, err := ctx.FormFile("file1")
	if err != nil {
		fmt.Println(err)
		return
	}

	fasthttp.SaveMultipartFile(header, fmt.Sprintf("./tmp/%s", header.Filename))
}

func uploadMultiHandler(ctx *fasthttp.RequestCtx) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Error("上传文件时发生错误", fasthttp.StatusInternalServerError)
		return
	}

	for _, v := range form.File {
		for _, header := range v {
			fasthttp.SaveMultipartFile(header, fmt.Sprintf("./tmp/%s", header.Filename))
		}
	}
}

func barHandler(ctx *fasthttp.RequestCtx) {
	// set some headers and status code first
	ctx.SetContentType("foo/bar")
	ctx.SetStatusCode(fasthttp.StatusOK)

	// then write the first part of body
	fmt.Fprintf(ctx, "this is the first part of body\n")

	// then set more headers
	ctx.Response.Header.Set("Foo-Bar", "baz")

	// then write more body
	fmt.Fprintf(ctx, "this is the second part of body\n")

	// then override already written body
	ctx.SetBody([]byte("this is completely new body contents"))

	// then update status code
	ctx.SetStatusCode(fasthttp.StatusNotFound)

	// basically, anything may be updated many times before
	// returning from RequestHandler.
	//
	// Unlike net/http fasthttp doesn't put response to the wire until
	// returning from RequestHandler.
}

func demo1() {
	httpServer := &fasthttp.Server{
		Handler: fastHTTPHandler,
		Name:    "kingkong",
	}

	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = httpServer.Serve(ln)
	if err != nil {
		fmt.Println(err)
	}
}
