package pkg

import (
	"os"

	"github.com/huichen/sego"
)

var segmenter sego.Segmenter

func SegmentInit() {
	if os.Getenv("MODE") == "dev" {
		return
	}
	path := os.Getenv("DICT_PATH")
	segmenter.LoadDictionary(path)

}

// 返回分词后的 n 或者 v
func WordSegment(text []byte) []string {
	// 载入词典  只执行一次

	// 分词

	segments := segmenter.Segment(text)

	// 处理分词结果
	// 支持普通模式和搜索模式两种分词，见代码中SegmentsToString函数的注释。
	// fmt.Println(sego.SegmentsToString(segments, true))

	r := make([]string, 0)
	for _, v := range segments {
		t := v.Token()
		if t.Pos() == "n" || t.Pos() == "v" || t.Pos() == "t" {
			r = append(r, t.Text())
		}
	}
	return r
}
