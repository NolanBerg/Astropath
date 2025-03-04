package processing

import (
	"image"
)

// Takes two frames and either brightens or darkens the resultant image depending on the blemding mode.
// Mutates the output frame by comparing itself to the incoming frame
func LuminanceBlendFrame(brighten bool, output *image.NRGBA, incoming image.Image) {
	bounds := output.Bounds()
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			px1 := output.At(x, y)
			px2 := incoming.At(x, y)

			r1, g1, b1, _ := px1.RGBA()
			r2, g2, b2, _ := px2.RGBA()

			bright1 := r1 + g1 + b1
			bright2 := r2 + g2 + b2

			if brighten {
				if bright1 >= bright2 {
					output.Set(x, y, px1)
				} else {
					output.Set(x, y, px2)
				}
			} else {
				if bright1 <= bright2 {
					output.Set(x, y, px1)
				} else {
					output.Set(x, y, px2)
				}
			}
		}
	}
}
