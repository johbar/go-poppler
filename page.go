package poppler

// #cgo pkg-config: poppler-glib
// #include <poppler.h>
// #include <glib.h>
import "C"
import (
	"runtime"
	"unsafe"
)

type Page struct {
	p *C.struct__PopplerPage
}

func (p *Page) Text() string {
	cCharPtr := C.poppler_page_get_text(p.p)
	pageText := C.GoString(cCharPtr)
	C.free(unsafe.Pointer(cCharPtr))
	return pageText
}

func (p *Page) Size() (width, height float64) {
	var w, h C.double
	C.poppler_page_get_size(p.p, &w, &h)
	return float64(w), float64(h)
}

func (p *Page) Index() int {
	return int(C.poppler_page_get_index(p.p))
}

func (p *Page) Label() string {
	return toString(C.poppler_page_get_label(p.p))
}

func (p *Page) Duration() float64 {
	return float64(C.poppler_page_get_duration(p.p))
}

// Close frees memory allocated when Poppler opened the page
func (p *Page) Close() {
	// GC/finalizer shouldn't try to free C memory that has been freed already
	runtime.SetFinalizer(p, nil)
	closePage(p)
	p.p = nil
}

func closePage(p *Page) {
	C.g_object_unref(C.gpointer(p.p))
}
