// PKG_CONFIG_PATH=/usr/local/Cellar/opencv/4.4.0_1/lib/pkgconfig/ go run .

package main

import (
	"fmt"
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

// docker build -f Dockerfile-base -t ohko/gocv-base-440 .
// docker build -f Dockerfile -t ohko/opencv_test .
// docker run --rm -it ohko/opencv_test

func main() {
	settings := map[string][]string{
		// "data/1/3.png":  {"data/1/2.png", "data/1/4.png", "data/1/5.png", "data/1/6.png", "data/1/7.png", "data/1/8.png"},
		"data/1/3.png":  {"data/2/9.png", "data/2/10.png", "data/2/11.png", "data/2/12.png", "data/2/13.png", "data/2/15.png"},
		"data/2/14.png": {"data/2/9.png", "data/2/10.png", "data/2/11.png", "data/2/12.png", "data/2/13.png", "data/2/15.png"},
		"data/3/20.png": {"data/3/16.png", "data/3/17.png", "data/3/18.png", "data/3/19.png", "data/3/21.png", "data/3/22.png"},
		"data/4/1.png":  {"data/4/2.png", "data/4/3.png", "data/4/4.png", "data/4/5.png", "data/4/6.png", "data/4/7.png"},
		"data/5/7.png":  {"data/5/1.png", "data/5/2.png", "data/5/3.png", "data/5/4.png", "data/5/5.png", "data/5/6.png"},
		"data/12/1.png": {"data/12/2.png", "data/12/3.png", "data/12/4.png", "data/12/5.png", "data/12/6.png", "data/12/7.png"},
	}

	for sk, sv := range settings {
		for _, v := range sv {
			rate, percent, out := check(sk, v)
			defer out.Close()
			fmt.Printf("[% 15s - % 15s] Percent:%0.3f Rate:%.0f %s\n", sk, v, percent, rate, map[bool]string{true: "<== x", false: ""}[percent <= 0.8])

			// if percent < 0.8 {
			// show(out)
			// return
			// }
		}
	}
}

func check(rightFile, checkFile string) (float64, float64, gocv.Mat) {
	w, h, angle := 114, 114, 5
	out := gocv.NewMat()
	defer out.Close()

	img := gocv.IMRead(rightFile, gocv.IMReadColor)
	defer img.Close()

	animal := findAnimal(img, w, h)
	defer animal.Close()

	testImg := gocv.IMRead(checkFile, gocv.IMReadColor)
	defer testImg.Close()

	test := findAnimal(testImg, w, h)
	defer test.Close()
	gocv.Hconcat(animal, test, &out)

	okRate := float64(0)
	okPercent := float32(0)
	var okMat gocv.Mat
	defer okMat.Close()

	for r := -180.0; r < 180; r += float64(angle) {
		ro := rotationImg(test, r)
		defer ro.Close()
		matResult := gocv.NewMat()
		mask := gocv.NewMat()
		gocv.MatchTemplate(ro, animal, &matResult, gocv.TmCcoeffNormed, mask)
		mask.Close()
		minConfidence, _, _, _ := gocv.MinMaxLoc(matResult)
		// minConfidence, maxConfidence, minLoc, maxLoc := gocv.MinMaxLoc(matResult)
		// fmt.Println(r, minConfidence, maxConfidence, minLoc, maxLoc)
		gocv.PutText(&ro, fmt.Sprintf("%.0f", r), image.Point{0, ro.Cols()}, gocv.FontHersheyPlain, 1, color.RGBA{R: 255, G: 0, B: 0, A: 0}, 1)
		gocv.PutText(&ro, fmt.Sprintf("%.2f", minConfidence), image.Point{0, ro.Cols() / 3}, gocv.FontHersheyPlain, 1, color.RGBA{R: 255, G: 0, B: 0, A: 0}, 1)
		// gocv.PutText(&ro, fmt.Sprintf("%.2f", maxConfidence), image.Point{0, ro.Cols() / 2}, gocv.FontHersheyPlain, 0.75, color.RGBA{R: 255, G: 0, B: 0, A: 0}, 1)
		if minConfidence > okPercent {
			okPercent = minConfidence
			okRate = r
			okMat = ro.Clone()
		}
		gocv.Hconcat(out, ro, &out)
	}

	// fmt.Println("[good] rate:", okRate, "percent:", okPercent)
	return okRate, float64(okPercent), out.Clone()
}

func rotationImg(src gocv.Mat, angle float64) gocv.Mat {
	out := gocv.NewMat()
	center := image.Point{src.Rows() / 2, src.Cols() / 2}
	M := gocv.GetRotationMatrix2D(center, angle, 1.0)
	gocv.WarpAffineWithParams(src, &out, M, image.Point{src.Rows(), src.Cols()}, gocv.InterpolationLinear, gocv.BorderConstant, color.RGBA{255, 255, 255, 0})
	return out
}

func findAnimal(img gocv.Mat, w, h int) gocv.Mat {

	imgClone := img.Clone()
	defer imgClone.Close()

	// 灰度化 CvtColor
	// grayImage := gocv.NewMat()
	// defer grayImage.Close()
	gocv.CvtColor(imgClone, &imgClone, gocv.ColorBGRToGray)

	// 二值化 Threshold
	// destImage := gocv.NewMat()
	gocv.Threshold(imgClone, &imgClone, 200, 255, gocv.ThresholdToZero)

	// 缩小图片 Resize
	// resultImage := gocv.NewMatWithSize(500, 400, gocv.MatTypeCV8U)
	// gocv.Resize(destImage, &resultImage, image.Pt(resultImage.Rows(), resultImage.Cols()), 0, 0, gocv.InterpolationCubic)

	// 膨胀 Dilate
	// gocv.Dilate(resultImage, &resultImage, gocv.NewMat())
	// gocv.Dilate(imgClone, &imgClone, gocv.NewMat())
	// gocv.Dilate(imgClone, &imgClone, gocv.NewMat())

	// 高斯模糊 GaussianBlur
	// gocv.GaussianBlur(resultImage, &resultImage, image.Pt(5, 5), 0, 0, gocv.BorderWrap)

	// 查找轮廓 FindContours
	results := gocv.FindContours(imgClone, gocv.RetrievalTree, gocv.ChainApproxSimple)

	// imageForShowing := gocv.NewMatWithSize(imgClone.Rows(), imgClone.Cols(), gocv.MatChannels4)
	var rect image.Rectangle
	for _, element := range results {
		// 绘制轮廓 DrawContours
		// gocv.DrawContours(&imgClone, results, index, color.RGBA{R: 0, G: 0, B: 255, A: 255}, 1)

		// 绘制轮廓的最小外接矩形 Rectangle
		rect = gocv.BoundingRect(element)
		w := rect.Max.X - rect.Min.X
		h := rect.Max.Y - rect.Min.Y
		if w == h && w > 100 {
			compress := gocv.NewMatWithSize(w, h, gocv.MatTypeCV8U)
			gocv.Resize(imgClone.Region(rect), &compress, image.Pt(compress.Rows(), compress.Cols()), 0, 0, gocv.InterpolationCubic)
			return compress
			gocv.Rectangle(&img,
				gocv.BoundingRect(element),
				color.RGBA{R: 0, G: 255, B: 0, A: 100}, 1)
			// break
		}
	}

	compress := gocv.NewMatWithSize(w, h, gocv.MatTypeCV8U)
	gocv.Resize(img.Region(rect), &compress, image.Pt(compress.Rows(), compress.Cols()), 0, 0, gocv.InterpolationCubic)
	return compress
}

func show(img gocv.Mat) {
	win := gocv.NewWindow("preview")
	defer win.Close()
	for {
		win.IMShow(img)
		if win.WaitKey(0) > 0 {
			return
		}
	}
}
