package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/willie68/basl/internal/logging"

	flag "github.com/spf13/pflag"
)

type Stack []int
type ReaderEntry struct {
	Reader     *bufio.Reader
	Console    bool
	Subroutine bool
}
type ReaderStack []ReaderEntry

const COMMANDS string = " 0123456789abcdefghijklmnopqrstuvwxyz!\"/§%&={}+?*~-_#:;.,^°|><'`´\\[]@$"
const USRCMD string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

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
	reader       *bufio.Reader
	readerStack  ReaderStack
	v            int
	inNumber     bool
	inSubroutine bool
	pins         []byte
)

func init() {
	// variables for parameter override
	fs = flag.NewFlagSet("main", flag.ContinueOnError)
	fs.StringVarP(&baslFile, "file", "f", "", "source file to compile")
	fs.SortFlags = false

	stack = make([]int, 0)
	store = make([]int, 1024)
	readerStack = make([]ReaderEntry, 0)
	definitions = make(map[string]string)
	v = 0
	inNumber = false
	console = true
	inSubroutine = false
	pins = []byte("iiiiooooippixoaa")
}

func main() {
	log.Info("starting basl for pc")
	log.Logger.SetLevel(log.LvDebug)
	fs.Parse(os.Args[1:])

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
		c, err := readNme()
		if err != nil {
			fmt.Println("error: ", err)
		}
		if c > 0 {
			execute(c)
		}

		if console && (reader.Buffered() == 0) {
			fmt.Print(":")
		}
	}
}

func readNme() (byte, error) {
	rune := make([]byte, 1)
	_, err := reader.Read(rune)
	if (err == io.EOF) && ((baslFile != "") || inSubroutine) {
		if inSubroutine {
			re, ok := readerStack.Pop()
			if !ok {
				return 0, errors.New("no reader in stack")
			}
			inSubroutine = re.Subroutine
			reader = re.Reader
			console = re.Console
			return rune[0], nil
		}
		reader = bufio.NewReader(os.Stdin)
		console = true
		err = nil
	}
	if err != nil {
		return 0, err
	}
	return rune[0], nil
}

func execute(chr byte) {
	schr := string(chr)

	if chr == ':' {
		log.Debug("start defining user cmd")
		inDefinition = true
		return
	}

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

	if strings.Contains(COMMANDS, schr) {
		processNme(schr)
		return
	}

	if strings.Contains(USRCMD, schr) {
		log.Debug("user command: " + schr)
		def, ok := definitions[schr]
		if !ok {
			log.Error("user command not found! " + schr)
			return
		}
		re := ReaderEntry{
			Reader:     reader,
			Console:    console,
			Subroutine: inSubroutine,
		}
		readerStack.Push(re)
		reader = bufio.NewReader(strings.NewReader(def))
		inSubroutine = true
		console = false
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
		n, ok := getNumber()
		if !ok {
			fmt.Println("Error on input, can't get value.")
			return
		}
		stack.Push(n)
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
	case "&", "|", "^", "+", "-", "*", "/", "%":
		math(nme)
	case "~":
		v1, ok := stack.Pop()
		if !ok {
			fmt.Println("Error on stack, can't get value.")
			return
		}
		stack.Push(^v1)
	case "_":
		for {
			c, err := readNme()
			if err != nil {
				fmt.Println("Error: ", err)
				break
			}
			if c == '_' {
				fmt.Println()
				break
			}
			fmt.Print(string(c))
		}
	case "#":
		v, ok := stack.Pop()
		if !ok {
			fmt.Println("Error on stack, can't get value.")
			return
		}
		fmt.Println("loop from 0 to ", v-1)
		loopValue = 1
		c, err := readNme()
		if err != nil {
			fmt.Println("error: ", err)
			return
		}
		if c == '{' {
			// process a block
			text, err := readBlock()
			if err != nil {
				fmt.Println("error: ", err)
				return
			}
			for loopValue = 0; loopValue < v; loopValue++ {
				for _, c := range text {
					processNme(string(c))
				}
			}
		} else {
			// process a single command
			for loopValue = 0; loopValue < v; loopValue++ {
				processNme(string(c))
			}
		}
		loopValue = 0
	case "?":
		log.Debug("skip if > 0")
		v, ok := stack.Pop()
		if !ok {
			fmt.Println("Error on stack, can't get value.")
			return
		}
		nextBlockOrNot(v > 0)
	case "=", ">", "<":
		log.Debug("skip if v1 " + nme + " v2")
		v2, ok := stack.Pop()
		if !ok {
			fmt.Println("Error on stack, can't get value.")
			return
		}
		v1, ok := stack.Pop()
		if !ok {
			fmt.Println("Error on stack, can't get value.")
			return
		}
		doWork := false
		switch nme {
		case "=":
			doWork = v1 == v2
		case ">":
			doWork = v1 > v2
		case "<":
			doWork = v1 < v2
		}
		nextBlockOrNot(doWork)
	case "@":
		log.Debug("config command")
		pins = make([]byte, 0)
		for {
			c, err := readNme()
			if err != nil {
				fmt.Println("Error: ", err)
				break
			}
			switch c {
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
			if (c == '\r') || (c == '\n') {
				fmt.Println()
				break
			}
			fmt.Print(string(c))
		}
	case "$":
		log.Debug("output pin configuration")
		for _, p := range pins {
			fmt.Print(string(p))
		}
		fmt.Println()
	case "{", "}":
		// nothing to do here
	default:
		log.Errorf("unknown command: %v", nme)
	}
}

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
	c, err := readNme()
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	if c == '{' {
		// process a block
		text, err := readBlock()
		if err != nil {
			fmt.Println("error: ", err)
			return
		}
		if doWork {

			re := ReaderEntry{
				Reader:     reader,
				Console:    console,
				Subroutine: inSubroutine,
			}
			readerStack.Push(re)
			reader = bufio.NewReader(strings.NewReader(text))
			inSubroutine = true
			console = false
		}
	} else {
		if doWork {
			execute(c)
		}
	}
}

func readBlock() (string, error) {
	text := ""
	for {
		c, err := readNme()
		if err != nil {
			return "", err
		}
		if c == '}' {
			break
		}
		text = text + string(c)
	}
	return text, nil
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

func (s *Stack) Clear() {
	*s = make([]int, 0)
}

// IsEmpty: check if stack is empty
func (s *ReaderStack) IsEmpty() bool {
	return len(*s) == 0
}

// Push a new value onto the stack
func (s *ReaderStack) Push(r ReaderEntry) {
	*s = append(*s, r) // Simply append the new value to the end of the stack
}

// Remove and return top element of stack. Return false if stack is empty.
func (s *ReaderStack) Pop() (ReaderEntry, bool) {
	if s.IsEmpty() {
		return ReaderEntry{}, false
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		*s = (*s)[:index]      // Remove it from the stack by slicing it off.
		return element, true
	}
}

func (s *ReaderStack) Clear() {
	*s = make([]ReaderEntry, 0)
}
