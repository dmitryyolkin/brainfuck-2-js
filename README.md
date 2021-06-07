# Environment
The project was implemented and tested with:
* nodejs v14.8.0
* go v1.16.5 darwin/amd64

# Task
Implement Brainfuck translator to JavaScript on Golang.
Brainfuck is a programming language. Details:
* [Brainfuck description](https://en.wikipedia.org/wiki/Brainfuck)
* [ASCII code table](https://snipp.ru/handbk/table-ascii) - Brainfuck input/output ASCII symbols

# Implementation details
* Some examples exist in `{project_dir}/examples` folder
* There are two implementations:
  * `Brainfuck.js` is JavaScript implementation. It's NOT required by the test task but it was easier to implement on JavaScript because
    * JavaScript is more familiar to me and it helps me saving the time on debugging the main algorithm
    * JavaScript allows to run JS code from runtime with `Function` structure
    * As long as it ws done - I've decided to keep it in final version
  * `Brainfuck.go` is Golang implementation
* Both implementations put output files into `{project_dir}/output` folder in format `input.b.js` where `input.b` is initial file with Brainfuck instruction.

# Run
You can run examples as follows:
* `Brainfuck.js`:
  * go to `Brainfuck.js`, uncomment in `main()` function required examples (one or more) 
  * run in terminal `node Brainfuck.js` 
  * results are be printed in console log and saved in `{project_dir}/output` folder
    * file from output dir can be run the same way with `node`
* `Brainfuck.go`:
  * go to `Brainfuck.go`, uncomment in `main()` function required examples (one or more)
  * run in terminal `go run Brainfuck.go`
  * results are saved in `{project_dir}/output` folder