package main

/*
#cgo linux openbsd freebsd pkg-config: gtk+-3.0 webkit2gtk-4.0
#include <gtk/gtk.h>
*/
import "C"

import (
	"embed"
	"os"

	"github.com/ghostiam/systray"
	"github.com/webview/webview"
)

//go:embed icon64.png
var efs embed.FS

const url = "https://www.evernote.com/client/web"

func main() {

	iconBytes, _ := efs.ReadFile("icon64.png")

	w := webview.New(false)
	defer w.Destroy()
	w.SetTitle("Everview")
	w.SetSize(1200, 800, webview.HintNone)
	C.gtk_window_set_icon_from_file((*C.GtkWindow)(w.Window()), C.CString("icon64.png"), nil)

	w.Navigate(url)
	systray.Register(onReady(w, iconBytes))

	if len(os.Args) == 2 && os.Args[1] == "--hidden" {
		C.gtk_widget_hide((*C.GtkWidget)(w.Window()))
	}

	w.Run()
}

func onReady(w webview.WebView, iconBytes []byte) func() {
	return func() {
		systray.SetTitle("Everview")
		systray.SetTooltip("Everview")
		systray.SetIcon(iconBytes)

		mShow := systray.AddMenuItem("Show", "Show")
		mHide := systray.AddMenuItem("Hide", "Hide")
		mQuit := systray.AddMenuItem("Quit", "Quit")

		go func() {
			for {
				select {
				case <-mShow.ClickedCh:
					C.gtk_widget_show_all((*C.GtkWidget)(w.Window()))
				case <-mHide.ClickedCh:
					C.gtk_widget_hide((*C.GtkWidget)(w.Window()))
				case <-mQuit.ClickedCh:
					w.Terminate()
					return
				}
			}
		}()
	}
}
