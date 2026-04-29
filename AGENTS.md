# go-libav Agent Reference

## Background

Go language bindings for FFmpeg libraries (libavutil, libavcodec, libavformat,
libavfilter) via CGo. Used exclusively by the synchroniser service for media
processing. Targets the Spalk custom FFmpeg 4.3.1 fork at
`libs/ffmpeg/spalk-ffmpeg/` and supports multiple FFmpeg versions via build tags
(`ffmpeg30`, `ffmpeg33`, `ffmpeg43`). Only `ffmpeg43` is used in production.

Architecture doc: `docs/libs/go-libav.md`.

## Repository Layout

```
go-libav/
├── avcodec/                   # Bindings for libavcodec (encoding/decoding)
│   ├── avcodec.go
│   └── avcodec_test.go
├── avfilter/                  # Bindings for libavfilter (filter graphs)
│   ├── avfilter.go
│   └── avfilter_test.go
├── avformat/                  # Bindings for libavformat (muxing/demuxing)
│   ├── avformat.go
│   └── avformat_test.go
├── avutil/                    # Bindings for libavutil (utilities, pixel formats)
│   ├── avutil.go
│   └── avutil_test.go
├── Makefile                   # Build targets: gofmt, golint, govet, test, cover
├── go.mod                     # Go 1.22.1, module: github.com/SpalkLtd/go-libav
├── CHANGELOG.md
└── LICENSE
```

## Development Commands

```bash
# Run all tests (requires spalk-ffmpeg dev libraries installed)
cd libs/go/go-libav && go test -tags=ffmpeg43 ./...

# Run tests with race detector
cd libs/go/go-libav && go test -tags=ffmpeg43 -race ./...

# Run tests with coverage
cd libs/go/go-libav && make FFMPEG_TAG=ffmpeg43 cover-test
```

## Testing Guidance

- All tests require the `-tags=ffmpeg43` build tag.
- The Spalk custom FFmpeg libraries must be installed; standard system packages are insufficient.
- Tests use `testify` for assertions.

## Code Style

- Follow the shared Go style guide at `docs/onboarding/style/golang-style.md`.
- CGo bindings wrap C structs with explicit allocate/free pairs for memory management.
- Each FFmpeg library gets its own Go package mirroring the C library structure.

## Additional Notes

- Will not compile without the Spalk custom FFmpeg 4.3.1 fork installed from `libs/ffmpeg/spalk-ffmpeg/`.
- Module path is `github.com/SpalkLtd/go-libav` (pre-monorepo); do not change without updating synchroniser imports.
- CGo makes builds slower and cross-compilation harder; C-side memory leaks escape Go's garbage collector.
- The `ffmpeg30` and `ffmpeg33` tags exist for history but are not used in production.
