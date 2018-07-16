package main

import (
    "log"
    "net"
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

