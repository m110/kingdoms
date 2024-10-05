//go:build ios

package storage

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation
#include <stdlib.h>
const char* GetApplicationSupportDirectory();
*/
import "C"
import (
	"path"
)

func GetApplicationSupportDirectoryPath() string {
	dir := C.GetApplicationSupportDirectory()
	return C.GoString(dir)
}

func storageDirectory() (string, error) {
	return path.Join(GetApplicationSupportDirectoryPath(), "kingdoms"), nil
}
