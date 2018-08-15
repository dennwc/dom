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
	sp, err := NewSplitter(bin)
	if err != nil {
		return err
	}
	m, err := sp.SplitByPrefix(prefixes)
	if err != nil {
		return err
	}
	_ = m
	return nil
}

func NewSplitter(mod *wasm.Module) (*Splitter, error) {
	sp := &Splitter{mod: mod}
	if err := sp.decodeNames(); err != nil {
		return nil, err
	}
	sp.countImported()
	if err := sp.buildFuncTable(); err != nil {
		return nil, err
	}
	sp.statsCallIndirect()
	return sp, nil
}

type Splitter struct {
	mod       *wasm.Module
	funcs     wasm.NameMap // function names; indexes are in a function index space (with funcsImp offset)
	funcsImp  int          // number of imported functions
	funcTable []int        // global table with function indexes

	bodies map[int][]disasm.Instr
}

func findInstr(code []disasm.Instr, typ byte) int {
	for i, op := range code {
		if op.Op.Code == typ {
			return i
		}
	}
	return -1
}

func findInstrRev(code []disasm.Instr, typ byte) int {
	for i := len(code) - 1; i >= 0; i-- {
		op := code[i]
		if op.Op.Code == typ {
			return i
		}
	}
	return -1
}

func sumCallStubs(instr []disasm.Instr) int {
	var sum int
	for {
		// find call_indirect
		ci := findInstr(instr, operators.CallIndirect)
		if ci < 0 {
			ci = findInstr(instr, operators.Call)
		}
		if ci < 0 {
			return sum
		}
		// find SP -= 8
		sp := findInstrRev(instr[:ci], operators.I32Sub)
		if sp < 0 || instr[sp-1].Op.Code != operators.I32Const || instr[sp-1].Immediates[0].(int32) != 8 {
			instr = instr[ci+1:]
			continue
		}
		sp -= 2
		expr := instr[sp : ci+1]

		data, err := disasm.Assemble(expr)
		if err != nil {
			panic(err)
		}
		sum += len(data) - 2
		instr = instr[ci+1:]
	}
}

func sumReturnStubs(instr []disasm.Instr) int {
	var sum int
	codes := []byte{
		operators.I32Const,

		operators.SetGlobal,
		operators.I32Add,
		operators.I32Const,
		operators.GetGlobal,

		operators.SetGlobal,
		operators.I32Load16u,
		operators.GetGlobal,

		operators.SetGlobal,
		operators.I32Load16u,
		operators.GetGlobal,
	}

	opt := []byte{
		operators.SetGlobal,
		operators.I32Add,
		operators.I32Const,
		operators.GetGlobal,
	}
loop:
	for {
		// find return
		ci := findInstr(instr, operators.Return)
		if ci < 0 {
			return sum
		}
		next := func() {
			instr = instr[ci+1:]
		}
		// check previous ops
		for i, opc := range codes {
			op := instr[ci-1-i]
			if op.Op.Code != opc {
				next()
				continue loop
			}
		}
		end := ci - len(codes)
		ok := true
		for i, opc := range opt {
			ind := end - 1 - i
			if ind < 0 {
				ok = false
				break
			}
			op := instr[ind]
			if op.Op.Code != opc {
				ok = false
				break
			}
		}
		if ok {
			end -= len(opt)
		}

		expr := instr[end : ci+1]

		data, err := disasm.Assemble(expr)
		if err != nil {
			panic(err)
		}
		sum += len(data) - 2
		next()
	}
}

func (sp *Splitter) statsCallIndirect() {
	save := 0
	for ind, name := range sp.funcs {
		instr, err := sp.disassemble(int(ind))
		if err != nil {
			log.Println(name, err)
			continue // we only gathering stats here
		}
		save += sumCallStubs(instr)
		save += sumReturnStubs(instr)
	}
	log.Println("wrapping to a function will save:", save, humanize.Bytes(uint64(save)))
}

func (sp *Splitter) decodeNames() error {
	sec := sp.mod.Custom(wasm.CustomSectionName)
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
	sp.funcs = sub.(*wasm.FunctionNames).Names
	return nil
}

