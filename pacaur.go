package main

import (
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

func pacaurSegment(segment *segment) {
	conn, err := net.Dial("unix", "/tmp/pacaurd.sock")
	check(err)
	defer func() { check(conn.Close()) }()

	packages, err := ioutil.ReadAll(conn)
	check(err)
	num, err := strconv.Atoi(string(packages))
	check(err)

	if num > 5 {
		segment.value = icons["package"] + string(packages)
	} else {
		segment.value = strings.Repeat(icons["package"], num)
	}
}