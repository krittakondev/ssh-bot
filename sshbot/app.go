package sshbot

import (
	"bytes"
	"fmt"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

func Runcommand(host string, password string, command string) {
	// รับพารามิเตอร์จากคำสั่งรัน

	port := "22"
	username := "root"

	// สร้างการเชื่อมต่อโดยใช้ ssh.Dial
	conn, err := ssh.Dial("tcp", host+":"+port, &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}

	session, err := conn.NewSession()
	if err != nil {
		log.Fatal("Failed to Session: ", err)
	}
	var buff bytes.Buffer
	session.Stdout = &buff
	if err := session.Run(command); err != nil {
		log.Fatal(err)
	}
	fmt.Println(buff.String())
}
