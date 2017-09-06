package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/unixpickle/essentials"
	"github.com/unixpickle/mazenv"
)

func main() {
	var mazesPath string
	var lengthOnly bool
	flag.StringVar(&mazesPath, "in", "", "file containing mazes (instead of stdin)")
	flag.BoolVar(&lengthOnly, "length", false, "only print the solution length")
	flag.Parse()

	for maze := range readMazes(mazesPath) {
		solution := mazenv.Solve(maze)
		if lengthOnly {
			fmt.Println(len(solution))
		} else {
			fmt.Println(solution)
		}
	}
}

func readMazes(path string) <-chan *mazenv.Maze {
	res := make(chan *mazenv.Maze, 1)

	go func() {
		defer close(res)

		reader, err := mazeReader(path)
		essentials.Must(err)
		defer reader.Close()

		br := bufio.NewReader(reader)

		var curMaze string
		sendCurMaze := func() {
			if len(curMaze) == 0 {
				return
			}
			maze, err := mazenv.ParseMaze(curMaze)
			essentials.Must(err)
			res <- maze
			curMaze = ""
		}

		for {
			line, err := br.ReadString('\n')
			if err == io.EOF {
				break
			}
			essentials.Must(err)
			line = strings.TrimSpace(line)
			if len(line) == 0 {
				sendCurMaze()
			} else {
				curMaze += line + "\n"
			}
		}

		sendCurMaze()
	}()

	return res
}

func mazeReader(path string) (io.ReadCloser, error) {
	if path == "" {
		fmt.Fprintln(os.Stderr, "reading from standard input...")
		return os.Stdin, nil
	} else {
		return os.Open(path)
	}
}
