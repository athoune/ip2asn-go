package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	_tree "github.com/athoune/ip2asn-go/tree"
	"github.com/athoune/ip2asn-go/tsv"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	r, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	tree, err := _tree.New(256)
	if err != nil {
		panic(err)
	}
	src := tsv.New(r)
	cpt := 0
	for src.Next() {
		line, err := src.Values()
		if err != nil {
			panic(err)
		}
		n := line.Network()
		//fmt.Println(n, line.ASNumber, line.CountryCode, line.ASDescription)
		if line.ASNumber != 0 {
			if n.IP.To4() != nil {
				tree.Append(&n, line)
				cpt++
			}
		}
	}

	fmt.Println("Indexation done :", cpt)

	if len(os.Args) == 3 {
		f, err := os.Open(os.Args[2])
		if err != nil {
			panic(err)
		}
		lines := bufio.NewScanner(f)
		chrono := time.Now()
		cpt := 0
		for lines.Scan() {
			tree.Get(net.ParseIP(lines.Text()))
			cpt++
		}
		dt := time.Now().Sub(chrono)
		fmt.Println(cpt, "in", dt, "=>", int64(dt)/int64(cpt)/1000, "µs")

		//tree.Dump(os.Stdout)
	} else {
		fmt.Println("Listening 0.0.0.0:1234")
		listen, err := net.Listen("tcp", "0.0.0.0:1234")
		if err != nil {
			panic(err)
		}
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			go func(conn net.Conn) {
				scan := bufio.NewScanner(conn)
				defer conn.Close()
				for scan.Scan() {
					line := scan.Text()
					line = strings.TrimSpace(line)
					log.Println(line)
					if line == "" {
						continue
					}
					chrono := time.Now()
					ok, data := tree.Get(net.ParseIP(line))
					log.Printf("%v", time.Now().Sub(chrono))
					if ok {
						fmt.Fprintf(conn, "%s => %s\n", line, data)
					}
				}
			}(conn)
		}
	}

}
