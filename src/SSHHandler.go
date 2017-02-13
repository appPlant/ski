package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/hypersleep/easyssh.v0"
)

/**
*	Executes a command on a remote ssh server
*	@params:
*		connDet: connection details in following form: user@hostname
*		cmd: command to be executed
 */
func execCommand(user string, hostname string, command string, strucOut *StructuredOuput, opts *Opts) {

	keyPath := os.Getenv("ORBIT_KEY")
	if keyPath == "" {
		if runtime.GOOS == "windows" {
			keyPath = os.Getenv("TEMP") + "\\tempTabFormat.py"
		} else {
			keyPath = strings.TrimPrefix(path.Join(os.Getenv("ORBIT_HOME"), "config", "ssh", "orbit.key"), os.Getenv("HOME"))
		}
	}

	ssh := &easyssh.MakeConfig{
		User:   user,
		Server: hostname,
		Key:    keyPath,
		Port:   "22",
	}
	var cmd string
	if opts.loadFlag {
		cmd = fmt.Sprintf(`sh -lc "echo -----APPPLANT-ORBIT----- && %s "`, command)
	} else {
		cmd = command
	}
	// Call Run method with command you want to run on remote server.
	out, err := ssh.Run(cmd)
	// Handle errors

	if err != nil {
		message := fmt.Sprintf("called from execCommand.\nKeypath: %s\nCommand: %s", keyPath, cmd)
		log.Fatalf("%s\nAdditional info: %s\n", err, message)
	} else {
		cleanedOut := out
		if opts.loadFlag {
			splitOut := strings.Split(out, "-----APPPLANT-ORBIT-----\n")
			cleanedOut = splitOut[len(splitOut)-1]
		}
		maxLength := getMaxLength(out)
		strucOut.output = cleanedOut
		strucOut.maxOutLength = maxLength
	}
	log.Debugln("### execCommand complete ###")
	log.Debugln(fmt.Sprintf("user: %s", user))
	log.Debugln(fmt.Sprintf("hostname: %s", hostname))
	log.Debugln(fmt.Sprintf("orbit key: %s", os.Getenv("ORBIT_KEY")))
	log.Debugln(fmt.Sprintf("command: %s", command))
	log.Debugln(fmt.Sprintf("strucOut: %v", strucOut))
	log.Debugln(fmt.Sprintf("planet: %s\n maxLineLength: %d", strucOut.planet, strucOut.maxOutLength))
	// debugPrintStructuredOutput(strucOut)
	// debugPrintString("### execCommand complete ###")
}

/**
*	Uploads a file to the remote server
 */
func uploadFile(user string, hostname string, opts *Opts) {
	keyPath := os.Getenv("ORBIT_KEY")
	if keyPath == "" {
		if runtime.GOOS == "windows" {
			keyPath = os.Getenv("TEMP") + "\\tempTabFormat.py"
		} else {
			keyPath = strings.TrimPrefix(path.Join(os.Getenv("ORBIT_HOME"), "config", "ssh", "orbit.key"), os.Getenv("HOME"))
		}
	}
	ssh := &easyssh.MakeConfig{
		User:   user,
		Server: hostname,
		Key:    keyPath,
		Port:   "22",
	}

	// Call Scp method with file you want to upload to remote server.
	err := ssh.Scp(path.Join(os.Getenv("ORBIT_HOME"), scriptDirectory, opts.scriptName))

	// Handle errors
	if err != nil {
		message := fmt.Sprintf("called from uploadFile. Keypath: %s", keyPath)
		log.Fatalf("%s\nAdditional info: %s\n", err, message)
	}
}

/**
*	Uploads and executes a script on a given planet
*	@params:
*		connDet: 	Connection details to planet
*		scriptPath: Path to script
 */
func execScript(user string, hostname string, strucOut *StructuredOuput, opts *Opts) {
	uploadFile(user, hostname, opts)
	placeholder := StructuredOuput{}
	scriptName := opts.scriptName
	executionCommand := fmt.Sprintf("sh %s", scriptName)
	delCommand := fmt.Sprintf("rm %s", scriptName)
	execCommand(user, hostname, executionCommand, strucOut, opts)
	execCommand(user, hostname, delCommand, &placeholder, opts)
}
