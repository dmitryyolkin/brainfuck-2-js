'use strict'

const fs = require('fs')

// Links:
//  * Brainfuck from wiki -> https://ru.wikipedia.org/wiki/Brainfuck#%D0%AF%D0%B7%D1%8B%D0%BA%D0%B8_%D0%BD%D0%B0_%D0%BE%D1%81%D0%BD%D0%BE%D0%B2%D0%B5_Brainfuck
//  * ASCII codes table -> https://snipp.ru/handbk/table-ascii

const BRAINFUCK_ASCII_CODE = {
    Plus: 43, // +
    Minus: 45, // -
    Dot: 46, // .
    Comma: 44, // ,
    Left: 60, // <
    Right: 62, // >
    OBracket: 91, // [
    CBracket: 93, // ]
}

const ALL_BRAINFUCK_COMMANDS = new Set([
    BRAINFUCK_ASCII_CODE.Plus,
    BRAINFUCK_ASCII_CODE.Minus,
    BRAINFUCK_ASCII_CODE.Dot,
    BRAINFUCK_ASCII_CODE.Comma,
    BRAINFUCK_ASCII_CODE.Left,
    BRAINFUCK_ASCII_CODE.Right,
    BRAINFUCK_ASCII_CODE.OBracket,
    BRAINFUCK_ASCII_CODE.CBracket
])

const OUTPUT_DIR = './output/';

function addJavaScriptCode(command, commandCount) {
    switch (command) {
        case BRAINFUCK_ASCII_CODE.Plus:
            return 'plus(' + commandCount + ');\n'
        case BRAINFUCK_ASCII_CODE.Minus:
            return 'minus(' + commandCount + ');\n'
        case BRAINFUCK_ASCII_CODE.Left:
            return 'ptr-= ' + commandCount + ';\n';
        case BRAINFUCK_ASCII_CODE.Right:
            return 'ptr+=' + commandCount + ';\n' + 'arrayExtend(ptr);\n';
        case BRAINFUCK_ASCII_CODE.Dot:
            return 'printCurrentData();\n';
        case BRAINFUCK_ASCII_CODE.Comma:
            return 'data[ptr]=readSymbolFromStdIn();\n'
        case BRAINFUCK_ASCII_CODE.OBracket:
            return 'while( data[ptr] ){ \n'.repeat(commandCount);
        case BRAINFUCK_ASCII_CODE.CBracket:
            return '}\n'.repeat(commandCount);
        default:
            throw new Error(`symbol "${currSymbol}" is ignored`)
    }
}

function main(fileName) {
    console.log(`[Start]: ${fileName}`)

    let javaScriptSrc = "// ------ Java script code for: " + fileName + "\n"
    javaScriptSrc += fs.readFileSync('./brainfuck.translator.template', 'utf8');

    // translate code
    let prevSymbol = undefined
    let prevSymbolCount = 0
    const brainfuckSrc = fs.readFileSync(fileName)
    for (let i = 0; i < brainfuckSrc.length; i++) {
        const currSymbol = brainfuckSrc.readUInt8(i)
        if (!ALL_BRAINFUCK_COMMANDS.has(currSymbol)) {
            // continue insufficient symbol
            console.log(`current symbol "${currSymbol}" is not valid Brainfuck command`)
            continue
        }

        if (prevSymbol === undefined || prevSymbol === currSymbol) {
            // Count repetitions number
            // increase symbol count and move to next symbol
            if (prevSymbol === undefined) {
                prevSymbol = currSymbol
            }
            prevSymbolCount++
            continue
        }

        javaScriptSrc += addJavaScriptCode(prevSymbol, prevSymbolCount)

        // increase counters
        prevSymbol = currSymbol
        prevSymbolCount = 1
    }

    if (prevSymbol !== undefined) {
        javaScriptSrc += addJavaScriptCode(prevSymbol, prevSymbolCount)
    }


    let outputFileName = fileName
        .substring(fileName.lastIndexOf('/') + 1, fileName.length)
        .concat('.js')

    console.log(`java script src code:\n${javaScriptSrc}`)
    if (!fs.existsSync(OUTPUT_DIR)){
        fs.mkdirSync(OUTPUT_DIR);
    }
    fs.writeFileSync(OUTPUT_DIR + outputFileName, javaScriptSrc)

    const f = new Function('require', javaScriptSrc)
    f(require) // require fs to be able to import modules

    console.log(`[End]: ${fileName}`)
}

main('./examples/Hello1.b')
// main('./examples/Hello2.b')
// main('./examples/Beer.b')
// main('./examples/Mandelbrot.b')
