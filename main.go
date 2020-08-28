// PKG_CONFIG_PATH=/usr/local/Cellar/opencv/4.4.0_1/lib/pkgconfig/ go run .

package main

import (
	"fmt"
)

// docker build -f Dockerfile-base -t ohko/gocv-base-440 .
// docker build -f Dockerfile -t ohko/opencv_test .
// docker run --rm -it ohko/opencv_test

func main() {
	settings := map[string][]string{
		// "data/1/3.png": {"data/1/2.png", "data/1/4.png", "data/1/5.png", "data/1/6.png", "data/1/7.png", "data/1/8.png"},
		"data/1/3.png":  {"data/2/9.png", "data/2/10.png", "data/2/11.png", "data/2/12.png", "data/2/13.png", "data/2/15.png"},
		"data/2/14.png": {"data/2/9.png", "data/2/10.png", "data/2/11.png", "data/2/12.png", "data/2/13.png", "data/2/15.png"},
		"data/3/20.png": {"data/3/16.png", "data/3/17.png", "data/3/18.png", "data/3/19.png", "data/3/21.png", "data/3/22.png"},
		"data/4/1.png":  {"data/4/2.png", "data/4/3.png", "data/4/4.png", "data/4/5.png", "data/4/6.png", "data/4/7.png"},
		"data/5/7.png":  {"data/5/1.png", "data/5/2.png", "data/5/3.png", "data/5/4.png", "data/5/5.png", "data/5/6.png"},
		"data/12/1.png": {"data/12/2.png", "data/12/3.png", "data/12/4.png", "data/12/5.png", "data/12/6.png", "data/12/7.png"},
	}

	cv := &CV{Width: 114, Height: 114, Angle: 5, OutPerLine: 10}
	for sk, sv := range settings {
		// 基准图
		tpl := cv.analyseAnimal(sk)

		for _, v := range sv {
			// 测试图
			test := cv.analyseAnimal(v)

			rate, percent, out := cv.check(tpl, test)
			defer out.Close()
			// fmt.Println("out_"+filepath.Base(sk), out)
			fmt.Printf("[% 15s - % 15s] Percent:%0.3f Rate:%.0f %s\n", sk, v, percent, rate, map[bool]string{true: "<== x", false: ""}[percent < 0.9])

			// if percent < 0.8 {
			// cv.show(out)
			// return
			// }
		}
		fmt.Println()
	}
}
