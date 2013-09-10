package main

import (
	"fmt"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"os"
	"unsafe"
    "sort"
    //"math/rand"
    //"time"
    //"math"
	//"errors"
	//"github.com/go-gl/glu"
)

const (
	Title  = "Convex Hull"
	Width  = 800
	Height = 600
    Float32 float32 = 0.0
    Int32 int32 = 0
)

const fShaderSrc =`
#version 150

in vec3 Color;

out vec4 outColor;

void main()
{
    outColor = vec4( Color, 1.0 );
}`

const vShaderSrc = `
#version 150

uniform float time;

in vec2 position;
in vec3 color;

out vec3 Color;

void main()
{
    Color = color;
    gl_Position = vec4( position.x + sin(time), position.y + cos(4*time), 0.0, 1.0 );
}`

var (
    mouse [3]int
    posAttrib, colAttrib gl.AttribLocation
    ebo, vbo gl.Buffer
    uniTime gl.UniformLocation
    doHull bool
)

type Point struct {
    x, y float32
}

type Points []Point

func (p Points) Len() int {
    return len(p)
}
func (p Points) Swap(i, j int) {
    p[i], p[j] = p[j], p[i]
}
func (p Points) Less(i, j int) bool {
    if p[i].x == p[j].x {
        return p[i].y < p[j].y
    }
    return p[i].x < p[j].x
}

func setupShaders() {
    vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
    vertexShader.Source(vShaderSrc)
    vertexShader.Compile()

    fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
    fragmentShader.Source(fShaderSrc)
    fragmentShader.Compile()

    shaderProgram := gl.CreateProgram()
    shaderProgram.AttachShader(vertexShader)
    shaderProgram.AttachShader(fragmentShader)
    shaderProgram.BindFragDataLocation(0, "outColor")
    shaderProgram.Link()
    shaderProgram.Use()

    posAttrib = shaderProgram.GetAttribLocation("position")
    colAttrib = shaderProgram.GetAttribLocation("color")

    uniTime = shaderProgram.GetUniformLocation("time");
}

func onMouseBtn(button, state int) {
    mouse[button] = state
}

func onKey(key, state int) {
	switch key {
	case glfw.KeySpace:
        if state == 0 {
            doHull = true
        }
	}
}

func drawPoints(points []Point) {
    var (
        vertices []float32
        elements []int32
    )
    for i, p := range points {
        vertices = append(vertices, 2*p.x, 2*p.y, 1.0, 1.0, 1.0)
        elements = append(elements, int32(i))
    }
    time := float32(glfw.Time())
    uniTime.Uniform1f(time)

    vbo.Bind(gl.ARRAY_BUFFER)
    gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(Float32))*len(vertices), vertices, gl.DYNAMIC_DRAW)

    ebo.Bind(gl.ELEMENT_ARRAY_BUFFER)
    gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(unsafe.Sizeof(Int32))*len(elements), elements, gl.DYNAMIC_DRAW)

    posAttrib.EnableArray()
    posAttrib.AttribPointer(2, gl.FLOAT, false, int(unsafe.Sizeof(Float32))*5, uintptr(0))

    colAttrib.EnableArray()
    colAttrib.AttribPointer(3, gl.FLOAT, false, int(unsafe.Sizeof(Float32))*5, unsafe.Sizeof(Float32)*2)

    gl.Clear(gl.COLOR_BUFFER_BIT)
    gl.ClearColor(0.0, 0.0, 0.0, 1.0)

    gl.DrawElements(gl.POINTS, len(elements), gl.UNSIGNED_INT, uintptr(0));
}

func drawLines(points []Point) {
    var (
        vertices []float32
        elements []int32
    )
    for i, p := range points {
        vertices = append(vertices, p.x, p.y, 1.0, 1.0, 1.0)
        elements = append(elements, int32(i))
    }

    vbo.Bind(gl.ARRAY_BUFFER)
    gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(Float32))*len(vertices), vertices, gl.DYNAMIC_DRAW)

    ebo.Bind(gl.ELEMENT_ARRAY_BUFFER)
    gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(unsafe.Sizeof(Int32))*len(elements), elements, gl.DYNAMIC_DRAW)

    posAttrib.EnableArray()
    posAttrib.AttribPointer(2, gl.FLOAT, false, int(unsafe.Sizeof(Float32))*5, uintptr(0))

    colAttrib.EnableArray()
    colAttrib.AttribPointer(3, gl.FLOAT, false, int(unsafe.Sizeof(Float32))*5, unsafe.Sizeof(Float32)*2)

    gl.Clear(gl.COLOR_BUFFER_BIT)
    gl.ClearColor(0.0, 0.0, 0.0, 1.0)

    gl.DrawElements(gl.LINES, len(elements), gl.UNSIGNED_INT, uintptr(0));
}

func normalizeMouse(x, y int) (float32, float32) {
    x_out, y_out := float32(x)/Width, float32(y)/Height
    return 2*(x_out - 0.5), -2*(y_out - 0.5)
}

func fastConvexHull() {
    var points Points
    sort.Sort(points)
}

func turnsRight(a, b int) (bool) {
    return false
}

func main() {
	if err := glfw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "glfw: %s\n", err)
		return
	}
	defer glfw.Terminate()

	glfw.OpenWindowHint(glfw.WindowNoResize, 1)
	glfw.OpenWindowHint(glfw.OpenGLVersionMajor, 3)
	glfw.OpenWindowHint(glfw.OpenGLVersionMinor, 2)
    glfw.OpenWindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile);

	if err := glfw.OpenWindow(Width, Height, 0, 0, 0, 0, 16, 0, glfw.Windowed); err != nil {
		fmt.Fprintf(os.Stderr, "glfw: %s\n", err)
		return
	}
	defer glfw.CloseWindow()

	glfw.SetSwapInterval(1)
	glfw.SetWindowTitle(Title)
    glfw.SetMouseButtonCallback(onMouseBtn)
	glfw.SetKeyCallback(onKey)

    ver := gl.GetString(gl.VERSION)
    fmt.Println("GL Version", ver)

    if err := gl.Init(); err != 0 {
		fmt.Fprintf(os.Stderr, "OpenGL err:", err)
		return
    }
    e := gl.GetError() // Don't care about this error.. glfw bug

    vao := gl.GenVertexArray()
    vao.Bind()

    vbo = gl.GenBuffer()
    ebo = gl.GenBuffer()

    setupShaders()

    e = gl.GetError()
    fmt.Println(e)

    pressed := false
    var x, y float32
    var points Points
    //var hull Points
	for glfw.WindowParam(glfw.Opened) == 1 {
        // Add points if mouse is clicked
        if mouse[0] != 0 {
            x_, y_ := glfw.MousePos()
            x, y = normalizeMouse(x_, y_)
            pressed = true
        } else if pressed {
            points = append(points, Point{x, y})
            pressed = false
        }

        if doHull {
            doHull = false
            println("DO HUll")
        }
        pressed = false
        //drawLines(hull)
		glfw.SwapBuffers()
	}
}
