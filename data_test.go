package goinsight

import (
	"testing"
	"fmt"
	"unsafe"
)

func TestInsightMemString(t *testing.T) {
	a := "hahahahah"
	InsightMemString(&a)
}

func TestInsightMemArray(t *testing.T) {
	a := [4]string{"1","bb","tt","kk"}
	InsightMemArray(&a)
}

func TestInsightMemSlice(t *testing.T) {
	a := []string{"1","bb","tt","kk"}
	InsightMemSlice(a)
}


func TestInsightEmptyInterface(t *testing.T) {
	a := [4]string{"1","bb","tt","kk"}
	InsightEmptyInterface(a)
}


// Struct and interface definition just for testing non empty interface

type man struct {
	name string
	age int
	wife int
}

type woman struct {
	name string
	age int
	husband int
}

type people interface {
	showName()
}

func (i man) showName() {
	fmt.Println(i.age)
}

func (i woman) showName() {
	fmt.Println(i.age)
}

func (i woman) showHusband() {
	fmt.Println(i.husband)
}


func TestInsightNonEmptyInterface(t *testing.T) {

	xiaoming := man{
		name:"xiaoming",
		age:20,
		wife:0,
	}

	xiaohong := woman{
		name:"xiaohong",
		age:19,
		husband:2,
	}

	var p1,p2 people

	p1 = xiaoming
	p2 = xiaohong

	InsightNonEmptyInterface(unsafe.Pointer(&p1))
	InsightNonEmptyInterface(unsafe.Pointer(&p2))

}

func TestInsightMap(t *testing.T) {
	a := map[string]string{
		"name":"lishaopeng",
		"age":"1000",
		"in3":   "hhdddh",
		"in31":  "hhddh2",
		"in32":  "hhddh3",
		"in33":  "hhddh4",
		"in34":  "hhddh5",
		"in35":  "hhddh6",
		"in356": "hhddh7",

	}
	InsightMapString(a)

}