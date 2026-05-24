package libav

import (
	"github.com/SpalkLtd/go-libav/avfilter"
	"github.com/SpalkLtd/go-libav/avutil"
)

func DefaultFactory() Factory { return Factory{} }

type Factory struct{}

func (Factory) NewGraph() avfilter.IGraph        { return avfilter.NewGraph() }
func (Factory) NewDictionary() avutil.IDictionary { return avutil.NewDictionary() }
