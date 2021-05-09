package main

import (
	"fmt"
	"runtime"
	"strings"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/gl/v4.6-core/gl"
)

const (
	vertexShaderSource = `
		#version 410

		in vec2 vp;
		in vec3 color;

		out vec3 out_color;

		void main() {
			out_color = color;
			gl_Position = vec4(vp, 1.0, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 410

		in vec3 out_color;

		out vec4 frag_colour;

		void main() {
			frag_colour = vec4(out_color, 1.0);
		}
	` + "\x00"
	BOTTOM_BAR_HEIGHT = 32.0
)

var (
    mesh []float32
	window_w int
	window_h int
	COLOR_VAO_OFFSET = 2 * 4
)

func compileShader(source string, shaderType uint32) (uint32, error) {
    shader := gl.CreateShader(shaderType)

    csources, free := gl.Strs(source)
    gl.ShaderSource(shader, 1, csources, nil)
    free()
    gl.CompileShader(shader)

    var status int32
    gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
    if status == gl.FALSE {
        var logLength int32
        gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

        log := strings.Repeat("\x00", int(logLength+1))
        gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

        return 0, fmt.Errorf("failed to compile %v: %v", source, log)
    }

    return shader, nil
}

func assembleVAO() uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	var vao uint32
    gl.GenVertexArrays(1, &vao)
    gl.BindVertexArray(vao)

	gl.EnableVertexAttribArray(0)
    gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 5 * 4, gl.PtrOffset(0 * 4))

	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 5 * 4, gl.PtrOffset(2 * 4))

    return vao
}

func resize(win *glfw.Window, w int, h int) {
	fmt.Println(w, h)
	gl.Viewport(0, 0, int32(w), int32(h))
	window_w = w
	window_h = h
}

func posToGL(x float32, y float32) (float32, float32) {
	return x/float32(window_w) * 2 - 1, y/float32(window_h) * 2 - 1
}

func meshAddRect(x float32, y float32, w float32, h float32, r float32, g float32, b float32) {
	gx, gy := posToGL(x, y)
	gw, gh := posToGL(w, h)
	gw += 1
	gh += 1
	mesh = append(mesh, gx, gy, r, g, b)
	mesh = append(mesh, gx + gw, gy, r, g, b)
	mesh = append(mesh, gx, gy + gh, r, g, b)
	mesh = append(mesh, gx + gw, gy, r, g, b)
	mesh = append(mesh, gx, gy + gh, r, g, b)
	mesh = append(mesh, gx + gw, gy + gh, r, g, b)
}

func assembleMesh() {
	mesh = nil
	meshAddRect(0, 0, float32(window_w), BOTTOM_BAR_HEIGHT, 28 / 256.0, 31 / 256.0, 36 / 256.0)
}

func render(vao uint32, gl_prog uint32, mesh []float32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(gl_prog)

	gl.BufferData(gl.ARRAY_BUFFER, 4*len(mesh), gl.Ptr(mesh), gl.STREAM_DRAW)
    gl.DrawArrays(gl.TRIANGLES, 0, int32(len(mesh) / 5))
}

func main() {
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("could not initialize glfw: %v", err))
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	win, err := glfw.CreateWindow(800, 600, "Hello world", nil, nil)
	if err != nil {
		panic(fmt.Errorf("could not create opengl renderer: %v", err))
	}
	win.SetSizeCallback(resize)

	win.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		panic(err)
	}

    vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
    if err != nil {
        panic(err)
    }
    fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
    if err != nil {
        panic(err)
    }

    gl_prog := gl.CreateProgram()
    gl.AttachShader(gl_prog, vertexShader)
    gl.AttachShader(gl_prog, fragmentShader)
    gl.LinkProgram(gl_prog)

	gl.ClearColor(35/256.0, 39/256.0, 46/256.0, 1.0)

	vao := assembleVAO()
	gl.BindVertexArray(vao)

	for !win.ShouldClose() {
		assembleMesh()
		render(vao, gl_prog, mesh)
		win.SwapBuffers()
		glfw.PollEvents()
	}
}
