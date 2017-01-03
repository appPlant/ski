package main

import (
	"os"
)

/**
################################################################################
								Main-Section
################################################################################
*/

/**
*	StructuredOuput:
*	A thingy thing thing
 */
type StructuredOuput struct {
	planet       string
	output       string
	maxOutLength int
}

/**
*	Main function
 */
func main() {
	args := os.Args
	opts := Opts{}
	exec := Executor{}

	opts.procArgs(args)
	exec.execMain(&opts)

}
