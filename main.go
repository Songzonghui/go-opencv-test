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

var base = []string{"data/1/3.png", "data/2/14.png", "data/3/20.png", "data/4/1.png", "data/5/7.png", "data/12/1.png"}
var tests = [][]string{
	{"data/1/2.png", "data/1/4.png", "data/1/5.png", "data/1/6.png", "data/1/7.png", "data/1/8.png"},
	{"data/2/9.png", "data/2/10.png", "data/2/11.png", "data/2/12.png", "data/2/13.png", "data/2/15.png"},
	{"data/3/16.png", "data/3/17.png", "data/3/18.png", "data/3/19.png", "data/3/21.png", "data/3/22.png"},
	{"data/4/2.png", "data/4/3.png", "data/4/4.png", "data/4/5.png", "data/4/6.png", "data/4/7.png"},
	{"data/5/1.png", "data/5/2.png", "data/5/3.png", "data/5/4.png", "data/5/5.png", "data/5/6.png"},
	{"data/12/2.png", "data/12/3.png", "data/12/4.png", "data/12/5.png", "data/12/6.png", "data/12/7.png"},
}

func main() {
	check1()
	check2()
}

// 旋转角度相似度
func check1() {
	cv := &CV{Width: 114, Height: 114}
	for _, k := range base {
		func() {
			// 基准图
			tpl := cv.AnalyseAnimal(k, true)
			defer tpl.Close()

			for _, vs := range tests {
				for _, v := range vs {
					func() {
						// 测试图
						test := cv.AnalyseAnimal(v, true)
						defer test.Close()

						x := tpl.Clone()
						defer x.Close()
						rate, percent, out := cv.Check(x, test, 5, 10)
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
							gocv.IMWrite("debug1.png", out)
							// GUI显示
							cv.Show(out)
							return
						}
					}()
				}
			}
			fmt.Println()
		}()
	}
}

// 图片相似度
func check2() {
	cv2 := &CV{Width: 114, Height: 114}
	for _, k := range base {
		func() {
			// 基准图
			tpl := cv2.AnalyseAnimal(k, false)
			defer tpl.Close()

			for _, vs := range tests {
				for _, v := range vs {
					func() {
						// 测试图
						test := cv2.AnalyseAnimal(v, false)
						defer test.Close()

						x := tpl.Clone()
						defer x.Close()
						percent, out := cv2.Check2(x, test)
						defer out.Close()

						sign := "√"
						if filepath.Dir(k) == filepath.Dir(v) && percent < 0.1 {
							sign = "<== x"
						} else if percent < 0.1 {
							sign = "x"
						}
						fmt.Printf("[% 15s - % 15s] Percent:%0.3f %s\n", k, v, percent, sign)

						// 分值过低的结果
						if false {
							// 输出到文件查看
							gocv.IMWrite("debug2.png", out)
							// GUI显示
							cv2.Show(out)
							return
						}
					}()
				}
			}
			fmt.Println()
		}()
	}
}
