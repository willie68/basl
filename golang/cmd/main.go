package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

type Stack []int

const COMMANDS string = " 0123456789bcdhijkopqrst\"/%&={}+?*~-_#:;.,^|><'"
const USRCMD string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

var (
	baslFile     string
	fs           *flag.FlagSet
	stack        Stack
	store        []int
	loopValue    int
	inDefinition bool
	definitions  map[string]string
	defName      string
	defText      string
	reader       *bufio.Reader
	v            int
	inNumber     bool
)

func init() {
	// variables for parameter override
	fs = flag.NewFlagSet("main", flag.ContinueOnError)
	fs.StringVarP(&baslFile, "file", "f", "", "source file to compile")
	fs.SortFlags = false
}

func main() {
	fs.Parse(os.Args[1:])

	if baslFile != "" {
		f, err := os.Open(baslFile)
		if err != nil {
			panic(fmt.Sprintf("file not found: %s", baslFile))
		}
		defer f.Close()
		reader = bufio.NewReader(f)
	} else {
		reader = bufio.NewReader(os.Stdin)
	}

	stack = make([]int, 0)
	store = make([]int, 1024)
	definitions = make(map[string]string)
	v = 0
	inNumber = false
	fmt.Println("SPLRC Serial Programming Language for Micro Controller")
	for {
		fmt.Print(":")
		text, err := reader.ReadString('\n')
		if (err == io.EOF) && (baslFile != "") {
			reader = bufio.NewReader(os.Stdin)
			text, _ = reader.ReadString('\n')
		}
		// convert CRLF to LF
		text = strings.Replace(text, "\r", "", -1)
		text = strings.Replace(text, "\n", "", -1)
		execute(text)
	}
}

func execute(text string) {
	// slc := strings.Split(text, " ")
	for _, chr := range text {
		schr := string(chr)

		if chr == ':' {
			fmt.Println("start defining user cmd")
			inDefinition = true
			continue
		}

		if inDefinition {
			if chr == ';' {
				fmt.Println("stop defining user cmd")
				inDefinition = false
				definitions[defName] = strings.TrimSpace(defText)
				defName = ""
				defText = ""
				continue
			}
			if defName == "" {
				defName = schr
				continue
			}
			if defText == "" {
				defText = schr
				continue
			}
			defText = defText + schr
			continue
		}

		if strings.Contains(COMMANDS, schr) {
			if (chr >= '0') && (chr <= '9') {
				inNumber = true
				v = v*10 + (int(chr) - int('0'))
				continue
			} else {
				if inNumber {
					stack.Push(v)
					v = 0
					inNumber = false
				}
			}
			if chr == ' ' {
				continue
			}

			processNme(schr)
		}

		if strings.Contains(USRCMD, schr) {
			fmt.Println("user command: ", schr)
			def, ok := definitions[schr]
			if !ok {
				fmt.Println("user command not found! ", schr)
				continue
			}
			execute(def)
		}
	}
}

