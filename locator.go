package gaode_tile_locator

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"image/draw"
	"io/ioutil"
	"math"
	"net/http"
	"os"
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

func (t *Tile) tileXY2LongLat() {
	n := math.Pi - 2.0*math.Pi*float64(t.tileY)/math.Exp2(float64(t.level))
	t.lat = 180.0 / math.Pi * math.Atan(0.5*(math.Exp(n)-math.Exp(-n)))
	t.long = float64(t.tileX)/math.Exp2(float64(t.level))*360.0 - 180.0
	return
}
func math_sinh(x float64)(result float64){
	result = (math.Exp(x) - math.Exp(-x)) / 2
	return
}
func (t *Tile) pixelXY2LongLat()  {
	pixelXToTileAddition := float64(t.pixelX) / 256.0
	longitude := (float64(t.tileX) + pixelXToTileAddition) /  math.Exp2(float64(t.level)) * 360 - 180
	t.long = longitude
	pixelYToTileAddition:=float64(t.pixelY) /256.0
	latitude := math.Atan(math_sinh(math.Pi * (1 - 2 * (float64(t.tileY) + pixelYToTileAddition) / math.Exp2(float64(t.level))))) * 180.0 / math.Pi
	t.lat = latitude
	return
}

func ImageHandler(x float64,y float64, z int,sysbolSize int) (result *image.RGBA){
	tile:=Tile{x,y,0,0,z,0,0}
	tile.longLat2TileXY()
	fmt.Println(tile)
	tile.longLat2PixelXY()
	fmt.Println(tile)
	tile2:=Tile{0,0,218453,108154 ,18,85,29}
	tile2.tileXY2LongLat()
	fmt.Println(tile2)
	tile2.pixelXY2LongLat()
	fmt.Println(tile2)
	//url := "http://webrd02.is.autonavi.com/appmaptile?lang=zh_cn&size=1&scale=1&style=8&z=18&x=218453&y=108154"
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
	//Rect一定要从0，0开始
	//rec :=image.Rect(0,0,10,10)
	rec :=image.Rect(p1.X, p1.Y, p2.X, p2.Y)
	logo := image.NewRGBA(rec)
	blue := color.RGBA{255, 0, 255, 255}
	draw.Draw(logo, logo.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)
	result2, _ := os.Create("logo.png")
	png.Encode(result2,logo)
	m := image.NewRGBA(base_image.Bounds())
	draw.Draw(m, base_image.Bounds(), base_image, image.ZP, draw.Src)
	draw.Draw(m, logo.Bounds().Add(offset), logo, image.ZP, draw.Over)
	result=m
	return result
}