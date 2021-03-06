package imaging

import (
	"image"
	"image/color"
)

// Clone returns a copy of the given image.
func Clone(img image.Image) *image.NRGBA {
	dstBounds := img.Bounds().Sub(img.Bounds().Min)
	dst := image.NewNRGBA(dstBounds)

	switch src := img.(type) {
	case *image.NRGBA:
		copyNRGBA(dst, src)
	case *image.NRGBA64:
		copyNRGBA64(dst, src)
	case *image.RGBA:
		copyRGBA(dst, src)
	case *image.RGBA64:
		copyRGBA64(dst, src)
	case *image.Gray:
		copyGray(dst, src)
	case *image.Gray16:
		copyGray16(dst, src)
	case *image.YCbCr:
		copyYCbCr(dst, src)
	case *image.Paletted:
		copyPaletted(dst, src)
	default:
		copyImage(dst, src)
	}

	return dst
}

func copyNRGBA(dst *image.NRGBA, src *image.NRGBA) {
	srcMinX := src.Rect.Min.X
	srcMinY := src.Rect.Min.Y
	dstW := dst.Rect.Dx()
	dstH := dst.Rect.Dy()
	rowSize := dstW * 4
	parallel(dstH, func(partStart, partEnd int) {
		for dstY := partStart; dstY < partEnd; dstY++ {
			di := dst.PixOffset(0, dstY)
			si := src.PixOffset(srcMinX, srcMinY+dstY)
			copy(dst.Pix[di:di+rowSize], src.Pix[si:si+rowSize])
		}
	})
}

func copyNRGBA64(dst *image.NRGBA, src *image.NRGBA64) {
	srcMinX := src.Rect.Min.X
	srcMinY := src.Rect.Min.Y
	dstW := dst.Rect.Dx()
	dstH := dst.Rect.Dy()
	parallel(dstH, func(partStart, partEnd int) {
		for dstY := partStart; dstY < partEnd; dstY++ {
			di := dst.PixOffset(0, dstY)
			si := src.PixOffset(srcMinX, srcMinY+dstY)
			for dstX := 0; dstX < dstW; dstX++ {
				dst.Pix[di+0] = src.Pix[si+0]
				dst.Pix[di+1] = src.Pix[si+2]
				dst.Pix[di+2] = src.Pix[si+4]
				dst.Pix[di+3] = src.Pix[si+6]
				di += 4
				si += 8
			}
		}
	})
}

func copyRGBA(dst *image.NRGBA, src *image.RGBA) {
	srcMinX := src.Rect.Min.X
	srcMinY := src.Rect.Min.Y
	dstW := dst.Rect.Dx()
	dstH := dst.Rect.Dy()
	parallel(dstH, func(partStart, partEnd int) {
		for dstY := partStart; dstY < partEnd; dstY++ {
			di := dst.PixOffset(0, dstY)
			si := src.PixOffset(srcMinX, srcMinY+dstY)
			for dstX := 0; dstX < dstW; dstX++ {
				a := src.Pix[si+3]
				dst.Pix[di+3] = a

				switch a {
				case 0:
					dst.Pix[di+0] = 0
					dst.Pix[di+1] = 0
					dst.Pix[di+2] = 0
				case 0xff:
					dst.Pix[di+0] = src.Pix[si+0]
					dst.Pix[di+1] = src.Pix[si+1]
					dst.Pix[di+2] = src.Pix[si+2]
				default:
					var tmp uint16
					tmp = uint16(src.Pix[si+0]) * 0xff / uint16(a)
					dst.Pix[di+0] = uint8(tmp)
					tmp = uint16(src.Pix[si+1]) * 0xff / uint16(a)
					dst.Pix[di+1] = uint8(tmp)
					tmp = uint16(src.Pix[si+2]) * 0xff / uint16(a)
					dst.Pix[di+2] = uint8(tmp)
				}

				di += 4
				si += 4
			}
		}
	})
}

