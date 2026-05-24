package libav

import (
	"github.com/SpalkLtd/go-libav/avfilter"
	"github.com/SpalkLtd/go-libav/avutil"
)

// DefaultFactory returns a value whose NewGraph and NewDictionary methods wrap
// avfilter.NewGraph and avutil.NewDictionary. Consumers define their own
// factory interface and accept this value through it.
func DefaultFactory() Factory { return Factory{} }

// Factory is the concrete type returned by DefaultFactory.
type Factory struct{}

func (Factory) NewGraph() avfilter.IGraph        { return avfilter.NewGraph() }
func (Factory) NewDictionary() avutil.IDictionary { return avutil.NewDictionary() }
