package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/willie68/basl/internal/logging"

	flag "github.com/spf13/pflag"
)

type Stack []int
type Block struct {
	content string
	pos     int
}

const STACKSIZE int = 128
const STORESIZE int = 128
const COMMANDS string = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"

var (
	baslFile     string
	fs           *flag.FlagSet
	console      bool
	stack        Stack
	store        []int
	loopValue    int
	inDefinition bool
	definitions  map[string]string
	defName      string
	defText      string
	v            int
	inNumber     bool
	inSubroutine bool
	inOutput     bool
	inConfig     bool
	pins         []byte
	b2i          = map[bool]int{false: 0, true: 1}
	block        Block
)

func init() {
	// variables for parameter override
	fs = flag.NewFlagSet("main", flag.ContinueOnError)
	fs.StringVarP(&baslFile, "file", "f", "", "source file to compile")
	fs.SortFlags = false

	stack = make([]int, 0)
	store = make([]int, STORESIZE)
	definitions = make(map[string]string)
	v = 0
	inNumber = false
	console = true
	inSubroutine = false
	inOutput = false
	inConfig = false
	pins = []byte("iiiiooooippixoaa")
}

func main() {
	log.Info("starting basl for pc V0.1")
	log.Logger.SetLevel(log.LvInfo)
	fs.Parse(os.Args[1:])
	var reader *bufio.Reader

	if baslFile != "" {
		f, err := os.Open(baslFile)
		if err != nil {
			panic(fmt.Sprintf("file not found: %s", baslFile))
		}
		defer f.Close()
		console = false
		reader = bufio.NewReader(f)
	} else {
		reader = bufio.NewReader(os.Stdin)
		console = true
	}

	for {
		line, err := reader.ReadString('\n')
		if (err == io.EOF) && !console {
			reader = bufio.NewReader(os.Stdin)
			console = true
			err = nil
		}
		if err != nil {
			fmt.Println("error: ", err)
		}

		evaluate(line)

		if console && (reader.Buffered() == 0) {
			fmt.Print(":")
		}
	}
}

func evaluate(line string) {
	for _, c := range line {
		execute(byte(c))
	}
	inConfig = false
}

func execute(chr byte) {
	schr := string(chr)

	if inDefinition {
		if chr == ';' {
			log.Debug("stop defining user cmd")
			inDefinition = false
			definitions[defName] = strings.TrimSpace(defText)
			defName = ""
			defText = ""
			return
		}
		if defName == "" {
			defName = schr
			return
		}
		if defText == "" {
			defText = schr
			return
		}
		defText = defText + schr
		return
	}

	if inConfig {
		switch chr {
		case 'i', 'I':
			// digital input
			pins = append(pins, 'i')
		case 'o', 'O':
			// digital output
			pins = append(pins, 'o')
		case 'a', 'A':
			// analog input
			pins = append(pins, 'a')
		case 'p', 'P':
			// pwm output
			pins = append(pins, 'p')
		case 's', 'S':
			// servo output
			pins = append(pins, 's')
		case 'x', 'X':
			// pin not used
			pins = append(pins, 'x')
		}
		return
	}

	if schr == "_" {
		inOutput = !inOutput
		if !inOutput {
			fmt.Println()
		}
		return
	}

	if inOutput {
		fmt.Print(schr)
		return
	}

	if inBLock {

	}

	if strings.Contains(COMMANDS, schr) {
		processNme(schr)
		return
	}
}

