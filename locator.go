package gaode_tile_locator

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"image/draw"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
)

type Tile struct {
	long     float64
	lat      float64
	tileX    int
	tileY    int
	level    int
	pixelX   int
	pixelY   int
}
func (t *Tile) longLat2TileXY()  {
	t.tileX = int(math.Floor((t.long + 180.0) / 360.0 * (math.Exp2(float64(t.level)))))
	t.tileY = int(math.Floor((1.0 - math.Log(math.Tan(t.lat*math.Pi/180.0)+1.0/math.Cos(t.lat*math.Pi/180.0))/math.Pi) / 2.0 * (math.Exp2(float64(t.level)))))
	return
}
func (t *Tile) longLat2PixelXY()  {
	a1 := (t.long + 180.0) / 360.0 * (math.Exp2(float64(t.level)))*256
	b1 := int(a1) % 256
	t.pixelX = b1
	a := math.Sin(t.lat*math.Pi/180.0)
	b := 0.5 - math.Log((a + 1.0)/(1-a))/(4*math.Pi)
	c := int(b* (math.Exp2(float64(t.level)))*256) % 256
	t.pixelY=c
	return
}

func ImageHandler(x float64,y float64, z int,sysbolSize int) (result *image.RGBA){
	tile:=Tile{x,y,0,0,z,0,0}
	tile.longLat2TileXY()
	tile.longLat2PixelXY()
	url := "http://webrd02.is.autonavi.com/appmaptile?lang=zh_cn&size=1&scale=1&style=8&z="+strconv.Itoa(z)+"&x="+strconv.Itoa(tile.tileX)+"&y="+strconv.Itoa(tile.tileY)
	resp, err := http.Get(url)
	if err != nil {
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	}
	base_image, err :=png.Decode( bytes.NewReader(body))
	p1 := image.Point{0, 0}
	p2 := image.Point{sysbolSize, sysbolSize}
	offset := image.Pt(tile.pixelX-sysbolSize/2, tile.pixelY-sysbolSize/2)
	rec :=image.Rect(p1.X, p1.Y, p2.X, p2.Y)
	logo := image.NewRGBA(rec)
	blue := color.RGBA{255, 0, 255, 255}
	draw.Draw(logo, logo.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)
	m := image.NewRGBA(base_image.Bounds())
	draw.Draw(m, base_image.Bounds(), base_image, image.ZP, draw.Src)
	draw.Draw(m, logo.Bounds().Add(offset), logo, image.ZP, draw.Over)
	result=m
	return result
}
