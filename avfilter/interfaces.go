package avfilter

import (
	"unsafe"

	"github.com/SpalkLtd/go-libav/avutil"
)

// IContext is the interface a filter context exposes for setup, linking, and
// frame in/out. *Context implements it directly. Tests may substitute fakes.
type IContext interface {
	Name() string
	Init() error
	InitWithDictionary(opts avutil.IDictionary) error
	Link(srcPad uint, dst IContext, dstPad uint) error
	SendCommand(cmd, args string, returnLength int) (string, error)
	SetOption(name, value string) error
	SetChannelLayoutOption(layout avutil.ChannelLayout) error
	SetSampleFormatOption(name string, format avutil.SampleFormat) error
	SetRationalOption(name string, val *avutil.Rational) error
	SetInt64OptionC(name unsafe.Pointer, val int64) error
	AddFrameWithFlags(frame *avutil.Frame, flags BufferSrcFlags) error
	GetFrameWithFlags(frame *avutil.Frame, flags BufferSinkFlags) error
	GetFrame(frame *avutil.Frame) error
}

// IGraph is the interface a filter graph exposes for filter creation and
// lifecycle. *Graph implements it directly.
type IGraph interface {
	AddFilter(name, id string) (IContext, error)
	Config() error
	RequestOldest() error
	Free()
}
