package main

import (
	"bytes"
	"fmt"
	"go/types"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/mdempsky/gocode/gbimporter"
	"github.com/mdempsky/gocode/suggest"
	"github.com/mdempsky/gocode/vgo"
)

func doServer() {
	addr := *g_addr
	if *g_sock == "unix" {
		addr = getSocketPath()
	}

	lis, err := net.Listen(*g_sock, addr)
	if err != nil {
		log.Fatal(err)
	}

	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		exitServer()
	}()

	if err = rpc.Register(&Server{}); err != nil {
		log.Fatal(err)
	}
	rpc.Accept(lis)
}

func exitServer() {
	if *g_sock == "unix" {
		_ = os.Remove(getSocketPath())
	}
	os.Exit(0)
}

type Server struct {
}

type AutoCompleteRequest struct {
	Filename string
	Data     []byte
	Cursor   int
	Context  gbimporter.PackedContext
}

type AutoCompleteReply struct {
	Candidates []suggest.Candidate
	Len        int
}

func (s *Server) AutoComplete(req *AutoCompleteRequest, res *AutoCompleteReply) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic: %s\n\n", err)
			debug.PrintStack()

			res.Candidates = []suggest.Candidate{
				{Class: "PANIC", Name: "PANIC", Type: "PANIC"},
			}
		}
	}()
	var importer types.Importer
	vgoMode := false

	if !*g_vgo {
		// client sends an abs path
		dir := filepath.Dir(req.Filename)
		found := false
		for {
			if fi, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil && !fi.IsDir() {
				found = true
				vgoMode = true
				break
			}
			d, _ := filepath.Split(dir)
			if d == dir {
				break
			}
			dir = d
		}
		if !found {
			importer = gbimporter.New(&req.Context, req.Filename)
		}
	}

	var vgoOutput *strings.Builder

	if importer == nil {
		v := vgo.NewTestLoader(filepath.Dir(req.Filename))
		vgoOutput = new(strings.Builder)
		v.Debug = vgoOutput
		importer = v
	}

	if *g_debug {
		var buf bytes.Buffer
		log.Printf("Got autocompletion request for '%s'\n", req.Filename)
		log.Printf("vgo-mode: %v\n", vgoMode)
		log.Printf("Cursor at: %d\n", req.Cursor)
		buf.WriteString("-------------------------------------------------------\n")
		buf.Write(req.Data[:req.Cursor])
		buf.WriteString("#")
		buf.Write(req.Data[req.Cursor:])
		log.Print(buf.String())
		log.Println("-------------------------------------------------------")
	}
	now := time.Now()

	cfg := suggest.Config{
		Importer: importer,
	}
	if *g_debug {
		cfg.Logf = log.Printf
	}
	candidates, d := cfg.Suggest(req.Filename, req.Data, req.Cursor)
	elapsed := time.Since(now)
	if *g_debug {
		if vgoOutput != nil {
			log.Printf("vgo output was: \n%v\n", vgoOutput.String())
		}
		log.Printf("Elapsed duration: %v\n", elapsed)
		log.Printf("Offset: %d\n", res.Len)
		log.Printf("Number of candidates found: %d\n", len(candidates))
		log.Printf("Candidates are:\n")
		for _, c := range candidates {
			log.Printf("  %s\n", c.String())
		}
		log.Println("=======================================================")
	}
	res.Candidates, res.Len = candidates, d
	return nil
}

type ExitRequest struct{}
type ExitReply struct{}

func (s *Server) Exit(req *ExitRequest, res *ExitReply) error {
	go func() {
		time.Sleep(time.Second)
		exitServer()
	}()
	return nil
}
