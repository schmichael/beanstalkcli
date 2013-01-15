package main

import (
    "flag"
    "fmt"
    "os"
)

type Command struct {
    // Runs the command
    Run func(cmd *Command)

    Name string

    Desc string

    Flag flag.FlagSet
}


var commands = []*Command{
    cmdPut,
    cmdGet,
}
var beanstalkdAddress = flag.String("h", "localhost:11300", "hostname:port of beanstalkd")


// Based on go command's main.go
func main() {
    flag.Parse()

    args := flag.Args()

    if len(args) == 0 {
        printUsage()
        os.Exit(1)
    }

    if args[0] == "help" {
        printHelp()
        return
    }

    found := false
    for _, c := range commands {
        if args[0] == c.Name {
            found = true
            c.Flag.Init(c.Name, flag.ExitOnError)
            c.Flag.Parse(args[1:])
            c.Run(c)
        }
    }

    if !found {
        fatal(1, "%s is not a valid command", args[0])
    }
}

func fatal(status int, msg string, args ...interface{}) {
    writeStderr(msg, args...)
    os.Exit(status)
}

func writeStderr(msg string, args ...interface{}) {
    os.Stderr.WriteString(fmt.Sprintf(msg, args...))
    os.Stderr.Write([]byte("\n"))
    os.Stderr.Sync()
}

func printUsage() {
    writeStderr("%s COMMAND [OPTIONS]", os.Args[0])
}

func printHelp() {
    printUsage()
    writeStderr("")

    for _, c := range commands {
        writeStderr("    %-10s %s", c.Name, c.Desc)
    }
}