func processNme(nme string) {
	if (nme[0] >= '0') && (nme[0] <= '9') {
		inNumber = true
		v = v*10 + (int(nme[0]) - int('0'))
		return
	} else {
		if inNumber {
			stack.Push(v)
			v = 0
			inNumber = false
		}
	}
	if nme[0] == ' ' {
		return
	}

	switch nme {
	case ":":
		inDefinition = true
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
		v, ok := stack.Pop()
		if !ok {
			fmt.Println("Error on stack, can't get value.")
			return
		}
		fmt.Println("pin input ", v)
		stack.Push(1234)
	case "j":
		v, ok := stack.Pop()
		if !ok {
			fmt.Println("Error on stack, can't get value.")
			return
		}
		fmt.Println("pulse in ", v)
		stack.Push(1234)
	case "k":
		stack.Push(loopValue)
	case "n":
		fmt.Println("not implemented yet")
		/*
			n, ok := getNumber()
			if !ok {
				fmt.Println("Error on input, can't get value.")
				return
			}
			stack.Push(n)
		*/
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
		fmt.Println(v)
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
	case "!":
		v1, ok := stack.Pop()
		if !ok {
			fmt.Println("Error on stack, can't get value.")
			return
		}
		v2, ok := stack.Pop()
		if !ok {
			fmt.Println("Error on stack, can't get value.")
			return
		}
		stack.Push(v1)
		stack.Push(v2)
		fmt.Println("values swapped")
	case "z":
		stack.Clear()
		fmt.Println("stack cleared ")
	case "&", "|", "^", "+", "-", "*", "/", "%", "=", ">", "<":
		math(nme)
	case "~":
		v1, ok := stack.Pop()
		if !ok {
			fmt.Println("Error on stack, can't get value.")
			return
		}
		stack.Push(^v1)
	case "@":
		if !inConfig {
			pins = make([]byte, 0)
			inConfig = true
		} else {
			inConfig = false
		}
	case "$":
		log.Debug("output pin configuration")
		for _, p := range pins {
			fmt.Print(string(p))
		}
		fmt.Println()
	case "{":
		inBlock = true
		
	case "}":
		inBlock = false
		// nothing to do here
	case "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z":
		log.Debug("user command: " + nme)
		def, ok := definitions[nme]
		if !ok {
			log.Error("user command not found! " + nme)
			return
		}
		log.Debug("eval: " + def)
		evaluate(def)
	default:
		log.Errorf("unknown command: %v", nme)
	}
}

/*
func getNumber() (int, bool) {
	var r *bufio.Reader
	v := 0
	rune := make([]byte, 1)
	if !console {
		r = bufio.NewReader(os.Stdin)
	} else {
		r = reader
	}
	fmt.Print(">")
	for {
		_, err := r.Read(rune)
		if err != nil {
			log.Errorf("error on input: %v", err)
		}
		if (rune[0] >= '0') && (rune[0] <= '9') {
			v = v*10 + (int(rune[0]) - int('0'))
		}
		if (rune[0] == '\r') || (rune[0] == '\n') {
			break
		}
	}
	return v, true
}

func nextBlockOrNot(doWork bool) {
	_, err := readNme()
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
}

*/

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
	case ">":
		stack.Push(b2i[v2 > v1])
	case "=":
		stack.Push(b2i[v2 == v1])
	case "<":
		stack.Push(b2i[v2 < v1])
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
	fmt.Println("n: input a number")
	fmt.Println("o: output to pin")
	fmt.Println("@: set pin configuration")
	fmt.Println("$: output pin configuration")
	fmt.Println("r: retrive value from address")
	fmt.Println("s: store value to address")
	fmt.Println("t: tone, 0=off")
	fmt.Println("\": dupe stack value")
	fmt.Println("': drop stack value")
	fmt.Println("§: swap first 2 values on stack")
	fmt.Println("°: clear stack")
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
	fmt.Println("_: print text till next _")
	fmt.Println("=: skip if not equal")
	fmt.Println("?: skip if not null")
	fmt.Println(">: skip if not greater than")
	fmt.Println("<: skip if not lesser than")
	fmt.Println("{..}}: defining a block")
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
	if len(*s) < STACKSIZE {
		*s = append(*s, v) // Simply append the new value to the end of the stack
	}
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

func (s *Stack) Clear() {
	*s = make([]int, 0)
}
