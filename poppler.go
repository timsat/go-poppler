package poppler

// #cgo pkg-config: poppler-glib
// #include <poppler.h>
// #include <stdlib.h>
// #include <glib.h>
// #include <unistd.h>
import "C"

import (
	"errors"
	_ "log"
	"path/filepath"
	"runtime"
	"unsafe"
)

type poppDoc *C.struct__PopplerDocument

func Open(filename string) (doc *Document, err error) {
	filename, err = filepath.Abs(filename)
	if err != nil {
		return
	}
	var e *C.GError
	fn := C.g_filename_to_uri((*C.gchar)(C.CString(filename)), nil, nil)
	var d poppDoc = C.poppler_document_new_from_file((*C.char)(fn), nil, &e)
	if e != nil {
		err = errors.New(C.GoString((*C.char)(e.message)))
	}
	doc = &Document{
		doc: d,
	}
	//Ensure C memory is freed even if the user forgets to call Close()
	runtime.SetFinalizer(doc, closeDocument)
	return
}

func Load(data []byte) (doc *Document, err error) {
	var e *C.GError
	var d poppDoc
	var b *C.GBytes = C.g_bytes_new((C.gconstpointer)(unsafe.Pointer(&data[0])), (C.ulong)(len(data)))
	d = C.poppler_document_new_from_bytes(b, nil, &e)
	C.g_bytes_unref(b)
	if e != nil {
		err = errors.New(C.GoString((*C.char)(e.message)))
	}
	doc = &Document{
		doc:   d,
	}
	//Ensure C memory is freed even if the user forgets to call Close()
	runtime.SetFinalizer(doc, closeDocument)
	return
}

func Version() string {
	return C.GoString(C.poppler_get_version())
}
