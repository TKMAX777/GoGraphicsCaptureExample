package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/TKMAX777/GoGraphicsCaptureExample/winapi/dx11"
	"github.com/TKMAX777/GoGraphicsCaptureExample/winapi/winrt"
	"github.com/go-ole/go-ole"
	"github.com/lxn/win"
	"github.com/pkg/errors"
)

type CaptureHandler struct {
	device                 *winrt.IDirect3DDevice
	graphicsCaptureItem    *winrt.IGraphicsCaptureItem
	framePool              *winrt.IDirect3D11CaptureFramePool
	graphicsCaptureSession *winrt.IGraphicsCaptureSession
	framePoolToken         *winrt.EventRegistrationToken
}

func (c *CaptureHandler) StartCapture(hwnd win.HWND) error {
	err := winrt.RoInitialize(winrt.RO_INIT_SINGLETHREADED)
	if err != nil {
		return errors.Wrap(err, "RoInitialize")
	}

	var device *dx11.ID3D11Device
	err = dx11.D3DCreateDevice(nil, dx11.D3D_DRIVER_TYPE_HARDWARE, 0, dx11.D3D11_CREATE_DEVICE_BGRA_SUPPORT, nil, 0, dx11.D3D11_SDK_VERSION, &device, nil, nil)
	if err != nil {
		return errors.Wrap(err, "D3DCreateDevice")
	}

	var dxgiDevice *dx11.IDXGIDevice
	err = device.PutQueryInterface(dx11.IDXGIDeviceID, &dxgiDevice)
	if err != nil {
		return errors.Wrap(err, "PutQueryInterface")
	}

	var deviceRT *ole.IInspectable

	err = dx11.CreateDirect3D11DeviceFromDXGIDevice(dxgiDevice, &deviceRT)
	if err != nil {
		return errors.Wrap(err, "CreateDirect3D11DeviceFromDXGIDevice")
	}

	err = deviceRT.PutQueryInterface(winrt.IDirect3DDeviceID, &c.device)
	if err != nil {
		return errors.Wrap(err, "QueryInterface: IDirect3DDeviceID")
	}

	factory, err := ole.RoGetActivationFactory(winrt.GraphicsCaptureItemClass, winrt.IGraphicsCaptureItemInteropID)
	if err != nil {
		return errors.Wrap(err, "RoGetActivationFactory: IGraphicsCaptureItemID")
	}

	var interop *winrt.IGraphicsCaptureItemInterop
	err = factory.PutQueryInterface(winrt.IGraphicsCaptureItemInteropID, &interop)
	if err != nil {
		return errors.Wrap(err, "QueryInterface: IGraphicsCaptureItemInteropID")
	}

	var captureItemDispatch *ole.IInspectable
	// err = interop.CreateForWindow(hwnd, winrt.IGraphicsCaptureItemID, &captureItemDispatch)
	// if err != nil {
	// 	return errors.Wrap(err, "CreateForWindow")
	// }

	var hmoni = win.MonitorFromWindow(hwnd, win.MONITORINFOF_PRIMARY)

	err = interop.CreateForMonitor(hmoni, winrt.IGraphicsCaptureItemID, &captureItemDispatch)
	if err != nil {
		return errors.Wrap(err, "CreateForMonitor")
	}

	err = captureItemDispatch.PutQueryInterface(winrt.IGraphicsCaptureItemID, &c.graphicsCaptureItem)
	if err != nil {
		return errors.Wrap(err, "PutQueryInterface captureItemDispatch")
	}

	size, err := c.graphicsCaptureItem.Size()
	if err != nil {
		return errors.Wrap(err, "Size")
	}

	ins, err := ole.RoGetActivationFactory(winrt.Direct3D11CaptureFramePoolClass, winrt.IDirect3D11CaptureFramePoolStaticsID)
	if err != nil {
		return errors.Wrap(err, "RoGetActivationFactory: IDirect3D11CaptureFramePoolStatics Class Instance")
	}
	defer ins.Release()

	var framePoolStatic *winrt.IDirect3D11CaptureFramePoolStatics
	err = ins.PutQueryInterface(winrt.IDirect3D11CaptureFramePoolStaticsID, &framePoolStatic)
	if err != nil {
		return errors.Wrap(err, "PutQueryInterface: IDirect3D11CaptureFramePoolStaticsID")
	}

	c.framePool, err = framePoolStatic.Create(c.device, winrt.DirectXPixelFormat_B8G8R8A8UIntNormalized, 2, size)
	if err != nil {
		return errors.Wrap(err, "CreateFramePool")
	}

	var eventObject winrt.Direct3D11CaptureFramePoolVtbl
	eventObject.Invoke = syscall.NewCallback(c.onFrameArrived)

	c.framePoolToken, err = c.framePool.AddFrameArrived(eventObject)
	if err != nil {
		return errors.Wrap(err, "AddFrameArrived")
	}

	c.graphicsCaptureSession, err = c.framePool.CreateCaptureSession(c.graphicsCaptureItem)
	if err != nil {
		return errors.Wrap(err, "CreateCaptureSession")
	}

	err = c.graphicsCaptureSession.StartCapture()
	if err != nil {
		return errors.Wrap(err, "StartCapture")
	}

	return nil
}

func (c *CaptureHandler) onFrameArrived(sender *winrt.IDirect3D11CaptureFramePool, args *ole.IInspectable) uintptr {
	_, err := sender.TryGetNextFrame()
	if err != nil {
		os.Stderr.Write([]byte("Error: TryGetNextFrame: " + err.Error()))
		return 0
	}

	fmt.Println("Arrived")

	return 0
}

func (c *CaptureHandler) Close() error {
	if c.framePool != nil {
		err := c.framePool.RemoveFrameArrived(c.framePoolToken)
		if err != nil {
			return errors.Wrap(err, "RemoveFrameArrived")
		}

		var closable *winrt.IClosable
		err = c.framePool.PutQueryInterface(winrt.IClosableID, &closable)
		if err != nil {
			return errors.Wrap(err, "PutQueryInterface: graphicsCaptureSession")
		}
		defer closable.Release()

		closable.Close()

		c.framePool = nil
	}

	var closable *winrt.IClosable
	err := c.graphicsCaptureSession.PutQueryInterface(winrt.IClosableID, &closable)
	if err != nil {
		return errors.Wrap(err, "PutQueryInterface: graphicsCaptureSession")
	}
	defer closable.Release()

	closable.Close()

	c.graphicsCaptureItem = nil

	return nil
}
