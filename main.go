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
		in vec3 vp;
		void main() {
			gl_Position = vec4(vp, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 410
		out vec4 frag_colour;
		void main() {
			frag_colour = vec4(1, 1, 1, 1);
		}
	` + "\x00"
)

var (
    triangle = []float32{
        0, 0.5, 0, // top
        -0.5, -0.5, 0, // left
        0.5, -0.5, 0, // right
    }
	window_w int
	window_h int
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

func assembleVAO(mesh []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    gl.BufferData(gl.ARRAY_BUFFER, 4*len(mesh), gl.Ptr(mesh), gl.STATIC_DRAW)

	var vao uint32
    gl.GenVertexArrays(1, &vao)
    gl.BindVertexArray(vao)
    gl.EnableVertexAttribArray(0)
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

    return vao
}

func resize(win *glfw.Window, w int, h int) {
	fmt.Println(w, h)
	gl.Viewport(0, 0, int32(w), int32(h))
	window_w = w
	window_h = h
}

func render(vao uint32, gl_prog uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(gl_prog)

	gl.BindVertexArray(vao)
    gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle) / 3))
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

	vao := assembleVAO(triangle)

	for !win.ShouldClose() {
		render(vao, gl_prog)
		win.SwapBuffers()
		glfw.PollEvents()
	}
}
