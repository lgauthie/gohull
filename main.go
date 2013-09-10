package main

import (
	"fmt"
	"github.com/go-gl/glfw"
	"os"
)

const (
	Title  = "Convex Hull"
	Width  = 800
	Height = 600
    Float32 float32 = 0.0
    Int32 int32 = 0
)

var (
    mouse [3]int
    spaceReleased bool
)


func onMouseBtn(button, state int) {
    mouse[button] = state
}

func onKey(key, state int) {
	switch key {
	case glfw.KeySpace:
        if state == 0 {
            spaceReleased = true
        }
	}
}

func normalizeMouse(x, y int) (x_out, y_out float32) {
    x_out, y_out = float32(x)/Width, float32(y)/Height
    return 2*(x_out - 0.5), -2*(y_out - 0.5)
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
    setupShaders()

    pressed := false
    var x, y float32
    points := Points{}
    hull := Points{}
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

        if spaceReleased {
            hull = fastConvexHull(points)
            spaceReleased = false
        }

        clear()
        if len(points) > 0 {
            drawPoints(points)
        }
        if len(hull) > 0 {
            drawHull(hull)
        }
		glfw.SwapBuffers()
	}
}
