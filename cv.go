package main

import (
	"fmt"
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

// CV ...
type CV struct {
	Width, Height int // 生成最终比较图片的尺寸 default: 114
}

// Check 旋转检查两张图片相似度最高的角度和比例
// 参数：angle 比较时每次旋转角度 default: 5
// 参数：perLine 调试结果每行输出图片数量 default: 10
// 返回：角度、相似度、调试结果
func (o *CV) Check(tpl, chk gocv.Mat, angle float64, perLine int) (float64, float64, gocv.Mat) {
	out := gocv.NewMat()
	defer out.Close()

	gocv.Hconcat(tpl, chk, &out)

	okRate := float64(0)
	okPercent := float64(0)
	var okMat gocv.Mat
	defer okMat.Close()

	for r := -180.0; r < 180; r += angle {
		ro := o.RotationImg(chk, r)
		defer ro.Close()
		minConfidence, ro := o.Check1(tpl, ro)
		gocv.PutText(&ro, fmt.Sprintf("%.0f", r), image.Point{0, ro.Cols()}, gocv.FontHersheyPlain, 1, color.RGBA{R: 0, G: 0, B: 255, A: 0}, 1)
		gocv.PutText(&ro, fmt.Sprintf("%.2f", minConfidence), image.Point{0, ro.Cols() / 3}, gocv.FontHersheyPlain, 1, color.RGBA{R: 0, G: 0, B: 255, A: 0}, 1)
		// gocv.PutText(&ro, fmt.Sprintf("%.2f", maxConfidence), image.Point{0, ro.Cols() / 2}, gocv.FontHersheyPlain, 0.75, color.RGBA{R: 0, G: 0, B: 255, A: 0}, 1)
		if minConfidence > okPercent {
			okPercent = minConfidence
			okRate = r
			okMat = ro.Clone()
		}
		gocv.Hconcat(out, ro, &out)
	}

	fix := (out.Cols() / out.Rows()) % (perLine * 2)
	if fix != 0 {
		append := (perLine * 2) - fix
		emp := gocv.NewMatWithSize(tpl.Rows(), tpl.Cols(), tpl.Type())
		gocv.FillPoly(&emp, [][]image.Point{{{0, 0}, {0, tpl.Cols()}, {tpl.Rows(), tpl.Cols()}, {tpl.Rows(), 0}}}, color.RGBA{255, 255, 255, 1})
		defer emp.Close()
		for i := 0; i < append; i++ {
			gocv.Hconcat(out, emp, &out)
		}
	}

	tmp := gocv.NewMatWithSize(0, out.Rows()*perLine, out.Type())
	defer tmp.Close()

	line := out.Cols() / (out.Rows() * perLine)
	for i := 0; i < line; i += 2 {
		a := out.Region(image.Rect(out.Rows()*perLine*(i+0), 0, out.Rows()*perLine*(i+1), out.Rows()))
		b := out.Region(image.Rect(out.Rows()*perLine*(i+1), 0, out.Rows()*perLine*(i+2), out.Rows()))
		gocv.Vconcat(tmp, a, &tmp)
		gocv.Vconcat(tmp, b, &tmp)
	}

	gocv.PutText(&tmp, fmt.Sprintf("%.0f", okRate), image.Point{0, tpl.Rows()}, gocv.FontHersheyPlain, 1, color.RGBA{R: 255, G: 0, B: 0, A: 0}, 1)
	gocv.PutText(&tmp, fmt.Sprintf("%.2f", okPercent), image.Point{0, tpl.Rows() / 3}, gocv.FontHersheyPlain, 1, color.RGBA{R: 255, G: 0, B: 0, A: 0}, 1)
	return okRate, okPercent, tmp.Clone()
}

// Check1 检查两张图片相似度，与角度有关
// 返回：相似度、调试结果
func (o *CV) Check1(tpl, chk gocv.Mat) (float64, gocv.Mat) {
	matResult := gocv.NewMat()
	mask := gocv.NewMat()
	gocv.MatchTemplate(chk, tpl, &matResult, gocv.TmCcoeffNormed, mask)
	mask.Close()
	minConfidence, _, _, _ := gocv.MinMaxLoc(matResult)
	// minConfidence, maxConfidence, minLoc, maxLoc := gocv.MinMaxLoc(matResult)
	// fmt.Println(r, minConfidence, maxConfidence, minLoc, maxLoc)
	return float64(minConfidence), chk
}

// Check2 检查两张图片相似度，与角度无关
// 返回：相似度、调试结果
func (o *CV) Check2(tpl, chk gocv.Mat) (float64, gocv.Mat) {
	out := gocv.NewMat()
	defer out.Close()

	orb := gocv.NewORB()
	kp1, des1 := orb.DetectAndCompute(tpl, gocv.NewMat())
	kp2, des2 := orb.DetectAndCompute(chk, gocv.NewMat())
	if len(kp1) > 0 {
		gocv.DrawKeyPoints(tpl, kp1, &tpl, color.RGBA{0, 0, 255, 1}, gocv.DrawRichKeyPoints)
	}
	if len(kp2) > 0 {
		gocv.DrawKeyPoints(chk, kp2, &chk, color.RGBA{0, 0, 255, 1}, gocv.DrawRichKeyPoints)
	}
	bf := gocv.NewBFMatcherWithParams(gocv.NormHamming, false)
	matches := bf.KnnMatch(des1, des2, 2)
	var good1, good2 []gocv.KeyPoint
	for _, v := range matches {
		if v[0].Distance < 0.5*v[1].Distance {
			good1 = append(good1, kp1[v[0].QueryIdx])
			good2 = append(good2, kp2[v[0].TrainIdx])
		}
	}
	similary := 0.0
	if len(good1) > 0 {
		similary = float64(len(good1)) / float64(len(matches))
	}

	// fmt.Println(len(good1), len(matches))
	gocv.Hconcat(tpl, chk, &out)

	for k, v := range good1 {
		gocv.Line(&out,
			image.Pt(int(v.X), int(v.Y)),
			image.Pt(tpl.Cols()+int(good2[k].X), int(good2[k].Y)),
			color.RGBA{255, 0, 0, 1}, 1)
	}

	return similary, out.Clone()
}

// AnalyseAnimal 分析出主体内容
// 参数：Optimize 比较前是否再次优化图片，角度旋转比较Check1前应再优化一次，相似度比较Check2不需要再次优化
func (o *CV) AnalyseAnimal(filename string, optimize bool) gocv.Mat {
	img := gocv.IMRead(filename, gocv.IMReadColor)
	defer img.Close()
	animal := o.FindAnimal(img)
	defer animal.Close()

	compress := gocv.NewMatWithSize(o.Width, o.Height, gocv.MatTypeCV8U)
	gocv.Resize(animal, &compress, image.Pt(compress.Rows(), compress.Cols()), 0, 0, gocv.InterpolationCubic)

	if optimize {
		gocv.GaussianBlur(compress, &compress, image.Pt(3, 3), 0, 0, gocv.BorderWrap)
		gocv.Dilate(compress, &compress, gocv.NewMat())
	}
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
			return img.Region(rect)
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
