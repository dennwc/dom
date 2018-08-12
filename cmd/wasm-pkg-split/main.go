package main

import (
	"flag"
	"log"
	"os"

	"bytes"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/go-interpreter/wagon/disasm"
	"github.com/go-interpreter/wagon/wasm"
	"github.com/go-interpreter/wagon/wasm/operators"
	"path/filepath"
	"strings"
)

func main() {
	flag.Parse()
	if err := split(flag.Arg(0)); err != nil {
		log.Fatal(err)
	}
}

func split(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	m, err := wasm.DecodeModule(f)
	f.Close()
	if err != nil {
		return fmt.Errorf("cannot decode module: %v", err)
	}

	if err := splitPackages(m, filepath.Dir(path), []string{
		"runtime.",
		"runtime_",
		"callRet",
		"memeqbody", "cmpbody", "memcmp", "memchr",
		"time.now",
		"sync.event",
		"internal_bytealg",
		"internal_cpu",
	}); err != nil {
		return err
	}

	ext := filepath.Ext(path)
	f, err = os.Create(strings.TrimSuffix(path, ext) + "_out" + ext)
	if err != nil {
		return err
	}
	defer f.Close()

	return wasm.EncodeModule(f, m)
}

func splitPackages(bin *wasm.Module, dir string, prefixes []string) error {
	fmt.Println("section sizes:")
	for _, s := range bin.Sections {
		raw := s.GetRawSection()
		size := len(raw.Bytes)
		fmt.Printf("%10v:  %8v\n", raw.ID, humanize.Bytes(uint64(size)))
	}
	sec := bin.Custom(wasm.CustomSectionName)
	if sec == nil {
		return fmt.Errorf("cannot find names section")
	}
	var names wasm.NameSection
	if err := names.UnmarshalWASM(bytes.NewReader(sec.Data)); err != nil {
		return fmt.Errorf("cannot decode names section: %v", err)
	}
	sub, err := names.Decode(wasm.NameFunction)
	if err != nil {
		return err
	} else if sub == nil {
		return fmt.Errorf("no function names")
	}
	funcs := sub.(*wasm.FunctionNames)

	fimp := uint32(0)
	if bin.Import != nil {
		for _, imp := range bin.Import.Entries {
			if _, ok := imp.Type.(wasm.FuncImport); ok {
				fimp++
			}
		}
	}

	toMove := make(map[uint32]struct{})

	for ind, name := range funcs.Names {
		ind -= fimp
		if false {
			ok := true
			for _, pref := range prefixes {
				if strings.HasPrefix(name, pref) {
					ok = false
					break
				}
			}
			if !ok {
				continue
			}
		} else {
			ok := false
			for _, pref := range prefixes {
				if strings.HasPrefix(name, pref) {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}
		toMove[ind] = struct{}{}
		//log.Println(name)
	}
	fmt.Printf("total %d functions\n", len(toMove))
	var (
		direct, indirect int
		external         int
		inPkg            int
	)

	for ind := range toMove {
		if int(ind) >= len(bin.Code.Bodies) {
			continue
		}
		b := bin.Code.Bodies[ind]
		d, err := disasm.DisassembleRaw(b.Code)
		if err != nil {
			return err
		}
		for _, op := range d {
			switch op.Op.Code {
			case operators.Call:
				direct++
				callee := op.Immediates[0].(uint32)
				if callee < fimp {
					continue // call of imported function
				}
				if _, ok := toMove[callee-fimp]; !ok {
					log.Printf("external call: %v", funcs.Names[callee])
					external++
				}
			case operators.CallIndirect:
				indirect++
				//tind := op.Immediates[0].(uint32)
			}
		}
		inPkg += len(b.Code)
	}
	if external != 0 {
		return fmt.Errorf("cannot split: %d external calls", external)
	}
	log.Printf("direct: %v, indirrect: %v, external: %v\n", direct, indirect, external)

	var (
		total int
	)
	for _, b := range bin.Code.Bodies {
		total += len(b.Code)
	}
	log.Printf("will split %v/%v", humanize.Bytes(uint64(inPkg)), humanize.Bytes(uint64(total)))
	return nil
}
