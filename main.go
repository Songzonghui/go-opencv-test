// PKG_CONFIG_PATH=/usr/local/Cellar/opencv/4.4.0_1/lib/pkgconfig/ go run .

package main

import (
	"fmt"
	"path/filepath"

	"gocv.io/x/gocv"
)

// docker build -f Dockerfile-base -t ohko/gocv-base-440 .
// docker build -f Dockerfile -t ohko/opencv_test .
// docker run --rm -it ohko/opencv_test

func main() {
	base := []string{"data/1/3.png", "data/2/14.png", "data/3/20.png", "data/4/1.png", "data/5/7.png", "data/12/1.png"}
	tests := [][]string{
		{"data/1/2.png", "data/1/4.png", "data/1/5.png", "data/1/6.png", "data/1/7.png", "data/1/8.png"},
		{"data/2/9.png", "data/2/10.png", "data/2/11.png", "data/2/12.png", "data/2/13.png", "data/2/15.png"},
		{"data/3/16.png", "data/3/17.png", "data/3/18.png", "data/3/19.png", "data/3/21.png", "data/3/22.png"},
		{"data/4/2.png", "data/4/3.png", "data/4/4.png", "data/4/5.png", "data/4/6.png", "data/4/7.png"},
		{"data/5/1.png", "data/5/2.png", "data/5/3.png", "data/5/4.png", "data/5/5.png", "data/5/6.png"},
		{"data/12/2.png", "data/12/3.png", "data/12/4.png", "data/12/5.png", "data/12/6.png", "data/12/7.png"},
	}

	cv := &CV{Width: 114, Height: 114, Angle: 5, OutPerLine: 10}
	for _, k := range base {
		// 基准图
		tpl := cv.AnalyseAnimal(k)

		for _, vs := range tests {
			for _, v := range vs {
				// 测试图
				test := cv.AnalyseAnimal(v)

				rate, percent, out := cv.Check(tpl, test)
				defer out.Close()

				sign := "√"
				if filepath.Dir(k) == filepath.Dir(v) && percent < 0.9 {
					sign = "<== x"
				} else if percent < 0.9 {
					sign = "x"
				}
				fmt.Printf("[% 15s - % 15s] Percent:%0.3f Rate:% 4.0f %s\n", k, v, percent, rate, sign)

				// 分值过低的结果
				if percent < 0.5 {
					// 输出到文件查看
					gocv.IMWrite("test_out.png", out)
					// GUI显示
					cv.Show(out)
					return
				}
			}
		}
		fmt.Println()
	}
}
