package main

import (
    "github.com/kr/beanstalk"
    "time"
)

var cmdPut = &Command{
    Name: "put",
    Desc: "put a job into the tube",
}
var putTube = cmdPut.Flag.String("t", "default", "tube")

func init() {
    cmdPut.Run = runPut
}

func runPut(cmd *Command) {
    conn := DialBeanstalk()
    t := &beanstalk.Tube{Conn: conn, Name: *putTube}
    id, err := t.Put([]byte(cmd.Flag.Args()[0]), 1, 0, 120 * time.Second)
    if err != nil {
        fatal(2, "Error putting job:\n%v\n", err)
    }
    writeStderr("Put ID: %d", id)
}