func copyRGBA64(dst *image.NRGBA, src *image.RGBA64) {
	srcMinX := src.Rect.Min.X
	srcMinY := src.Rect.Min.Y
	dstW := dst.Rect.Dx()
	dstH := dst.Rect.Dy()
	parallel(dstH, func(partStart, partEnd int) {
		for dstY := partStart; dstY < partEnd; dstY++ {
			di := dst.PixOffset(0, dstY)
			si := src.PixOffset(srcMinX, srcMinY+dstY)
			for dstX := 0; dstX < dstW; dstX++ {
				a := src.Pix[si+6]
				dst.Pix[di+3] = a

				switch a {
				case 0:
					dst.Pix[di+0] = 0
					dst.Pix[di+1] = 0
					dst.Pix[di+2] = 0
				case 0xff:
					dst.Pix[di+0] = src.Pix[si+0]
					dst.Pix[di+1] = src.Pix[si+2]
					dst.Pix[di+2] = src.Pix[si+4]
				default:
					var tmp uint16
					tmp = uint16(src.Pix[si+0]) * 0xff / uint16(a)
					dst.Pix[di+0] = uint8(tmp)
					tmp = uint16(src.Pix[si+2]) * 0xff / uint16(a)
					dst.Pix[di+1] = uint8(tmp)
					tmp = uint16(src.Pix[si+4]) * 0xff / uint16(a)
					dst.Pix[di+2] = uint8(tmp)
				}

				di += 4
				si += 8
			}
		}
	})
}

func copyGray(dst *image.NRGBA, src *image.Gray) {
	srcMinX := src.Rect.Min.X
	srcMinY := src.Rect.Min.Y
	dstW := dst.Rect.Dx()
	dstH := dst.Rect.Dy()
	parallel(dstH, func(partStart, partEnd int) {
		for dstY := partStart; dstY < partEnd; dstY++ {
			di := dst.PixOffset(0, dstY)
			si := src.PixOffset(srcMinX, srcMinY+dstY)
			for dstX := 0; dstX < dstW; dstX++ {
				c := src.Pix[si]
				dst.Pix[di+0] = c
				dst.Pix[di+1] = c
				dst.Pix[di+2] = c
				dst.Pix[di+3] = 0xff
				di += 4
				si++
			}
		}
	})
}

func copyGray16(dst *image.NRGBA, src *image.Gray16) {
	srcMinX := src.Rect.Min.X
	srcMinY := src.Rect.Min.Y
	dstW := dst.Rect.Dx()
	dstH := dst.Rect.Dy()
	parallel(dstH, func(partStart, partEnd int) {
		for dstY := partStart; dstY < partEnd; dstY++ {
			di := dst.PixOffset(0, dstY)
			si := src.PixOffset(srcMinX, srcMinY+dstY)
			for dstX := 0; dstX < dstW; dstX++ {
				c := src.Pix[si]
				dst.Pix[di+0] = c
				dst.Pix[di+1] = c
				dst.Pix[di+2] = c
				dst.Pix[di+3] = 0xff
				di += 4
				si += 2
			}
		}
	})
}

