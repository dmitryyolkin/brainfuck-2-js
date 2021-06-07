package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

const BrainfuckPlus = 43     // +
const BrainfuckMinus = 45    // -
const BrainfuckDot = 46      // .
const BrainfuckComma = 44    // ,
const BrainfuckLeft = 60     // <
const BrainfuckRight = 62    // >
const BrainfuckOBracket = 91 // [
const BrainfuckCBracket = 93 // ]

var AllBrainfuckCommands = map[uint8]bool{
	BrainfuckPlus:     true,
	BrainfuckMinus:    true,
	BrainfuckDot:      true,
	BrainfuckComma:    true,
	BrainfuckLeft:     true,
	BrainfuckRight:    true,
	BrainfuckOBracket: true,
	BrainfuckCBracket: true,
}

const JsTemplateFilePath = "./brainfuck.translator.template"
const OutputDir = "./output/"

func readFile(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// Efficient string concatenation taken from
	// https://yourbasic.org/golang/build-append-concatenate-strings-efficiently/
	scanner := bufio.NewScanner(file)
	var strBuilder strings.Builder
	strBuilder.Grow(32)
	for scanner.Scan() { // internally, it advances token based on seperator
		_, err = fmt.Fprintf(&strBuilder, scanner.Text()+"\n") // token in unicode-char
		if err != nil {
			log.Fatal(err)
		}
	}
	return strBuilder.String()
}

func writeFile(fileName string, content string) {
	// create output dif if it's not exist
	if _, err := os.Stat(OutputDir); os.IsNotExist(err) {
		err := os.Mkdir(OutputDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// create output file
	outputFileName := fileName[strings.LastIndex(fileName, "/"):] + ".js"
	file, err := os.Create(OutputDir + outputFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// save file to output dir
	_, err2 := file.WriteString(content)
	if err2 != nil {
		log.Fatal(err2)
	}

}

func addJavaScriptCode(command uint8, commandCount int) (string, error) {
	switch command {
	case BrainfuckPlus:
		return fmt.Sprint("plus(", commandCount, ");\n"), nil
	case BrainfuckMinus:
		return fmt.Sprint("minus(", commandCount, ");\n"), nil
	case BrainfuckLeft:
		return fmt.Sprint("if(!ptr) { throw Error('Data pointer can not be negative'); } ptr-= ", commandCount, ";\n"), nil
	case BrainfuckRight:
		return fmt.Sprint("ptr+=", commandCount, ";\n", "arrayExtend(ptr);\n"), nil
	case BrainfuckDot:
		return "printCurrentData();\n", nil
	case BrainfuckComma:
		return "data[ptr]=readSymbolFromStdIn();\n", nil
	case BrainfuckOBracket:
		return strings.Repeat("while( data[ptr] ){ \n", commandCount), nil
	case BrainfuckCBracket:
		return strings.Repeat("}\n", commandCount), nil
	default:
		return "", errors.New(fmt.Sprintf("command '%d' is ignored", command))
	}
}

func translate(fileName string) {
	fmt.Printf("[Start]: %s\n", fileName)

	javaScriptSrc := "// ------ Java script code for: " + fileName + "\n"
	javaScriptSrc += readFile(JsTemplateFilePath)

	// translate code
	brainfuckSrc := strings.TrimSpace(readFile(fileName))
	// brainfuck statement is not empty
	var prevSymbol uint8
	prevSymbolCount := 0
	for i := 0; i < len(brainfuckSrc); i++ {
		currSymbol := brainfuckSrc[i]
		if !AllBrainfuckCommands[currSymbol] {
			// continue insufficient symbol
			fmt.Printf("current symbol %d is not valid Brainfuck command\n", currSymbol)
			continue
		}

		if prevSymbol == 0 || prevSymbol == currSymbol {
			// Count repetitions number
			if prevSymbol == 0 {
				prevSymbol = currSymbol
			}
			prevSymbolCount++
			continue
		}

		code, err := addJavaScriptCode(prevSymbol, prevSymbolCount)
		if err != nil {
			log.Fatal(err)
		}
		javaScriptSrc += code

		// increase counters
		prevSymbol = currSymbol
		prevSymbolCount = 1
	}

	code, err := addJavaScriptCode(prevSymbol, prevSymbolCount)
	if err != nil {
		log.Fatal(err)
	}
	javaScriptSrc += code

	writeFile(fileName, javaScriptSrc)
	fmt.Printf("[End]: %s\n", fileName)
}

func main() {
	translate("./examples/Hello1.b")
	//translate("./examples/Hello2.b")
	//translate("./examples/Beer.b")
	//translate("./examples/Mandelbrot.b")
}
