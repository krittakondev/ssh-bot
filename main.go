package main

import (
	"b13sshbot/sshbot"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func get_list() ([]string, []string) {
	hosts, err := os.ReadFile("./host.txt")
	check(err)
	pass, err := os.ReadFile("./pass.txt")
	check(err)
	hosts_split := strings.Split(string(hosts), "\n")
	pass_split := strings.Split(string(pass), "\n")
	return hosts_split, pass_split
}

func main() {
	hosts, pass := get_list()
	if len(hosts) != len(pass) {
		fmt.Println("จำนวน host.txt กับ pass.txt ไม่เท่ากับ")
		return
	}
	var command string
	var confirm string
	fmt.Printf("ใส่คำสั่งที่ต้องการรัน: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		command = scanner.Text()
	}
	fmt.Printf("ยืนยันจะรันทั้งหมด %d server (y/N): ", len(hosts))
	fmt.Scanln(&confirm)
	if strings.ToLower(confirm) == "y" || strings.ToLower(confirm) == "yes" {
		fmt.Println("เริ่มรันคำสั่ง... ")
		for i := 0; i < len(hosts); i++ {
			fmt.Printf("(%d) ip: %s \n", i+1, hosts[i])
			fmt.Println(sshbot.Runcommand(hosts[i], pass[i], command))
		}
		return
	} else {
		fmt.Println("ทำการยกเลิก!! ")
		return
	}

}
