package gbimporter

import "go/build"

// PackedContext is a copy of build.Context without the func fields.
//
// TODO(mdempsky): Not sure this belongs here.
type PackedContext struct {
	GOARCH        string
	GOOS          string
	GOROOT        string
	GOPATH        string
	CgoEnabled    bool
	UseAllFiles   bool
	Compiler      string
	BuildTags     []string
	ReleaseTags   []string
	InstallSuffix string
}

func PackContext(ctx *build.Context) PackedContext {
	return PackedContext{
		GOARCH:        ctx.GOARCH,
		GOOS:          ctx.GOOS,
		GOROOT:        ctx.GOROOT,
		GOPATH:        ctx.GOPATH,
		CgoEnabled:    ctx.CgoEnabled,
		UseAllFiles:   ctx.UseAllFiles,
		Compiler:      ctx.Compiler,
		BuildTags:     ctx.BuildTags,
		ReleaseTags:   ctx.ReleaseTags,
		InstallSuffix: ctx.InstallSuffix,
	}
}

func (pc PackedContext) BuildContext() *build.Context {
	return &build.Context{
		GOARCH:        pc.GOARCH,
		GOOS:          pc.GOOS,
		GOROOT:        pc.GOROOT,
		GOPATH:        pc.GOPATH,
		CgoEnabled:    pc.CgoEnabled,
		UseAllFiles:   pc.UseAllFiles,
		Compiler:      pc.Compiler,
		BuildTags:     pc.BuildTags,
		ReleaseTags:   pc.ReleaseTags,
		InstallSuffix: pc.InstallSuffix,
	}
}
