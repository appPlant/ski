package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
    "golang.org/x/crypto/ssh"
    //"bytes"
    "io"
    "io/ioutil"
    //"time"
)

func printVersion(){
	fmt.Println("0.0.1")
}

func throwErr(out []byte, err error){
	fmt.Print(fmt.Sprint(err) + ": " + string(out) + "path is " + os.Getenv("PATH"))
	os.Stderr.WriteString(fmt.Sprint(err) + ": " + string(out))
	os.Exit(1)
}

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}




func getHost(connDet string) string{
	toReturn := strings.Split(connDet,"@")
	return toReturn[1]
}

func getUser(connDet string) string{
	toReturn := strings.Split(connDet,"@")
	return toReturn[0]
}

func execCommand(connDet string, cmd string){
    //timeout := time.After(5 * time.Second) // in 5 seconds the message will come to timeout channel

	hostname := getHost(connDet)
    user := getUser(connDet)
    sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
		PublicKeyFile("/root/.ssh/id_rsa"),
		},
	}

    connection, err := ssh.Dial("tcp", hostname + ":22", sshConfig)

	if err != nil {
		fmt.Println("error on dialup")
		os.Exit(1)
	}

	session, err := connection.NewSession()
	if err != nil {
		fmt.Println("Failed to create session: %s", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		fmt.Println("Unable to setup stdin for session: %v", err)
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		fmt.Println("Unable to setup stdout for session: %v", err)
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		fmt.Println("Unable to setup stderr for session: %v", err)
	}
	go io.Copy(os.Stderr, stderr)

	err = session.Run("ls -l $LC_USR_DIR")

	//fmt.Println(string(stdout))
	//fmt.Println(string(stderr))


}



func getType(id string) string{
	cmd 	 := exec.Command("ff","-t" ,id)
	out, err := cmd.CombinedOutput()
	if err != nil {
    	throwErr(out,err)
	}
	return strings.TrimSpace(string(out))
}

func getConnDet(id string) string{
	cmd 	 := exec.Command("ff",id)
	out, err := cmd.CombinedOutput()
	if err != nil {
    	throwErr(out,err)
	}
	return strings.TrimSpace(string(out))
}

func getArg(args []string, argType string) string{
	switch argType{
		case "command":
			var command string  = args[2]
			var cmdArgs []string
			if(len(args) >= 3){
				cmdArgs = args[2:(len(args))]
				for _, argument := range cmdArgs {
					command += (" " + argument)
				}
			}
			return command
		case "id":
			return args[1]
		default:
			return "Unhandled Arg"
	}

}

func main() {
	var args []string = os.Args
	if(len(args) <2){
		fmt.Println("usage: goo <Server-ID> <Command> [Arguments]")
		os.Exit(1)
	}
	if (len(args) <=2) {
		switch args[1]{
		case "-h":
			fmt.Println("usage: goo <Server-ID> <Command> [Arguments]")
			os.Exit(0)
		case "-v":
			printVersion()
			os.Exit(0)
		default:
			fmt.Println("Could not interpret arguments. For help, run goo -h")
			os.Exit(1)
		}
	}
	switch getType(getArg(args,"id")) {
		case "server":
			var connDet string = getConnDet(getArg(args,"id"))
			var command string = getArg(args,"command")
			execCommand(connDet,command)
		case "db":
			fmt.Println("This Type of Connection is not yet supported.")
			os.Exit(1)
		case "web":
			fmt.Println("This Type of Connection is not supported.")
			os.Exit(1)
		default:
			fmt.Println(getType(getArg(args,"id")))
			fmt.Println("###")
			fmt.Println(getArg(args,"id"))

	}
}
