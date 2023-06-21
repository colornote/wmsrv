package pkg

import (
	"strconv"
	"strings"
	"time"

	"github.com/danhper/structomap"
	"golang.org/x/crypto/bcrypt"
)

func Setup() error {
	// init segmenter
	SegmentInit()

	return nil
}

func BaseRes() structomap.Serializer {
	s := structomap.New().UseSnakeCase()
	return s

}
func NewRes() structomap.Serializer {
	s := structomap.New().UseSnakeCase().Omit("DefaultModel")
	s.PickFunc(func(t interface{}) interface{} {
		return t.(time.Time).UTC()
	}, "CreatedAt", "UpdatedAt").Pick("ID")
	return s
}

// generate user password
func GeneratePasswordHash(password string) (string, error) {
	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// compare user password
func ComparePasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func StrToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// 百度爬虫判断
func IsBaiduSpider(agent string) bool {
	// 百度爬虫判断
	agent = strings.ToLower(agent)
	if strings.Contains(agent, "baiduspider") {
		// 百度爬虫
		return true
	}
	return false
}