func (sp *Splitter) countImported() {
	if sp.mod.Import == nil {
		return
	}
	var n int
	for _, imp := range sp.mod.Import.Entries {
		if _, ok := imp.Type.(wasm.FuncImport); ok {
			n++
		}
	}
	sp.funcsImp = n
}

func (sp *Splitter) buildFuncTable() error {
	var (
		tbl *wasm.Table
		ind int
	)
	for i, t := range sp.mod.Table.Entries {
		if t.ElementType == wasm.ElemTypeAnyFunc {
			ind, tbl = i, &t
			break
		}
	}
	if tbl == nil {
		return nil
	}
	table := make([]int, tbl.Limits.Initial)
	for _, e := range sp.mod.Elements.Entries {
		if int(e.Index) != ind {
			continue
		}
		stack, err := evalCode(e.Offset)
		if err != nil {
			return fmt.Errorf("cannot evaluate table offset: %v", err)
		}
		off := stack[0]
		for i, v := range e.Elems {
			table[int(off)+i] = int(v)
		}
	}
	sp.funcTable = table
	return nil
}

func (sp *Splitter) importFuncName(i int) string {
	ind := 0
	for _, e := range sp.mod.Import.Entries {
		_, ok := e.Type.(wasm.FuncImport)
		if !ok {
			continue
		}
		if ind == i {
			return e.ModuleName + "." + e.FieldName
		}
		ind++
	}
	return ""
}

func (sp *Splitter) funcName(i int) string {
	name, ok := sp.funcs[uint32(i)]
	if !ok && sp.isImported(i) {
		name = sp.importFuncName(i)
	}
	return name
}

func (sp *Splitter) funcNameRel(i int) string {
	i += sp.funcsImp
	return sp.funcName(i)
}

func (sp *Splitter) lookupFuncTable(i int) int {
	return sp.funcTable[i]
}

func (sp *Splitter) SplitByPrefix(prefixes []string) (*wasm.Module, error) {
	// indexes are in a function index space
	var funcs []int
	for ind, name := range sp.funcs {
		if hasAnyPrefix(name, prefixes) {
			funcs = append(funcs, int(ind))
		}
	}
	return sp.SplitFunctions(funcs)
}

func (sp *Splitter) isImported(fnc int) bool {
	return fnc < sp.funcsImp
}
func (sp *Splitter) toFuncTable(fnc int) int {
	return fnc - sp.funcsImp
}
func (sp *Splitter) toFuncSpace(fnc int) int {
	return fnc + sp.funcsImp
}
func (sp *Splitter) codeOf(fnc int) ([]byte, error) {
	if sp.isImported(fnc) {
		return nil, fmt.Errorf("attempting to disassemble imported function")
	}
	fnc = sp.toFuncTable(fnc)
	if fnc >= len(sp.mod.Code.Bodies) {
		return nil, fmt.Errorf("function index out of bounds")
	}
	b := sp.mod.Code.Bodies[fnc]
	return b.Code, nil
}
func (sp *Splitter) disassemble(fnc int) ([]disasm.Instr, error) {
	if instr, ok := sp.bodies[fnc]; ok {
		return instr, nil
	}
	code, err := sp.codeOf(fnc)
	if err != nil {
		return nil, err
	}
	d, err := disasm.DisassembleRaw(code)
	if err != nil {
		return nil, err
	}
	if sp.bodies == nil {
		sp.bodies = make(map[int][]disasm.Instr)
	}
	sp.bodies[fnc] = d
	return d, nil
}

