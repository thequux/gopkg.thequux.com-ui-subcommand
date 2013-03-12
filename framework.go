package subcommand

import (
	"flag"
	"fmt"
	"os"
	"sort"
)

var (
	registry =  map[string]Module{}
)

type Module struct {
	Flags *flag.FlagSet
	Main func()
	Name string
	Descr string
}

func usage() {
	fmt.Printf("Usage:\n  %s [global options] subcommand [subcommand arguments]\n\n", os.Args[0])

	
	modules := make([]string, 0, len(registry))
	for name, _ := range registry {
		modules = append(modules, name)
	}
	sort.StringSlice(modules).Sort()
	fmt.Printf("Available commands:\n")
	for _, name := range modules {
		fmt.Printf("\t%s\t%s\n", name, registry[name].Descr)
	}

	fmt.Printf("\nGlobal options:\n")
	flag.PrintDefaults()
}

func init() {
	flag.Usage = usage
}

func Register(m Module) {
	registry[m.Name] = m
}

func Subcommand(cmd string) *Module {
	if mod, exist := registry[cmd]; exist {
		return &mod
	}
	return nil
	
}

// Convenience function that can be called directly from main
func Main() {
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	submod := Subcommand(flag.Arg(0))
	if submod == nil {
		fmt.Printf("%s %s: subcommand not found\n", os.Args[0], flag.Arg(0))
		flag.Usage()
		return
	}

	submod.Flags.Parse(flag.Args()[1:])
	submod.Main()
}
