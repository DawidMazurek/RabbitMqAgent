package main

import (
    "log"
    "net"
    "fmt"
)

func main() {
    l, err := net.Listen("unix", "/tmp/rabbitmqagent.sock")
    if err != nil {
        log.Fatal("listen error:", err)
    }

    for {
        fd, err := l.Accept()
        if err != nil {
            log.Fatal("accept error:", err)
        }

        go deliverMessage(fd)
    }
}

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
        panic(fmt.Sprintf("%s: %s", msg, err))
    }
}
