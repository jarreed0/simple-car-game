package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	WIDTH     int32 = 1600
	HEIGHT    int32 = 900
	FONT_SIZE       = 32
	TILE_SIZE       = 64
	SPEED           = 30
	PI              = 3.14159265359
	MAP_SIZE        = 21
)

var (
	tile_img              int
	grid                  [MAP_SIZE][MAP_SIZE]int
	image                 []*sdl.Texture
	particles             []object
	running               bool
	up, down, right, left bool
	halfC                 = PI / 180

	renderer *sdl.Renderer
	window   *sdl.Window
	font     *ttf.Font
	fcolor   sdl.Color
	mouse    sdl.Point
	s1       = rand.NewSource(int64(sdl.GetTicks()))
	randNum  = rand.New(s1)
)

type object struct {
	dest, src sdl.Rect
	img       int
	c         sdl.Color
	X, Y      int32
	angle     float64
	vel       float64
	tick      int
	id        int
	count     int
	corners   [4]sdl.Point
}

var car, center, camera, car2, coin, particle object

func rotate(o object, c int, oX int32, oY int32, pw int32) sdl.Point {
	a := o.angle * (PI / 180)
	cX := o.X + (o.dest.W / 2)
	cY := o.Y + (o.dest.H / 2)
	tX := o.X - cX + oX
	tY := o.Y - cY + oY
	rX := (float64(tX) * math.Cos(a)) - (float64(tY) * math.Sin(a))
	rY := (float64(tX) * math.Sin(a)) + (float64(tY) * math.Cos(a))
	o.corners[c].X = int32(rX) + cX - (pw / 2)
	o.corners[c].Y = int32(rY) + cY - (pw / 2)
	return o.corners[c]
}

func calcCorners(o object, pw int32) object {
	o.corners[0] = rotate(o, 0, 0, 0, pw)
	o.corners[1] = rotate(o, 1, o.dest.W, 0, pw)
	o.corners[2] = rotate(o, 2, 0, o.dest.H, pw)
	o.corners[3] = rotate(o, 3, o.dest.W, o.dest.H, pw)
	return o
}

func get_degrees(input float64) float64 { return input * halfC }

func setGrid(r *object, gx int32, gy int32) (int32, int32) {
	return gx*int32(TILE_SIZE) - r.dest.W/2 + int32(TILE_SIZE/2), gy*int32(TILE_SIZE) - r.dest.H/2 + int32(TILE_SIZE/2)
}

func intersects(a, b sdl.Rect) bool {
	if (a.X < (b.X + b.W)) && ((a.X + a.W) > b.X) &&
		(a.Y < (b.Y + b.H)) && ((a.Y + a.H) > b.Y) {
		return true
	} else {
		return false
	}
}

var ca, cb sdl.Rect

func gridCol(a object, b object) bool {
	ca.X = a.X
	ca.Y = a.Y
	cb.X = b.X
	cb.Y = b.Y
	ca.W = a.dest.W
	ca.H = a.dest.H
	cb.W = b.dest.W
	cb.H = b.dest.H
	return intersects(ca, cb)
}

func pushParticle(px int32, py int32, angle float64) {
	particle.X = px
	particle.Y = py
	particle.angle = angle
	particles = append(particles, particle)
}

func updateCar(c object, u bool, d bool, l bool, r bool) object {
	if math.Abs(math.Round(c.vel)) < 1 {
		c.tick = 0
	} else {
		particle.src.X = 10
		particle.dest.W, particle.dest.H = 12, 12
		pushParticle(c.corners[0].X, c.corners[0].Y, c.angle)
		pushParticle(c.corners[2].X, c.corners[2].Y, c.angle)
	}
	c = calcCorners(c, particle.dest.W)
	if gridCol(c, coin) {
		c.count++
		coin.X, coin.Y = setGrid(&coin, int32(randNum.Intn(MAP_SIZE)), int32(randNum.Intn(MAP_SIZE)))
	}
	dx := int32(math.Cos(get_degrees(c.angle)) * c.vel)
	dy := int32(math.Sin(get_degrees(c.angle)) * c.vel)
	c.X += dx
	c.Y += dy

	if c.vel > 4 || c.vel < -4 {
		if l {
			c.angle -= 3
		}
		if r {
			c.angle += 3
		}
	}
	if u {
		c.vel += 0.5
	}
	if d {
		c.vel -= 0.5
	}
	if u || d {
		c.tick++
	}

	if !u && !d {
		if c.vel > 0 {
			c.vel -= 0.3
		}
		if c.vel < 0 {
			c.vel += 0.3
		}
	}
	if c.vel > float64(SPEED)/2 && c.tick < 100 {
		c.vel = float64(SPEED) / 2
	}
	if c.vel < float64(-SPEED)/2 && c.tick < 100 {
		c.vel = float64(-SPEED) / 2
	}
	if c.vel > float64(SPEED) {
		c.vel = float64(SPEED)
	}
	if c.vel < float64(-SPEED) {
		c.vel = float64(-SPEED)
	}
	return c
}

