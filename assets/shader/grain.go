//go:build ignore

package main

var ScreenSize vec2
var Tick float
var Strength float

func Fragment(position vec4, texCoord vec2, _ vec4) vec4 {
	uv := position.xy / ScreenSize.xy
	noise := float(fract(sin(dot(uv, vec2(12.9898, 78.233)*2.0)) * (43758.5453 + Tick)))
	final := noise * Strength
	return imageSrc0At(texCoord) - vec4(final, final, final, 0)
}
