package cui

import (
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"unsafe"
)

// Screen is handling terminal window output
type Screen struct {
	size       winsize
	Width      int
	Height     int
	clearFuncs map[string]func()
	buffer     [][]rune
}

// Winsize stores terminal window sizes
type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// Init Screen params
func (scr *Screen) Init() {

	// Set clear screen system calls functions
	scr.clearFuncs = make(map[string]func())
	scr.clearFuncs["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	scr.clearFuncs["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	scr.UpdateSize()

	scr.buffer = make([][]rune, scr.Height)

	for i := 0; i < scr.Height; i++ {
		scr.buffer[i] = make([]rune, scr.Width)
		for j := range scr.buffer[i] {
			scr.buffer[i][j] = ' '
		}
	}
}

// SetRune sets rune at given coords on screen
func (scr *Screen) SetRune(x, y int, r rune) {
	//fmt.Printf("SetRune: x=%v, y=%v, rune=%q\r\n", x, y, r)
	if x >= scr.Width || y >= scr.Height {
		return
	}

	scr.buffer[y][x] = r
}

// Clear screen by system call from clearFunc map
func (scr *Screen) Clear() {
	clearFunc, ok := scr.clearFuncs[runtime.GOOS]
	if ok {
		clearFunc()
	} else {
		panic("Cannot clear screen, unsupported OS!")
	}
}

// SendToDisplay send screen buffer to std out to display it
func (scr *Screen) SendToDisplay() {
	scr.Clear()
	for _, runesRow := range scr.buffer {
		os.Stdout.WriteString(string(runesRow))
	}
}

// UpdateSize get current info about terminal screen size
func (scr *Screen) UpdateSize() {

	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&scr.size)))

	if int(retCode) == -1 {
		panic(errno)
	}

	scr.Width = int(scr.size.Col)
	scr.Height = int(scr.size.Row)
}