func (sp *Splitter) SplitFunctions(funcs []int) (*wasm.Module, error) {
	job := sp.newSplitJob(funcs)
	if err := job.ValidateFuncs(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("TODO: split")
}

func (sp *Splitter) newSplitJob(funcs []int) *splitJob {
	splitFunc := make(map[int]struct{}, len(funcs))
	for _, ind := range funcs {
		splitFunc[ind] = struct{}{}
	}
	return &splitJob{sp: sp, split: splitFunc, funcs: funcs}
}

type splitJob struct {
	sp    *Splitter
	split map[int]struct{}
	funcs []int
}

func (s *splitJob) ValidateFuncs() error {
	for _, fnc := range s.funcs {
		if s.sp.isImported(fnc) {
			return fmt.Errorf("attempting to split imported function")
		}
		d, err := s.sp.disassemble(fnc)
		if err != nil {
			return err
		}
		for oi, op := range d {
			var callee int
			switch op.Op.Code {
			case operators.Call:
				callee = int(op.Immediates[0].(uint32))
			case operators.CallIndirect:
				n, err := backEvalN(d[:oi], 1)
				if err != nil {
					return fmt.Errorf("cannot eval an indirrect call target from '%v': %v", s.sp.funcName(fnc), err)
				} else if n < 0 {
					return fmt.Errorf("cannot eval an indirrect call target from '%v'", s.sp.funcName(fnc))
				}
				expr := d[n:oi]
				log.Println("indirect call:", s.sp.funcName(fnc), expr)
				stack, err := eval(expr)
				if err != nil {
					return fmt.Errorf("cannot eval an indirrect call target from '%v': %v", s.sp.funcName(fnc), err)
				}
				ind := int(stack[0])
				callee = s.sp.lookupFuncTable(ind)
				log.Printf("indirect call: '%v' -> '%v' (%d = %d)",
					s.sp.funcName(fnc), s.sp.funcName(callee), ind, callee)
			default:
				continue
			}
			if s.sp.isImported(callee) {
				continue // call of imported function
			}
			if _, ok := s.split[callee]; !ok {
				return fmt.Errorf("cannot split: external call to '%v'", s.sp.funcName(callee))
			}
		}
	}
	return nil
}

var stackVars = map[string]int{
	"i32.shr_u":    2,
	"i32.const":    0,
	"i64.const":    0,
	"i32.store":    1,
	"i64.store":    1,
	"i32.wrap/i64": 1,
	"set_global":   1,
}

func backEvalN(instr []disasm.Instr, req int) (int, error) {
	i := len(instr) - 1
	for ; i >= 0; i-- {
		op := instr[i]
		log.Println(req, op.Immediates, op.Op.Name, op.Op.Args, op.Op.Returns)
		st, ok := stackVars[op.Op.Name]
		if !ok {
			return 0, fmt.Errorf("unsupported op: %v", op.Op.Name)
		}
		req += st
		if op.Op.Returns != 0 && op.Op.Returns != wasm.ValueType(wasm.BlockTypeEmpty) {
			req--
		}
		if req == 0 {
			return i, nil
		}
	}
	return -1, nil
}

func eval(instr []disasm.Instr) ([]uint64, error) {
	var stack []uint64
	push := func(v uint64) {
		stack = append(stack, v)
	}
	pop := func() uint64 {
		i := len(stack) - 1
		v := stack[i]
		stack = stack[:i]
		return v
	}
	for i, op := range instr {
		switch op.Op.Code {
		case operators.I32Const:
			v := op.Immediates[0].(int32)
			push(uint64(v))
		case operators.I32WrapI64:
			push(uint64(uint32(pop())))
		case operators.I32ShrU:
			v2 := uint32(pop())
			v1 := uint32(pop())
			push(uint64(v1 >> v2))
		case operators.SetGlobal:
			_ = pop()
			// do nothing
		case operators.End:
			if i != len(instr)-1 {
				return nil, fmt.Errorf("unexpected end")
			}
		default:
			return nil, fmt.Errorf("unsupported eval operation: %v", op.Op.Name)
		}
	}
	return stack, nil
}

func evalCode(code []byte) ([]uint64, error) {
	instr, err := disasm.DisassembleRaw(code)
	if err != nil {
		return nil, err
	}
	return eval(instr)
}

func hasAnyPrefix(s string, prefixes []string) bool {
	for _, pref := range prefixes {
		if strings.HasPrefix(s, pref) {
			return true
		}
	}
	return false
}
