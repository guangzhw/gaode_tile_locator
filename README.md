# gaode_tile_locator
使用示例：
func main() {
	http.HandleFunc("/gaode_tile_locator/", doImageHandler)
	http.ListenAndServe("localhost:800", nil)
}
func doImageHandler(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	var x, y float64
	var z, sysmbolSize int
	if vars != nil {
		x, _ = strconv.ParseFloat(vars["x"][0], 64)
		y, _ = strconv.ParseFloat(vars["y"][0], 64)
		z, _ = strconv.Atoi(vars["z"][0])
		sysmbolSize, _ = strconv.Atoi(vars["size"][0])
		m := gaode_tile_locator.ImageHandler(x, y, z, sysmbolSize)
		header := w.Header()
		header.Add("Content-Type", "image/png")
		png.Encode(w, m)
	}
}

http://localhost:800/gaode_tile_locator/?z=18&x=120&y=30&size=20
说明：z是地图级别；
x是经度
y是经度
size是图标大小



