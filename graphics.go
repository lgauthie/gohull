package main

import (
    "github.com/go-gl/gl"
    "fmt"
    "os"
    "unsafe"
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

in vec2 position;
in vec3 color;

out vec3 Color;

void main()
{
    Color = color;
    gl_Position = vec4( position.x, position.y, 0.0, 1.0 );
}`

var (
    posAttrib, colAttrib gl.AttribLocation
    ebo, vbo gl.Buffer
)

func setupShaders() {

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

    e = gl.GetError()
    fmt.Println(e)
}

func drawPoints(points []Point) {
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

    gl.DrawElements(gl.POINTS, len(elements), gl.UNSIGNED_INT, uintptr(0));
}

func drawHull(points []Point) {
    var (
        vertices []float32
        elements []int32
    )
    vertices = append(vertices, points[0].x, points[0].y, 1.0, 1.0, 1.0)
    elements = append(elements, int32(0))
    for i, p := range points[1:] {
        vertices = append(vertices, p.x, p.y, 1.0, 1.0, 1.0)
        elements = append(elements, int32(i + 1), int32(i + 1))
    }
    elements = append(elements, int32(0))

    vbo.Bind(gl.ARRAY_BUFFER)
    gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(Float32))*len(vertices), vertices, gl.DYNAMIC_DRAW)

    ebo.Bind(gl.ELEMENT_ARRAY_BUFFER)
    gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(unsafe.Sizeof(Int32))*len(elements), elements, gl.DYNAMIC_DRAW)

    posAttrib.EnableArray()
    posAttrib.AttribPointer(2, gl.FLOAT, false, int(unsafe.Sizeof(Float32))*5, uintptr(0))

    colAttrib.EnableArray()
    colAttrib.AttribPointer(3, gl.FLOAT, false, int(unsafe.Sizeof(Float32))*5, unsafe.Sizeof(Float32)*2)

    gl.DrawElements(gl.LINES, len(elements), gl.UNSIGNED_INT, uintptr(0));
}

func clear() {
    gl.Clear(gl.COLOR_BUFFER_BIT)
    gl.ClearColor(0.0, 0.0, 0.0, 1.0)
}
