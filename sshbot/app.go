package sshbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"net/http"

	"golang.org/x/crypto/ssh"
)

var main_path string = ".b13"

var token string

var count_login int
var url_enpoint string = "https://host.b-13.co"
var url_api string = url_enpoint + "/api"

func checkDir() {
	if _, err := os.Stat(main_path); err != nil {
		os.Mkdir(main_path, 0755)
	}
}

func Login() {
	var email string
	var password string
	fmt.Println("เข้าสู่ระบบ หากยังไม่มีสมาชิกสมัครได้ที่ " + url_enpoint + "/register")
	count_login++
	url := url_api + "/login"
	fmt.Printf("email: ")
	fmt.Scanf("%s \n", &email)
	fmt.Printf("pass: ")
	fmt.Scanf("%s \n", &password)

	var jsonStr = []byte(`{"email":"` + email + `", "password": "` + password + `"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	type resptype struct {
		success bool
		message string
		status  string
		data    map[string]interface{}
	}
	var resp_body map[string]interface{}

	content, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(string(content)), &resp_body)
	// fmt.Printf("%v", resp_body)
	// fmt.Printf("%T", resp_body)

	if resp_body["success"] == true {
		// return_token := resp_data["data"]["token"]
		md := (resp_body["data"]).(map[string]interface{})

		// fmt.Printf("%v\n", md)
		// fmt.Printf("%Tn", md)
		token = md["jwt_token"].(string)
		// fmt.Println(token)
		os.WriteFile(main_path+"/auth", []byte(token), 0644)
	} else {
		fmt.Printf("%v\n", resp_body["message"])
		if count_login >= 3 {
			fmt.Println("เข้าสู่ระบบผิดพลาด 3 ครั้ง!!!")
			os.Exit(0)
		} else {
			Login()

		}
	}
}

func Checkauth() {
	checkDir()
	readAuth, err := os.ReadFile(main_path + "/auth")
	if err != nil {
		// os.WriteFile(main_path+"/auth", []byte("this is token"), 0644)
		Login()
		// fmt.Println("set new token")
		// token = token
	}
	token = string(readAuth)
}

func Auth() bool {
	// url := url_api + "/bot/check_exec"
	url := url_api + "/user/me"

	req, _ := http.NewRequest("GET", url, nil)
	get_token := GetToken()
	req.Header.Set("token", get_token)
	client := &http.Client{}

	fmt.Println("checking token...")
	res, err := client.Do(req)
	// res, err := client.Get(url)
	if err != nil {
		fmt.Println("errrrrr")
		panic(err)
	}
	fmt.Println("check done")
	defer res.Body.Close()
	var resp_body map[string]interface{}

	content, _ := ioutil.ReadAll(res.Body)
	fmt.Println(content)
	json.Unmarshal([]byte(string(content)), &resp_body)
	if resp_body["success"] == false {
		fmt.Println(resp_body["message"])
		return false
	}
	return true
}

func GetToken() string {
	Checkauth()
	return token
}

func Runcommand(host string, password string, command string) string {
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
		fmt.Printf("%v", err)
		// log.Fatal(err)
	}
	session.Close()
	conn.Close()
	return buff.String()
}
