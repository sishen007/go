package main

import (
	"fmt"
	"reflect"
)

type IBaseModel interface {
	GetName() string
	GetAge() uint8
}

type Base struct {
}
type IBaseDao interface {
	GetList()
}
type BaseDao struct {
	EntityType reflect.Type
}

var IStuDaoImpl IStuDao

type IStuDao interface {
	IBaseDao
}
type StuDao struct {
	BaseDao
}

func (bd *BaseDao) GetList() {
	types := bd.EntityType
	model := reflect.New(types).Interface().(IBaseModel)
	fmt.Printf("%#v", model.GetName())
}

type Stu struct {
	Hobye string
}

func (s *Stu) GetName() string {
	return "name"
}
func (s *Stu) GetAge() uint8 {
	return 100
}

func NewStuDao() {
	if IStuDaoImpl == nil {
		v := new(StuDao)
		v.EntityType = reflect.TypeOf(new(Stu)).Elem()
		IStuDaoImpl = v
	}
}
func main() {
	NewStuDao()
	IStuDaoImpl.GetList()
}