func setImage(filename string) int {
	i, _ := img.LoadTexture(renderer, filename)
	image = append(image, i)
	return len(image) - 1
}

func inCamera(r sdl.Rect) bool {
	return intersects(r, camera.dest)
}

func uc(r object) object {
	r.dest.X = r.X - camera.X
	r.dest.Y = r.Y - camera.Y
	return r
}

func draw(o *object) {
	if inCamera(o.dest) {
		renderer.CopyEx(image[o.img], &o.src, &o.dest, o.angle, nil, sdl.FLIP_NONE)
	}
}

func drawRect(r sdl.Rect) {
	if inCamera(r) {
		renderer.FillRect(&r)
	}
}

func drawOutline(r sdl.Rect) {
	if inCamera(r) {
		renderer.DrawRect(&r)
	}
}

func fontColor(r, g, b int) {
	fcolor.R = uint8(r)
	fcolor.G = uint8(g)
	fcolor.B = uint8(b)
}

/*var (
	surface *sdl.Surface
	texture *sdl.Texture
	wrect   sdl.Rect
)
func write(text string, x, y int32) {
	//surface, _ = font.RenderUTF8Solid(text, fcolor);
	surface, err := font.RenderUTF8Blended(text, fcolor)
	if err != nil {
		fmt.Printf("Failed to render text: %s\n", err)
	}
	texture, _ = renderer.CreateTextureFromSurface(surface)
	wrect.W = surface.W
	wrect.H = surface.H
	wrect.X = x - wrect.W
	wrect.Y = y - wrect.H
	surface.Free()
	renderer.Copy(texture, nil, &wrect)
	texture.Destroy()
}
*/
func setCamera(ox int32, oy int32, ow int32, oh int32) {
	camera.X = ox - WIDTH/2 + (ow / 2)
	camera.Y = oy - HEIGHT/2 + (oh / 2)
}
func setCameraXY(ox int32, oy int32) { setCamera(ox, oy, 0, 0) }
func setCameraObj(o object)          { setCamera(o.X, o.Y, o.dest.W, o.dest.H) }

func update() {
	fontColor(0, 0, 0)
	if gridCol(car, center) {
		fontColor(0, 255, 0)
	}
	if gridCol(car2, center) {
		fontColor(255, 0, 0)
	}

	if camera.id == car.id {
		car = updateCar(car, up, down, left, right)
		car2 = updateCar(car2, true, false, false, true)
	}
	if camera.id == car2.id {
		car2 = updateCar(car2, up, down, left, right)
		car = updateCar(car, true, false, false, true)
	}

	if camera.id == car.id {
		setCameraObj(car)
	}
	if camera.id == car2.id {
		setCameraObj(car2)
	}
	if camera.id == center.id {
		setCameraObj(center)
	}
	if camera.id == -1 {
		if up {
			camera.Y -= int32(SPEED / 2)
		}
		if down {
			camera.Y += int32(SPEED / 2)
		}
		if left {
			camera.X -= int32(SPEED / 2)
		}
		if right {
			camera.X += int32(SPEED / 2)
		}
	}
}

func input() {
	left, right, down, up = false, false, false, false
	keystates := sdl.GetKeyboardState()
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			println("Quit")
			running = false
			break
		}
	}

	if keystates[sdl.SCANCODE_ESCAPE] != 0 {
		running = false
	}
	if keystates[sdl.SCANCODE_W] != 0 || keystates[sdl.SCANCODE_UP] != 0 {
		up = true
	}
	if keystates[sdl.SCANCODE_S] != 0 || keystates[sdl.SCANCODE_DOWN] != 0 {
		down = true
	}
	if keystates[sdl.SCANCODE_A] != 0 || keystates[sdl.SCANCODE_LEFT] != 0 {
		left = true
	}
	if keystates[sdl.SCANCODE_D] != 0 || keystates[sdl.SCANCODE_RIGHT] != 0 {
		right = true
	}

	if keystates[sdl.SCANCODE_P] != 0 {
		camera.id = center.id
	}
	if keystates[sdl.SCANCODE_O] != 0 {
		camera.id = car.id
	}
	if keystates[sdl.SCANCODE_I] != 0 {
		camera.id = car2.id
	}
	if keystates[sdl.SCANCODE_U] != 0 {
		camera.id = -1
	}
	mouse.X, mouse.Y, _ = sdl.GetMouseState()
}

var d object
var hud, controls string