func copyYCbCr(dst *image.NRGBA, src *image.YCbCr) {
	srcMinX := src.Rect.Min.X
	srcMinY := src.Rect.Min.Y
	dstW := dst.Rect.Dx()
	dstH := dst.Rect.Dy()
	parallel(dstH, func(partStart, partEnd int) {
		for dstY := partStart; dstY < partEnd; dstY++ {
			di := dst.PixOffset(0, dstY)
			for dstX := 0; dstX < dstW; dstX++ {
				srcX := srcMinX + dstX
				srcY := srcMinY + dstY

				siy := (srcY-src.Rect.Min.Y)*src.YStride + (srcX - src.Rect.Min.X)

				var sic int
				switch src.SubsampleRatio {
				case image.YCbCrSubsampleRatio444:
					sic = (srcY-src.Rect.Min.Y)*src.CStride + (srcX - src.Rect.Min.X)
				case image.YCbCrSubsampleRatio422:
					sic = (srcY-src.Rect.Min.Y)*src.CStride + (srcX/2 - src.Rect.Min.X/2)
				case image.YCbCrSubsampleRatio420:
					sic = (srcY/2-src.Rect.Min.Y/2)*src.CStride + (srcX/2 - src.Rect.Min.X/2)
				case image.YCbCrSubsampleRatio440:
					sic = (srcY/2-src.Rect.Min.Y/2)*src.CStride + (srcX - src.Rect.Min.X)
				default:
					sic = src.COffset(srcX, srcY)
				}

				yy1 := int32(src.Y[siy]) * 0x10101
				cb1 := int32(src.Cb[sic]) - 128
				cr1 := int32(src.Cr[sic]) - 128

				r := yy1 + 91881*cr1
				if uint32(r)&0xff000000 == 0 {
					r >>= 16
				} else {
					r = ^(r >> 31)
				}

				g := yy1 - 22554*cb1 - 46802*cr1
				if uint32(g)&0xff000000 == 0 {
					g >>= 16
				} else {
					g = ^(g >> 31)
				}

				b := yy1 + 116130*cb1
				if uint32(b)&0xff000000 == 0 {
					b >>= 16
				} else {
					b = ^(b >> 31)
				}

				dst.Pix[di+0] = uint8(r)
				dst.Pix[di+1] = uint8(g)
				dst.Pix[di+2] = uint8(b)
				dst.Pix[di+3] = 0xff

				di += 4
			}
		}
	})
}

func copyPaletted(dst *image.NRGBA, src *image.Paletted) {
	srcMinX := src.Rect.Min.X
	srcMinY := src.Rect.Min.Y
	dstW := dst.Rect.Dx()
	dstH := dst.Rect.Dy()
	plen := len(src.Palette)
	pnew := make([]color.NRGBA, plen)
	for i := 0; i < plen; i++ {
		pnew[i] = color.NRGBAModel.Convert(src.Palette[i]).(color.NRGBA)
	}
	parallel(dstH, func(partStart, partEnd int) {
		for dstY := partStart; dstY < partEnd; dstY++ {
			di := dst.PixOffset(0, dstY)
			si := src.PixOffset(srcMinX, srcMinY+dstY)
			for dstX := 0; dstX < dstW; dstX++ {
				c := pnew[src.Pix[si]]
				dst.Pix[di+0] = c.R
				dst.Pix[di+1] = c.G
				dst.Pix[di+2] = c.B
				dst.Pix[di+3] = c.A
				di += 4
				si++
			}
		}
	})
}

func copyImage(dst *image.NRGBA, src image.Image) {
	srcMinX := src.Bounds().Min.X
	srcMinY := src.Bounds().Min.Y
	dstW := dst.Bounds().Dx()
	dstH := dst.Bounds().Dy()
	parallel(dstH, func(partStart, partEnd int) {
		for dstY := partStart; dstY < partEnd; dstY++ {
			di := dst.PixOffset(0, dstY)
			for dstX := 0; dstX < dstW; dstX++ {
				c := color.NRGBAModel.Convert(src.At(srcMinX+dstX, srcMinY+dstY)).(color.NRGBA)
				dst.Pix[di+0] = c.R
				dst.Pix[di+1] = c.G
				dst.Pix[di+2] = c.B
				dst.Pix[di+3] = c.A
				di += 4
			}
		}
	})
}

// toNRGBA converts any image type to *image.NRGBA with min-point at (0, 0).
func toNRGBA(img image.Image) *image.NRGBA {
	if img, ok := img.(*image.NRGBA); ok && img.Bounds().Min.Eq(image.ZP) {
		return img
	}
	return Clone(img)
}
