package gnomon

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestError1(t *testing.T) {
	err1 := errors.New("new error")
	err2 := fmt.Errorf("err2: [%w]", err1)
	err3 := fmt.Errorf("err3: [%w]", err2)
	fmt.Println(err3)
}

func TestError2(t *testing.T) {
	// As 在err的链中找到与目标匹配的第一个错误，如果有则返回true，否则返回false
	fmt.Println("-------As--------")
	if _, err := os.Open("non-existing"); err != nil {
		var pathError *os.PathError
		if errors.As(err, &pathError) {
			fmt.Println("Failed at path:", pathError.Path)
		} else {
			fmt.Println(err)
		}
	}

	// New
	fmt.Println("-------New--------")
	err1 := errors.New("error1")
	err2 := errors.New("error2")

	// Is 判断两个error是否相等
	fmt.Println("-------Is--------")
	fmt.Println(errors.Is(err1, err1))
	fmt.Println(errors.Is(err1, err2))
	fmt.Println(errors.As(err1, &err2))
	fmt.Println(errors.Is(err1, errors.New("error1")))

	// Unwrap 如果传入的err对象中有%w关键字的格式化类容，则会在返回值中解析出这个原始error，多层嵌套只返回第一个，否则返回nil
	fmt.Println("-------Unwrap--------")
	e := errors.New("e")
	e1 := fmt.Errorf("e1: %w", e)
	e2 := fmt.Errorf("e2: %w", e1)
	fmt.Println(e2)
	fmt.Println(errors.Unwrap(e2))
	fmt.Println(e1)
	fmt.Println(errors.Unwrap(e1))
	Log().Error("error", Log().Err(e2))
}
