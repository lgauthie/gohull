package main

import (
	"fmt"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"os"
	//"unsafe"
    //"math"
    //"time"
	//"errors"
	//"github.com/go-gl/glu"
)

const (
	Title  = "The Start"
	Width  = 800
	Height = 600
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
    gl_Position = vec4( position, 0.0, 1.0 );
}`

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

    ver := gl.GetString(gl.VERSION)
    fmt.Println("GL Version", ver)

    if err := gl.Init(); err != 0 {
		fmt.Fprintf(os.Stderr, "OpenGL err:", err)
		return
    }
    e := gl.GetError() // Don't care about this error.. glfw bug

    vao := gl.GenVertexArray()
    vao.Bind()

    // Create an array buffer of points that will be the
    // vertexes of a triangle.
    vertices := []float32{
         0.0,  0.5, 1.0, 0.0, 0.0, // Vertex 1: Red
         0.5, -0.5, 0.0, 1.0, 0.0, // Vertex 2: Green
        -0.5, -0.5, 0.0, 0.0, 1.0, // Vertex 3: Blue
    }
    vbo := gl.GenBuffer()
    vbo.Bind(gl.ARRAY_BUFFER)
    sizeofFloat := 4
    gl.BufferData(gl.ARRAY_BUFFER, sizeofFloat*len(vertices), vertices, gl.STATIC_DRAW)

    // SETUP SHADERS
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

    posAttrib := shaderProgram.GetAttribLocation("position")
    var offset uintptr = 0
    posAttrib.EnableArray()
    posAttrib.AttribPointer(2, gl.FLOAT, false, sizeofFloat*5, offset)

    colAttrib := shaderProgram.GetAttribLocation("color")
    offset = uintptr(sizeofFloat*2)
    colAttrib.EnableArray()
    colAttrib.AttribPointer(3, gl.FLOAT, false, sizeofFloat*5, offset)

    e = gl.GetError()
    fmt.Println(e)

	for glfw.WindowParam(glfw.Opened) == 1 {
        gl.Clear(gl.COLOR_BUFFER_BIT)
        gl.ClearColor(0.0, 0.0, 0.0, 1.0)
        gl.DrawArrays(gl.LINE_LOOP, 0, 3);
		glfw.SwapBuffers()
	}
}
