package main

/*
#cgo linux openbsd freebsd pkg-config: gtk+-3.0 webkit2gtk-4.0
#include <gtk/gtk.h>
#include <gdk-pixbuf/gdk-pixbuf.h>

void setWindowIconFromBytes(GtkWindow* window, unsigned char* buf, int len) {
    GInputStream* stream = g_memory_input_stream_new_from_data(buf, len, NULL);
    GdkPixbuf* pixbuf = gdk_pixbuf_new_from_stream(stream, NULL, NULL);
    gtk_window_set_icon(window, pixbuf);
    g_object_unref(stream);
    g_object_unref(pixbuf);
}

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

	icon, _ := efs.ReadFile("icon64.png")

	w := webview.New(false)
	defer w.Destroy()
	w.SetTitle("Everview")
	w.SetSize(1200, 800, webview.HintNone)
	C.setWindowIconFromBytes((*C.GtkWindow)(w.Window()), (*C.uchar)(&icon[0]), C.int(len(icon)))

	w.Navigate(url)
	systray.Register(onReady(w, icon))

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
