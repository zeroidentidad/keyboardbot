package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/keyboard"
)

func main() {
	keys := keyboard.NewDriver()
	currentStr := make([]byte, 32)

	work := func() {
		_ = keys.On(keyboard.Key, func(data interface{}) {
			key := data.(keyboard.KeyEvent)

			if key.Key == 0 {
				fmt.Println(" => [Intro/Enter] hacer búsqueda...")
				_ = keys.Halt()

				str := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(string(currentStr), "")
				search := fmt.Sprintf("https://google.com/search?q=%s", str)
				err := open(search)
				if err != nil {
					fmt.Println(err.Error())
				}

				_ = keys.Start()
				currentStr = []byte("")
			} else {
				fmt.Print(key.Char)
				currentStr = append(currentStr, key.Char...)
			}
		})
	}

	robot := gobot.NewRobot("keyboardbot",
		[]gobot.Connection{},
		[]gobot.Device{keys},
		work,
	)

	_ = robot.Start()
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		// "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	run := exec.Command(cmd, args...)
	err := run.Start()
	return err
}
