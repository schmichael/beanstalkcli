package main

import (
    "github.com/kr/beanstalk"
)

const PRI = 0

func DialBeanstalk() *beanstalk.Conn {
    conn, err := beanstalk.Dial("tcp", *beanstalkdAddress)
    if err != nil {
        fatal(1, "Error connecting to beanstalkd:\n%v\n", err)
    }
    return conn
}
