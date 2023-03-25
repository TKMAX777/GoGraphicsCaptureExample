package main

import (
	"errors"
	"fmt"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/TKMAX777/GoGraphicsCaptureExample/winapi/winrt"
	"github.com/go-ole/go-ole"
	"github.com/lxn/win"
)

// Direct3D11CaptureFramePool

type Direct3D11CaptureFramePool struct {
	ole.IUnknown
}

type Direct3D11CaptureFramePoolVtbl struct {
	ole.IUnknownVtbl
	Invoke     uintptr
	originalV  *Direct3D11CaptureFramePool
	counter    int
	marshaller *ole.IUnknown
}

func NewDirect3D11CaptureFramePool(invoke winrt.Direct3D11CaptureFramePoolFrameArrivedProcType) *Direct3D11CaptureFramePool {
	var v = &Direct3D11CaptureFramePoolVtbl{
		Invoke:     syscall.NewCallback(invoke),
		originalV:  nil,
		counter:    1,
		marshaller: nil,
	}

	var newV = new(Direct3D11CaptureFramePool)

	newV.RawVTable = (*interface{})(unsafe.Pointer(v))

	v.QueryInterface = syscall.NewCallback(Direct3D11CaptureFramePoolQueryInterface)
	v.AddRef = syscall.NewCallback(Direct3D11CaptureFramePoolAddRef)
	v.Release = syscall.NewCallback(Direct3D11CaptureFramePoolRelease)
	// Protect from gabage collecter
	v.originalV = newV

	fmt.Println(newV.VTable())
	var err error
	var Unknown *ole.IUnknown
	err = newV.PutQueryInterface(ole.IID_IUnknown, &Unknown)
	if err != nil {
		panic(err)
	}

	v.marshaller, err = winrt.CoCreateFreeThreadedMarshaler(Unknown)
	if err != nil {
		panic(err)
	}

	return newV
}

func (v *Direct3D11CaptureFramePool) VTable() *Direct3D11CaptureFramePoolVtbl {
	return (*Direct3D11CaptureFramePoolVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *Direct3D11CaptureFramePool) Invoke(sender *winrt.IDirect3D11CaptureFramePool, args *ole.IInspectable) error {
	r1, _, _ := syscall.SyscallN(v.VTable().Invoke, uintptr(unsafe.Pointer(sender)), uintptr(unsafe.Pointer(args)))
	return ole.NewError(r1)
}

// QueryInterface(vp *Direct3D11CaptureFramePool, riid ole.GUID, lppvObj **ole.Inspectable)
func Direct3D11CaptureFramePoolQueryInterface(lpMyObj *uintptr, riid *uintptr, lppvObj **uintptr) (hr uintptr) {
	// Validate input
	if lpMyObj == nil {
		return win.E_INVALIDARG
	}

	var V = (*Direct3D11CaptureFramePool)(unsafe.Pointer(lpMyObj))
	var err error
	// Check dereferencability
	func() {
		defer func() {
			if recover() != nil {
				err = errors.New("InvalidObject")
			}
		}()
		// if object cannot be dereferenced, then panic occurs
		V.VTable()
	}()
	if err != nil {
		return win.E_INVALIDARG
	}

	*lppvObj = nil
	var id = (*ole.GUID)(unsafe.Pointer(riid))

	// Convert
	fmt.Println(id.String())
	switch id.String() {
	case ole.IID_IUnknown.String(), winrt.ITypedEventHandlerID.String():
		fmt.Println("OK")
		V.AddRef()
		*lppvObj = (*uintptr)(unsafe.Pointer(V))

		fmt.Println(V.VTable())
		return win.S_OK
	case winrt.IMarshalID.String():
		qi, err := V.VTable().marshaller.QueryInterface(winrt.IMarshalID)
		if err != nil {
			fmt.Println("MarshallingError: ", err)
			return win.E_UNEXPECTED
		}
		*lppvObj = (*uintptr)(unsafe.Pointer(qi))
		return win.S_OK
	default:
		return win.E_NOINTERFACE
	}
}

func Direct3D11CaptureFramePoolAddRef(lpMyObj *uintptr) uintptr {
	// Validate input
	if lpMyObj == nil {
		return 0
	}

	var V = (*Direct3D11CaptureFramePool)(unsafe.Pointer(lpMyObj))
	V.VTable().counter++

	return uintptr(V.VTable().counter)
}

func Direct3D11CaptureFramePoolRelease(lpMyObj *uintptr) uintptr {
	// Validate input
	if lpMyObj == nil {
		return 0
	}

	var V = (*Direct3D11CaptureFramePool)(unsafe.Pointer(lpMyObj))
	V.VTable().counter--

	if V.VTable().counter == 0 {
		V.VTable().originalV = nil
		runtime.GC()
	}

	return uintptr(V.VTable().counter)
}
