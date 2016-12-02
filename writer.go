// Copyright 2016 chrisfelesoid. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flif

import (
	"io"
	"image"
	"errors"
	"bytes"
	"image/color"
	"github.com/chrisfelesoid/go-flif/wrapper"
)

func Encode(w io.Writer, m image.Image) error {
	b := m.Bounds()
	if b.Dx() <= 0 || b.Dy() <= 0 {
		return errors.New("invalid format")
	}

	var buf bytes.Buffer
	var img *wrapper.FlifImage
	var bit int

	switch m.(type) {
	case *image.RGBA, *image.NRGBA:
		bit = 8
	//case *image.RGBA64, *image.NRGBA64:
	//case *image.Gray:
	//case *image.Gray16:
	default:
		bit = 16
	}

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			switch bit {
			case 8:
				c := color.NRGBAModel.Convert(m.At(x, y)).(color.NRGBA)
				buf.WriteByte(uint8(c.R))
				buf.WriteByte(uint8(c.G))
				buf.WriteByte(uint8(c.B))
				buf.WriteByte(uint8(c.A))
			default:
				c := color.NRGBA64Model.Convert(m.At(x, y)).(color.NRGBA64)
				buf.WriteByte(uint8(c.R >> 8))
				buf.WriteByte(uint8(c.R))
				buf.WriteByte(uint8(c.G >> 8))
				buf.WriteByte(uint8(c.G))
				buf.WriteByte(uint8(c.B >> 8))
				buf.WriteByte(uint8(c.B))
				buf.WriteByte(uint8(c.A >> 8))
				buf.WriteByte(uint8(c.A))
			}
		}
	}

	switch bit {
	case 8:
		img = wrapper.NewFlifImageFromRGBA(b.Dx(), b.Dy(), buf.Bytes())
	default:
		img = wrapper.NewFlifImageHDR(b.Dx(), b.Dy())
		px := buf.Bytes()
		for row := 0; row < b.Dy(); row++ {
			is := row*b.Dx()*8
			ie := is+b.Dx()*8+1
			img.WriteRowRGBA16(row, px[is:ie], 0)
		}
	}

	enc := wrapper.NewFlifEncoder()
	enc.AddImage(img)

	var ret []byte
	enc.EncodeMemory(&ret)

	io.Copy(w, bytes.NewBuffer(ret))

	return nil
}

