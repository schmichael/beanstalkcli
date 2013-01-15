package main

import (
    "github.com/kr/beanstalk"
    "fmt"
    "strings"
)

var cmdGet = &Command{
    Name: "get",
    Desc: "get a job (reserve)",
}
var getTubes = cmdGet.Flag.String("t", "default", "comma separated list of tubes")
var getNum = cmdGet.Flag.Uint64("n", 1, "number to get, 0 gets all")
var getAction = cmdGet.Flag.String("x", "r", "action to take: [r]elease, [d]elete, [b]ury, [n]othing")
var Actions = map[string]func(*beanstalk.Conn, uint64, []byte) {}

func init() {
    cmdGet.Run = runGet
    Actions["r"] = getRelease
    Actions["d"] = getDelete
    Actions["b"] = getBury
    Actions["n"] = getNoop
}

func getRelease(c *beanstalk.Conn, id uint64, body []byte) {
    c.Release(id, PRI, 0)
}

func getDelete(c *beanstalk.Conn, id uint64, body []byte) {
    c.Delete(id)
}

func getBury(c *beanstalk.Conn, id uint64, body []byte) {
    c.Bury(id, 0)
}

func getNoop(c *beanstalk.Conn, id uint64, body []byte) {}

func runGet(cmd *Command) {
    conn := DialBeanstalk()
    ts := beanstalk.NewTubeSet(conn, strings.Split(*getTubes, ",")...)
    n := *getNum
    var ok bool
    var action func(*beanstalk.Conn, uint64, []byte)
    if action, ok = Actions[*getAction]; !ok {
        fatal(2, "'%s' isn't a valid action", *getAction)
    }
    if *getAction == "r" && n == 0 {
        // Protect users from themselves
        fatal(2, "Using -n 0 and -x r together causes a tight loop and is disallowed")
    }
    for i := uint64(0); n == 0 || i < n; i++ {
        id, body, err := ts.Reserve(0)
        if err != nil {
            if cerr, ok := err.(beanstalk.ConnError); ok && cerr.Err == beanstalk.ErrTimeout {
                // Only write message if no jobs at all, but exit w/ 0
                if i == 0 {
                    writeStderr("No jobs")
                }
                return
            }
            fatal(2, "Error getting job:\n%v", err)
        }
        fmt.Printf("%s\n", body)
        action(conn, id, body)
    }
}
