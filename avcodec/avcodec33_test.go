package avcodec

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestCodecParameters_FreezeOutlivesSourceContext verifies that a snapshot
// taken with NewCodecParametersFromContext stays valid (and can be applied to a
// fresh context with ToContext) after the source context has been freed. This
// is the property SharedStreamRef relies on to survive demuxer/stream close.
func TestCodecParameters_FreezeOutlivesSourceContext(t *testing.T) {
	src := testNewContextWithCodec(t, "mpeg4")
	src.SetWidth(1280)
	src.SetHeight(720)
	wantCodecID := src.CodecID()

	params, err := NewCodecParametersFromContext(src)
	require.NoError(t, err)
	require.NotNil(t, params)
	defer params.Free()

	// Drop the source context: the frozen snapshot must remain usable.
	src.Free()

	dst := testNewContextWithCodec(t, "mpeg4")
	defer dst.Free()
	require.NoError(t, params.ToContext(dst))

	require.Equal(t, 1280, dst.Width())
	require.Equal(t, 720, dst.Height())
	require.Equal(t, wantCodecID, dst.CodecID())
}

// TestCodecParameters_RoundTripMatchesCopyTo checks the freeze/apply pair
// produces the same result as the existing Context.CopyTo round-trip.
func TestCodecParameters_RoundTripMatchesCopyTo(t *testing.T) {
	src := testNewContextWithCodec(t, "mpeg4")
	defer src.Free()
	src.SetWidth(640)
	src.SetHeight(480)

	viaCopyTo := testNewContextWithCodec(t, "mpeg4")
	defer viaCopyTo.Free()
	require.NoError(t, src.CopyTo(viaCopyTo))

	params, err := NewCodecParametersFromContext(src)
	require.NoError(t, err)
	defer params.Free()
	viaParams := testNewContextWithCodec(t, "mpeg4")
	defer viaParams.Free()
	require.NoError(t, params.ToContext(viaParams))

	require.Equal(t, viaCopyTo.Width(), viaParams.Width())
	require.Equal(t, viaCopyTo.Height(), viaParams.Height())
	require.Equal(t, viaCopyTo.CodecID(), viaParams.CodecID())
}