func processNme(nme string) {
	switch nme {
	case "h":
		showHelp()
	case ".":
		fmt.Printf("stacksize: %d\r\n", len(stack))
	case ",":
		fmt.Printf("stack: %v\r\n", Reverse(stack))
	case "b":
		fmt.Println("break, not implemented")
	case "c":
		fmt.Println("continue, not implemented")
	case "d":
		v, ok := stack.Pop()
		if !ok {
			fmt.Println("Error on stack, can't get value.")
			return
		}
		fmt.Println("delay: ", v, "ms")
	case "i":
		fmt.Println("pin input")
		stack.Push(1234)
	case "j":
		fmt.Println("pulse in")
		stack.Push(1234)
	case "k":
		stack.Push(loopValue)
	case "o":
		pin, ok := stack.Pop()
		if !ok {
			fmt.Println("Error on stack, can't get value.")
			return
		}
		v, ok := stack.Pop()
		if !ok {
			fmt.Println("Error on stack, can't get value.")
			return
		}
		fmt.Println("pin ", pin, ": ", v)
	case "p":
		v, ok := stack.Pop()
		if !ok {
			fmt.Println("Error on stack, can't get value.")
			return
		}
		fmt.Println("value:", v)
	case "r":
		x, ok := stack.Pop()
		if !ok {
			fmt.Println("Error on stack, can't get value.")
			return
		}
		if (x < 0) || (x > len(store)) {
			fmt.Println("invalid store address0")
			return
		}
		stack.Push(store[x])
		fmt.Println("value retrived")
	case "s":
		x, ok := stack.Pop()
		if !ok {
			fmt.Println("Error on stack, can't get value.")
			return
		}
		if (x < 0) || (x > len(store)) {
			fmt.Println("invalid store address0")
			return
		}
		v, ok := stack.Pop()
		if !ok {
			fmt.Println("Error on stack, can't get value.")
			return
		}
		store[x] = v
		fmt.Println("value stored")
	case "t":
		v, ok := stack.Pop()
		if !ok {
			fmt.Println("Error on stack, can't get value.")
			return
		}
		if v > 0 {
			fmt.Println("tone ", v, "Hz")
		} else {
			fmt.Println("tone off")
		}
	case "q":
		fmt.Println("subroutines: ")
		for k, v := range definitions {
			fmt.Printf("%s: %s\r\n", k, v)
		}
	case "\"":
		v, ok := stack.Pop()
		if !ok {
			fmt.Println("Error on stack, can't get value.")
			return
		}
		stack.Push(v)
		stack.Push(v)
	case "'":
		v, ok := stack.Pop()
		if !ok {
			fmt.Println("Error on stack, can't get value.")
			return
		}
		fmt.Println("value dropped ", v)
	case "&", "|", "^", "+", "-", "*", "/", "%":
		math(nme)
	case "~":
		v1, ok := stack.Pop()
		if !ok {
			fmt.Println("Error on stack, can't get value.")
			return
		}
		stack.Push(v1)
	case "#":
		v, ok := stack.Pop()
		if !ok {
			fmt.Println("Error on stack, can't get value.")
			return
		}
		fmt.Println("loop from 1 to ", v)
	case "?", "=", ">", "<":
		fmt.Println("skip if not ", nme)
	}
}

func math(mne string) bool {
	v1, ok := stack.Pop()
	if !ok {
		fmt.Println("Error on stack, can't get value.")
		return false
	}
	v2, ok := stack.Pop()
	if !ok {
		fmt.Println("Error on stack, can't get value.")
		return false
	}
	switch mne {
	case "&":
		stack.Push(v2 & v1)
	case "|":
		stack.Push(v2 | v1)
	case "^":
		stack.Push(v2 ^ v1)
	case "+":
		stack.Push(v2 + v1)
	case "-":
		stack.Push(v2 - v1)
	case "*":
		stack.Push(v2 * v1)
	case "/":
		stack.Push(v2 / v1)
	case "%":
		stack.Push(v2 % v1)
	}
	return true
}

func showHelp() {
	fmt.Println("Help")
	fmt.Println("[#]: push # to stack")
	fmt.Println("b: break actual block")
	fmt.Println("c: continue with next interation in loop")
	fmt.Println("d: delay in ms")
	fmt.Println("h: print help")
	fmt.Println("i: input from pin")
	fmt.Println("j: read pulse length")
	fmt.Println("k: push loop index")
	fmt.Println("o: output to pin")
	fmt.Println("r: retrive value from address")
	fmt.Println("s: store value to address")
	fmt.Println("t: tone, 0=off")
	fmt.Println("\": dupe stack value")
	fmt.Println("': drop stack value")
	fmt.Println("&: AND")
	fmt.Println("|: OR")
	fmt.Println("^: XOR")
	fmt.Println("~: NOT")
	fmt.Println("+: ADD")
	fmt.Println("-: DEL")
	fmt.Println("*: MUL")
	fmt.Println("/: DIV")
	fmt.Println("%: MOD")
	fmt.Println("p: print value")
	fmt.Println("q: print all user functions")
	fmt.Println(".: print stack size")
	fmt.Println(",: print stack")
	fmt.Println("_: print text")
	fmt.Println("=: skip if not equal")
	fmt.Println("?: skip if not null")
	fmt.Println(">: skip if not greater than")
	fmt.Println("?: skip if not lesser than")
}

func Reverse[T any](original []T) (reversed []T) {
	reversed = make([]T, len(original))
	copy(reversed, original)

	for i := len(reversed)/2 - 1; i >= 0; i-- {
		tmp := len(reversed) - 1 - i
		reversed[i], reversed[tmp] = reversed[tmp], reversed[i]
	}

	return
}

// IsEmpty: check if stack is empty
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// Push a new value onto the stack
func (s *Stack) Push(v int) {
	*s = append(*s, v) // Simply append the new value to the end of the stack
}

// Remove and return top element of stack. Return false if stack is empty.
func (s *Stack) Pop() (int, bool) {
	if s.IsEmpty() {
		return 0, false
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		*s = (*s)[:index]      // Remove it from the stack by slicing it off.
		return element, true
	}
}
