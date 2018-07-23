// +build !wasm

package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/dennwc/webidl/ast"
	"github.com/dennwc/webidl/parser"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

var (
	f_pkg = flag.String("p", "dom", "package name")
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatal("expected at least one file")
	}
	for _, arg := range flag.Args() {
		if err := process(arg); err != nil {
			log.Fatalln(arg, ":", err)
		}
	}
}

func process(fname string) error {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}
	root := parser.Parse(string(data))
	data = nil
	name := filepath.Base(fname)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	name = strings.NewReplacer(
		" - ", "-",
		" ", "_",
	).Replace(strings.ToLower(name))
	f, err := os.Create(name + ".go")
	if err != nil {
		return err
	}
	defer f.Close()

	if err := generate(f, root); err != nil {
		return err
	}
	return f.Close()
}

func generate(w io.Writer, root *ast.File) error {
	fmt.Fprintln(w, `package`, *f_pkg)
	fmt.Fprintln(w)
	for _, d := range root.Declarations {
		if err := generateDecl(w, d); err != nil {
			return err
		}
	}
	return nil
}

func generateDecl(w io.Writer, d ast.Decl) error {
	switch d := d.(type) {
	case *ast.Implementation:
		fmt.Fprintf(w, "var _ %s = (%s)(nil)\n\n", convName(d.Source), convName(d.Name))
	case *ast.Typedef:
		fmt.Fprintf(w, "type %s ", convName(d.Name))
		s, err := generateType(w, d.Type, convName(d.Name))
		if err != nil {
			return err
		}
		fmt.Fprint(w, "\n\n")
		w.Write([]byte(s))
	case *ast.Enum:
		name := convName(d.Name)
		fmt.Fprintf(w, "type %s string\n\n", name)
		fmt.Fprintf(w, "const (\n")
		for i, v := range d.Values {
			fmt.Fprintf(w, "\t%s_%d = %s(%v)\n", name, i+1, name, printLit(v))
		}
		fmt.Fprint(w, ")\n\n")
	case *ast.Dictionary:
		fmt.Fprintf(w, "type %s struct{\n", d.Name)
		if d.Inherits != "" {
			fmt.Fprintf(w, "\t%s\n", d.Inherits)
		}
		buf := bytes.NewBuffer(nil)
		for _, m := range d.Members {
			name := convName(m.Name)
			fmt.Fprintf(w, "\t%s ", name)
			s, err := generateType(w, m.Type, name)
			if err != nil {
				return err
			}
			buf.WriteString(s)
			fmt.Fprintln(w)
		}
		fmt.Fprintf(w, "}\n\n")
		w.Write(buf.Bytes())
	case *ast.Interface:
		fmt.Fprintf(w, "type %s interface{\n", d.Name)
		if d.Inherits != "" {
			fmt.Fprintf(w, "\t%s\n", d.Inherits)
		}
		buf := bytes.NewBuffer(nil)
		consts := 0
		for _, m := range d.Members {
			switch m := m.(type) {
			case *ast.Member:
				if m.Const {
					consts++
					continue
				}
				name := convName(m.Name)
				if m.Attribute {
					// getter
					fmt.Fprintf(w, "\t%s() ", name)
					s, err := generateType(w, m.Type, name)
					if err != nil {
						return err
					}
					buf.WriteString(s)
					fmt.Fprintln(w)
					if !m.Readonly {
						// setter
						fmt.Fprintf(w, "\tSet%s(v ", name)
						s, err = generateType(w, m.Type, name)
						if err != nil {
							return err
						}
						buf.WriteString(s)
						fmt.Fprintf(w, ")\n")
					}
				} else {
					fmt.Fprintf(w, "\t%s(", name)
					s, err := generateParams(w, m.Parameters, name)
					if err != nil {
						return err
					}
					buf.WriteString(s)
					fmt.Fprintf(w, ") ")
					s, err = generateType(w, m.Type, name)
					if err != nil {
						return err
					}
					buf.WriteString(s)
					fmt.Fprintln(w)
				}
			}
		}
		fmt.Fprintf(w, "}\n\n")
		if consts != 0 {
			fmt.Fprint(w, "const (\n")
			for _, m := range d.Members {
				switch m := m.(type) {
				case *ast.Member:
					if !m.Const {
						continue
					}
					fmt.Fprintf(w, "\t%s = %v\n", constName(d.Name, m.Name), printLit(m.Init))
				}
			}
			fmt.Fprint(w, ")\n\n")
		}
		w.Write(buf.Bytes())
	default:
		log.Printf("skip declaration: %T", d)
	}
	return nil
}

func generateParams(w io.Writer, params []*ast.Parameter, nameHint string) (string, error) {
	buf := bytes.NewBuffer(nil)
	for i, p := range params {
		if i != 0 {
			fmt.Fprint(w, ", ")
		}
		fmt.Fprint(w, localName(p.Name), " ")
		if p.Variadic {
			fmt.Fprint(w, "...")
		}
		s, err := generateType(w, p.Type, nameHint+convName(p.Name))
		if err != nil {
			return "", err
		}
		buf.WriteString(s)
	}
	return buf.String(), nil
}

var typeNames = map[string]string{
	"void":           "",
	"boolean":        "bool",
	"DOMString":      "string",
	"USVString":      "string",
	"unsigned short": "uint16",
	"short":          "int16",
	"unsigned long":  "uint32",
	"long":           "int32",
	"double":         "float64",
}

func generateType(w io.Writer, t ast.Type, nameHint string) (string, error) {
	switch t := t.(type) {
	case *ast.AnyType:
		fmt.Fprint(w, "interface{}")
	case *ast.TypeName:
		name := t.Name
		if s, ok := typeNames[name]; ok {
			name = s
		} else {
			name = convName(name)
		}
		if name != "" {
			fmt.Fprint(w, name)
		}
	case *ast.NullableType:
		//fmt.Fprint(w, "*")
		return generateType(w, t.Type, nameHint)
	case *ast.SequenceType:
		fmt.Fprint(w, "[]")
		return generateType(w, t.Elem, nameHint)
	case *ast.UnionType:
		fmt.Fprintf(w, "interface{ is%s() }", nameHint)
		buf := bytes.NewBuffer(nil)
		for i, e := range t.Types {
			tname := fmt.Sprintf("%s%d", nameHint, i+1)
			fmt.Fprintf(buf, "type %s struct{\n\tValue ", tname)
			s, err := generateType(buf, e, nameHint+"_")
			if err != nil {
				return "", err
			}
			fmt.Fprintf(buf, "\n}\n")
			fmt.Fprint(buf, s)
			fmt.Fprintf(buf, "func (%s) is%s(){}\n\n", tname, nameHint)
		}
		return buf.String(), nil
	default:
		log.Printf("skip type: %T", t)
	}
	return "", nil
}

func convName(name string) string {
	r := []rune(name)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func constName(typ, name string) string {
	return convName(typ) + "_" + name
}

func localName(s string) string {
	switch s {
	case "type":
		return "typ"
	case "interface":
		return "iface"
	}
	return s
}

func printLit(l ast.Literal) string {
	val := l.(*ast.BasicLiteral).Value
	if val == "null" {
		return "nil"
	}
	return val
}
