// Package libav exposes a cross-cutting AVFactory interface that bundles the
// libav constructors most users need: filter graphs and option dictionaries.
//
// The default factory wraps the real libav constructors. Tests substitute fakes
// to drive error paths without calling into FFmpeg.
package libav

import (
	"github.com/SpalkLtd/go-libav/avfilter"
	"github.com/SpalkLtd/go-libav/avutil"
)

// AVFactory creates the libav primitives used by code that builds filter graphs.
type AVFactory interface {
	NewGraph() avfilter.IGraph
	NewDictionary() avutil.IDictionary
}

// DefaultFactory returns an AVFactory that wraps the real libav constructors.
func DefaultFactory() AVFactory { return defaultFactory{} }

type defaultFactory struct{}

func (defaultFactory) NewGraph() avfilter.IGraph        { return avfilter.NewGraph() }
func (defaultFactory) NewDictionary() avutil.IDictionary { return avutil.NewDictionary() }
