package main

import (
	"fmt"
	"strings"

	"github.com/TKMAX777/GoGraphicsCaptureExample/winapi"
	"github.com/lxn/win"
)

func main() {
	var rdHwnd win.HWND
	for {
		rdHwnd = winapi.FindWindowEx(0, rdHwnd, winapi.MustUTF16PtrFromString("Chrome_WidgetWin_1"), nil)
		if rdHwnd == 0 {
			win.MessageBox(0, winapi.MustUTF16PtrFromString("Could not find window"), winapi.MustUTF16PtrFromString("RDP Relative Input"), win.MB_ICONERROR)
			return
		}
		var name = strings.TrimSpace(winapi.GetWindowTextString(rdHwnd))
		if name != "" {
			fmt.Println(name)
			break
		}
	}

	var handler = CaptureHandler{}

	err := handler.StartCapture(rdHwnd)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer handler.Close()

}
