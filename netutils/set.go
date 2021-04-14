package netutils

import (
	"log"
	"strings"
)

// Set 表示一个集合，里面所有元素不重复
type Set struct {
	data []string
}

// Add 给Set增加元素
func (set *Set) Add(target string) {
	if !set.IsExist(target) {
		set.data = append(set.data, target)
	}
}

// IsExist 判断set中是否存在
func (set *Set) IsExist(target string) bool {
	for _, str := range set.data {
		if strings.Compare(str, target) == 0 {
			return true
		}
	}
	return false
}

// GetAll 获取Set的所有元素
func (set *Set) GetAll() []string {
	return set.data
}

// Print 打印
func (set *Set) Print() {
	log.Print(set.data)
}

// AddAll 批量添加
func (set *Set) AddAll(bounds ...[]string) {
	for _, subBounds := range bounds {
		for _, bound := range subBounds {
			set.Add(bound)
		}
	}
}

// Size 返回set的大小
func (set *Set) Size() int {
	return len(set.data)
}
