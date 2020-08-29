package main

import (
	"fmt"
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

// CV ...
type CV struct {
	Width, Height int     // 生成最终比较图片的尺寸 default: 114
	Angle         float64 // 比较时每次旋转角度 default: 5
	OutPerLine    int     // 调试结果每行输出图片数量 default: 10
}

// Check 旋转检查两张图片相似度最高的角度和比例
// 返回：角度、相似度、调试结果
func (o *CV) Check(tpl, chk gocv.Mat) (float64, float64, gocv.Mat) {
	out := gocv.NewMat()
	defer out.Close()

	gocv.Hconcat(tpl, chk, &out)

	okRate := float64(0)
	okPercent := float32(0)
	var okMat gocv.Mat
	defer okMat.Close()

	for r := -180.0; r < 180; r += o.Angle {
		ro := o.RotationImg(chk, r)
		defer ro.Close()
		matResult := gocv.NewMat()
		mask := gocv.NewMat()
		gocv.MatchTemplate(ro, tpl, &matResult, gocv.TmCcoeffNormed, mask)
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

	fix := (out.Cols() / out.Rows()) % (o.OutPerLine * 2)
	if fix != 0 {
		append := (o.OutPerLine * 2) - fix
		emp := gocv.NewMatWithSize(tpl.Rows(), tpl.Cols(), tpl.Type())
		gocv.FillPoly(&emp, [][]image.Point{{{0, 0}, {0, tpl.Cols()}, {tpl.Rows(), tpl.Cols()}, {tpl.Rows(), 0}}}, color.RGBA{255, 255, 255, 1})
		defer emp.Close()
		for i := 0; i < append; i++ {
			gocv.Hconcat(out, emp, &out)
		}
	}

	tmp := gocv.NewMatWithSize(0, out.Rows()*o.OutPerLine, out.Type())
	defer tmp.Close()

	line := out.Cols() / (out.Rows() * o.OutPerLine)
	for i := 0; i < line; i += 2 {
		a := out.Region(image.Rect(out.Rows()*o.OutPerLine*(i+0), 0, out.Rows()*o.OutPerLine*(i+1), out.Rows()))
		b := out.Region(image.Rect(out.Rows()*o.OutPerLine*(i+1), 0, out.Rows()*o.OutPerLine*(i+2), out.Rows()))
		gocv.Vconcat(tmp, a, &tmp)
		gocv.Vconcat(tmp, b, &tmp)
	}

	gocv.PutText(&tmp, fmt.Sprintf("%.0f", okRate), image.Point{0, tpl.Rows()}, gocv.FontHersheyPlain, 1, color.RGBA{R: 255, G: 0, B: 0, A: 0}, 1)
	gocv.PutText(&tmp, fmt.Sprintf("%.2f", okPercent), image.Point{0, tpl.Rows() / 3}, gocv.FontHersheyPlain, 1, color.RGBA{R: 255, G: 0, B: 0, A: 0}, 1)
	return okRate, float64(okPercent), tmp.Clone()
}

// AnalyseAnimal 分析出主体内容
func (o *CV) AnalyseAnimal(filename string) gocv.Mat {
	img := gocv.IMRead(filename, gocv.IMReadColor)
	defer img.Close()
	animal := o.FindAnimal(img)
	defer animal.Close()

	compress := gocv.NewMatWithSize(o.Width, o.Height, gocv.MatTypeCV8U)
	gocv.Resize(animal, &compress, image.Pt(compress.Rows(), compress.Cols()), 0, 0, gocv.InterpolationCubic)

	gocv.GaussianBlur(compress, &compress, image.Pt(3, 3), 0, 0, gocv.BorderWrap)
	gocv.Dilate(compress, &compress, gocv.NewMat())
	return compress
}

// RotationImg 旋转图片
func (o *CV) RotationImg(src gocv.Mat, angle float64) gocv.Mat {
	out := gocv.NewMat()
	center := image.Point{src.Rows() / 2, src.Cols() / 2}
	M := gocv.GetRotationMatrix2D(center, angle, 1.0)
	gocv.WarpAffineWithParams(src, &out, M, image.Point{src.Rows(), src.Cols()}, gocv.InterpolationLinear, gocv.BorderConstant, color.RGBA{255, 255, 255, 0})
	return out
}

// FindAnimal 查找主体内容
func (o *CV) FindAnimal(img gocv.Mat) gocv.Mat {

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
	// gocv.GaussianBlur(imgClone, &imgClone, image.Pt(5, 5), 0, 0, gocv.BorderWrap)

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
			return imgClone.Region(rect)
			gocv.Rectangle(&img,
				gocv.BoundingRect(element),
				color.RGBA{R: 0, G: 255, B: 0, A: 100}, 1)
			// break
		}
	}

	return img
}

// Show GUI显示
func (o *CV) Show(img ...gocv.Mat) {
	var ws []*gocv.Window
	for k, v := range img {
		win := gocv.NewWindow(fmt.Sprintf("preview:%d", k))
		defer win.Close()
		win.IMShow(v)
		ws = append(ws, win)
	}
	for {
		for _, v := range ws {
			if v.WaitKey(0) > 0 {
				return
			}
			break
		}
	}
}
