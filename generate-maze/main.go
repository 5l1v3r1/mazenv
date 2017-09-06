package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/unixpickle/essentials"
	"github.com/unixpickle/mazenv"
)

type CommonFlags struct {
	Rows   int
	Cols   int
	Seed   int
	Num    int
	Border bool
}

func (c *CommonFlags) AddFlags(f *flag.FlagSet) {
	f.IntVar(&c.Rows, "rows", 11, "height of maze grid")
	f.IntVar(&c.Cols, "cols", 11, "width of maze grid")
	f.IntVar(&c.Seed, "seed", -1, "random number generator seed (-1 for random)")
	f.IntVar(&c.Num, "num", 1, "number of mazes to generate")
	f.BoolVar(&c.Border, "border", false, "add a border of walls around the maze")
}

type Generator interface {
	mazenv.Generator
	Description() string
	AddFlags(f *flag.FlagSet)
}

var Generators = map[string]Generator{
	"prim":   &mazenv.PrimGenerator{},
	"island": &mazenv.IslandGenerator{},
}

func main() {
	if len(os.Args) < 2 {
		dieUsage()
	}
	algoName, args := os.Args[1], os.Args[2:]
	if algo, ok := Generators[algoName]; ok {
		fs := flag.NewFlagSet(algoName, flag.ExitOnError)
		common := CommonFlags{}
		common.AddFlags(fs)
		algo.AddFlags(fs)
		fs.Parse(args)
		if common.Seed == -1 {
			rand.Seed(time.Now().UnixNano())
		} else {
			rand.Seed(int64(common.Seed))
		}
		for i := 0; i < common.Num; i++ {
			maze, err := algo.Generate(common.Rows, common.Cols)
			if err != nil {
				essentials.Die(err)
			}
			if common.Border {
				maze = maze.Bordered()
			}
			fmt.Println(maze.String())
			if i+1 < common.Num {
				fmt.Println()
			}
		}
	} else {
		fmt.Fprintln(os.Stderr, "unknown algorithm:", algoName)
		dieUsage()
	}
}

func dieUsage() {
	lines := []string{
		"Usage: generate-maze <generator> [flags | -help]",
		"",
		"Available generators:",
		"",
	}

	var longestName int
	var names []string
	for name := range Generators {
		names = append(names, name)
		longestName = essentials.MaxInt(longestName, len(name))
	}
	sort.Strings(names)
	for _, name := range names {
		desc := Generators[name].Description()
		for len(name) < longestName {
			name += " "
		}
		lines = append(lines, " "+name+"    "+desc)
	}

	lines = append(lines, "")
	for _, line := range lines {
		fmt.Fprintln(os.Stderr, line)
	}
	os.Exit(1)
}
