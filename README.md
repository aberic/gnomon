[![GoDoc](https://godoc.org/github.com/aberic/gnomon?status.svg)](https://godoc.org/github.com/aberic/gnomon)
[![Go Report Card](https://goreportcard.com/badge/github.com/aberic/gnomon)](https://goreportcard.com/report/github.com/aberic/gnomon)
[![GolangCI](https://golangci.com/badges/github.com/aberic/gnomon.svg)](https://golangci.com/r/github.com/aberic/gnomon)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/4f11995425294f42aec6a207b8aab367)](https://www.codacy.com/manual/aberic/gnomon?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=aberic/gnomon&amp;utm_campaign=Badge_Grade)
[![Travis (.org)](https://img.shields.io/travis/aberic/gnomon.svg?label=build)](https://www.travis-ci.org/aberic/gnomon)
[![Coveralls github](https://img.shields.io/coveralls/github/aberic/gnomon)](https://coveralls.io/github/aberic/gnomon?branch=master)

# Gnomon
通用编写go应用的公共库。

## 开发环境
* Go 1.12+
* Darwin/amd64

## 测试环境
* Go 1.11+
* Linux/x64

### 安装
``go get github.com/aberic/gnomon``

### 使用工具
```go
gnomon.Byte(). … // 字节
gnomon.Command(). … // 命令行
gnomon.Env(). … // 环境变量
gnomon.File(). … // 文件操作
gnomon.IP(). … // IP
gnomon.JWT(). … // JWT
gnomon.String(). … // 字符串
gnomon.CryptoHash(). … // Hash/散列
gnomon.CryptoRSA(). … // RSA
gnomon.CryptoECC(). … // ECC
gnomon.CryptoAES(). … // AES
gnomon.CryptoDES(). … // DES
gnomon.CA(). … // CA
gnomon.Log(). … // 日志
gnomon.Scale(). … // 算数/转换
gnomon.Time(). … // 时间
gnomon.Pool(). … // conn池
gnomon.GRPC(). … // grpc请求
gnomon.HTTPClient(). … // http请求
gnomon.SQL(). … // 数据库
```

### 使用HTTP Server
```go
func main() {
	httpServe := grope.NewHTTPServe(doFilter)
	router(httpServe)
	grope.ListenAndServe(":8888", httpServe)
}

func doFilter(ctx *grope.Context) {
	if ctx.HeaderGet("name") != "name" {
		log.Info("doFilter1", log.Field("resp", ctx.ResponseText(http.StatusForbidden, "filter name")))
	}
}

func router(hs *grope.GHttpServe) {
	// 仓库相关路由设置
	route := hs.Group("/one/test")
	route.Post("/test1", one1)
	route.Put("/test1", one2)
	route.Post("/test2/:a/:b", one2)
	route.Post("/test3/:a/:b", one3)
	route.Post("/test4/:a/:b", one4)
	route.Post("/test5/:a/:b", one5)
	route.Put("/test6/ok", one1)
	route.Put("/test6/ok/no", one6)
	route.Put("/test6/:a/:b", one6)
}


func one1(ctx *grope.Context) {
	ones := &TestOne{}
	_ = ctx.ReceiveJSON(ones)
	log.Info("one", log.Field("one", &ones),
		log.Field("url", ctx.Request().URL.String()), log.Field("paramMap", ctx.Params()))
	log.Info("one1", log.Field("resp", ctx.ResponseJSON(http.StatusOK, &TestTwo{
		Two:   "1",
		Twos:  false,
		TwoGo: 1,
	})))
}

……
```
更多详情参考：https://github.com/aberic/gnomon/blob/master/example/grope/grope_example.go

### 使用Balance
```go
func TestNewBalanceWeightRandom(t *testing.T) {
	b := NewBalance(Random)
	b.Add(1)
	b.Weight(1, 10)
	b.Add(2)
	b.Weight(2, 5)
	b.Add(3)
	b.Weight(3, 1)
	for i := 0; i < 50; i++ {
		t.Log(b.Acquire())
	}
}
```
更多详情参考：https://github.com/aberic/gnomon/blob/master/balance/balance_test.go


### 文档
参考 https://godoc.org/github.com/aberic/gnomon

<br><br>