// Copyright wasilibs authors
// SPDX-License-Identifier: MIT

//go:build gc.custom

package nottinygc

/*
void GC_register_finalizer(void* obj, void* fn, void* cd, void** ofn, void** ocn);
void onFinalizer(void* obj, void* fn);
void* malloc(unsigned int long);
void free(void* ptr);
*/
import "C"
import "unsafe"

var finalizers = map[uint64]interface{}{}

var finalizersIdx = uint64(0)

//go:linkname SetFinalizer runtime.SetFinalizer
func SetFinalizer(obj interface{}, finalizer interface{}) {
	if _, ok := finalizer.(func(interface{})); !ok {
		// Until function invocation with reflection is supported by TinyGo, we cannot support arbitrary finalizers.
		// We make a best-effort attempt to support the generic func(interface{}). For other finalizers, we silently
		// ignore, which would be the behavior for all finalizers with the default allocator.
		return
	}

	finalizersIdx++
	finKey := finalizersIdx
	finalizers[finKey] = finalizer

	in := (*_interface)(unsafe.Pointer(&obj))

	rf := (*registeredFinalizer)(C.malloc(C.ulong(unsafe.Sizeof(registeredFinalizer{}))))
	rf.typecode = in.typecode
	rf.finKey = finKey

	C.GC_register_finalizer(in.value, C.onFinalizer, unsafe.Pointer(rf), nil, nil)
}

//export onFinalizer
func onFinalizer(obj unsafe.Pointer, data unsafe.Pointer) {
	defer C.free(data)

	rf := (*registeredFinalizer)(data)
	finalizer := finalizers[rf.finKey]
	delete(finalizers, rf.finKey)

	var in interface{}
	inFields := (*_interface)(unsafe.Pointer(&in))
	inFields.typecode = rf.typecode
	inFields.value = obj

	switch f := finalizer.(type) {
	case func(interface{}):
		f(in)
	default:
		panic("BUG: SetFinalizer should have filtered out non-supported finalizer interface")
	}
}

type _interface struct {
	typecode uintptr
	value    unsafe.Pointer
}

type registeredFinalizer struct {
	typecode uintptr
	finKey   uint64
}
