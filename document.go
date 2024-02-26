package poppler

// #cgo pkg-config: poppler-glib
// #include <poppler.h>
// #include <stdlib.h>
// #include <glib.h>
// #include <unistd.h>
import "C"
import (
	"log"
	"os"
	"runtime"
)

type Document struct {
	doc poppDoc
	// keep this so it does not get GCed
	fd *os.File
}

type DocumentInfo struct {
	PdfVersion, Title, Author, Subject, KeyWords, Creator, Producer, Metadata string
	CreationDate, ModificationDate, Pages                                     int
	IsLinearized                                                              bool
}

func (d *Document) Info() DocumentInfo {
	return DocumentInfo{
		PdfVersion:       toString(C.poppler_document_get_pdf_version_string(d.doc)),
		Title:            toString(C.poppler_document_get_title(d.doc)),
		Author:           toString(C.poppler_document_get_author(d.doc)),
		Subject:          toString(C.poppler_document_get_subject(d.doc)),
		KeyWords:         toString(C.poppler_document_get_keywords(d.doc)),
		Creator:          toString(C.poppler_document_get_creator(d.doc)),
		Producer:         toString(C.poppler_document_get_producer(d.doc)),
		Metadata:         toString(C.poppler_document_get_metadata(d.doc)),
		CreationDate:     int(C.poppler_document_get_creation_date(d.doc)),
		ModificationDate: int(C.poppler_document_get_modification_date(d.doc)),
		Pages:            int(C.poppler_document_get_n_pages(d.doc)),
		IsLinearized:     toBool(C.poppler_document_is_linearized(d.doc)),
	}
}

func (d *Document) GetNPages() int {
	return int(C.poppler_document_get_n_pages(d.doc))
}

func (d *Document) GetPage(i int) (page *Page) {
	p := C.poppler_document_get_page(d.doc, C.int(i))
	page = &Page{p: p}
	// Make sure the C memory gets freed at some point
	// even if the user doesn't call p.Close():
	runtime.SetFinalizer(page, closePage)
	return
}

func (d *Document) HasAttachments() bool {
	return toBool(C.poppler_document_has_attachments(d.doc))
}

func (d *Document) GetNAttachments() int {
	return int(C.poppler_document_get_n_attachments(d.doc))
}

// Close releases memory and ressources allocated by Poppler
func (d *Document) Close() {
	// GC shouldn't try to free C memory that has been freed already:
	runtime.SetFinalizer(d, nil)
	closeDocument(d)
}

func closeDocument(d *Document) {
	C.g_object_unref(C.gpointer(d.doc))
	if d.fd != nil {
		log.Printf("Closing FD %d", d.fd.Fd())
		d.fd.Close()
	}
}