func render() {
	renderer.SetDrawColor(102, 75, 71, 255)
	renderer.Clear()

	d.dest.W, d.dest.H = int32(TILE_SIZE), int32(TILE_SIZE)
	d.src.W, d.src.H = int32(TILE_SIZE/2), int32(TILE_SIZE)
	d.src.X, d.src.Y = 0, 0
	for i := 0; i < MAP_SIZE; i++ {
		for j := 0; j < MAP_SIZE; j++ {
			d.X, d.Y = setGrid(&d, int32(i), int32(j))
			d = uc(d)
			if d.dest.X > WIDTH {
				break
			}
			if d.dest.X+int32(TILE_SIZE) < 0 {
				break
			}
			if d.dest.X+int32(TILE_SIZE) > 0 && d.dest.X-int32(TILE_SIZE) < WIDTH && d.dest.Y+int32(TILE_SIZE) > 0 && d.dest.Y-int32(TILE_SIZE) < HEIGHT {
				d.img = tile_img
				d.src.X = int32(grid[i][j] * 32)
				draw(&d)
				renderer.SetDrawColor(244, 147, 94, 255)
				drawOutline(d.dest)
			}
		}
	}

	d.X = int32((MAP_SIZE / 2) * TILE_SIZE)
	d.Y = int32((MAP_SIZE / 2) * TILE_SIZE)
	d = uc(d)
	renderer.SetDrawColor(0, 255, 0, 255)
	drawRect(d.dest)
	center = d
	center.id = 1
	for i := 0; i < len(particles); i++ {
		particles[i] = uc(particles[i])
		draw(&particles[i])
	}
	coin = uc(coin)
	renderer.SetDrawColor(255, 255, 0, 255)
	drawRect(coin.dest)
	car = uc(car)
	draw(&car)
	car2 = uc(car2)
	draw(&car2)
	//		if camera.id==car.id { write(string(car.X) + ", " + string(car.Y) + " " + string(camera.id) + " " + string(car.count), mouse.X, mouse.Y); }
	/*	if camera.id==car2.id { write(string(car2.X) + ", " + string(car2.Y) + " " + string(camera.id) + " " + string(car2.count), mouse.X, mouse.Y); }
		if camera.id==center.id { write(string(center.X) + ", " + string(center.Y) + " " + string(camera.id), mouse.X, mouse.Y); }
		if camera.id==-1 { write(string(camera.X+(WIDTH/2)) + ", " + string(camera.Y+(HEIGHT/2)) + " " + string(camera.id), mouse.X, mouse.Y); }
		if camera.id==car.id { hud="CAR 1"; }
		if camera.id==car2.id { hud="CAR 2"; }
		if camera.id==center.id { hud="CENTER"; }
		if camera.id==-1 { hud="CAMERA UNLOCKED"; }
		write(hud, int32(FONT_SIZE/2 * len(hud)) + 20, 50);
		write(controls, int32(FONT_SIZE/2 * len(controls)) + 20, HEIGHT - 25);
	*/
	renderer.Present()
}

func Init() {
	sw := true
	tile_img = setImage("tile.png")
	for i := 0; i < MAP_SIZE; i++ {
		for j := 0; j < MAP_SIZE; j++ {
			grid[j][i] = 0
			if sw {
				grid[j][i] = 1
			}
			sw = !sw
		}
	}
	camera.dest.X, camera.dest.H = 0, 0
	camera.dest.W = WIDTH
	camera.dest.H = HEIGHT
	car.src.X, car.src.Y = 0, 0
	car.src.W = 114/2
	car.src.H = 35
	car.dest.W = 40 * 2
	car.dest.H = 24 * 2
	car.dest.X = WIDTH / 2
	car.dest.Y = HEIGHT / 2
	car.img = setImage("car.png")
	car.X = int32(MAP_SIZE/2)*int32(TILE_SIZE) + int32(TILE_SIZE/2) - car.dest.W/2
	car.Y = int32(MAP_SIZE/2)*int32(TILE_SIZE) + int32(TILE_SIZE/2) - car.dest.H/2
	car.vel = 0
	fontColor(0, 0, 0)
	camera.id, car.id = 0, 0
	center.id = 1
	car2 = car
	car2.id = 2
	car2.X += 200
	car2.angle = 60
	car2.src.X = car2.src.W
	t := TILE_SIZE * .7
	coin.dest.W, coin.dest.H = int32(t), int32(t)
	coin.X, coin.Y = setGrid(&coin, int32(randNum.Intn(MAP_SIZE)), int32(randNum.Intn(MAP_SIZE)))
	car.X, car.Y = setGrid(&car, int32(randNum.Intn(MAP_SIZE)), int32(randNum.Intn(MAP_SIZE)))
	car.X, car.Y = setGrid(&car2, int32(randNum.Intn(MAP_SIZE)), int32(randNum.Intn(MAP_SIZE)))
	car.angle = float64(randNum.Intn(360))
	car2.angle = float64(randNum.Intn(360))
	car.count, car2.count = 0, 0
	controls = "WASD/ARROWS to move, POIU to change camera, ESC to close"
	particle.src.X, particle.src.Y = 0, 0
	particle.src.W, particle.src.H = 10, 10
	particle.dest.W, particle.dest.H = 12, 12
	particle.img = setImage("particle.png")
}

func main() {
	running = true

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0")
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()
	window, err := sdl.CreateWindow("Go Game", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WIDTH, HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	defer renderer.Destroy()

	window.SetFullscreen(sdl.WINDOW_FULLSCREEN)
	if err := ttf.Init(); err != nil {
		fmt.Println("failed at ttf init()")
	}
	font, err := ttf.OpenFont("pricedown.ttf", FONT_SIZE)
	if err != nil {
		fmt.Println("failed to open font")
	}
	defer font.Close()

	Init()
	for running {
		update()
		input()
		render()
	}
}
