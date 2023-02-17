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
	"golang.org/x/crypto/ssh/terminal"
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

func PrintCredit(){
	fmt.Println("Credit: https://github.com/krittakondev")
	fmt.Println("")
}

func Login() {
	var email string
	var password string
	fmt.Println("\nเข้าสู่ระบบ หากยังไม่มีสมาชิกสมัครได้ที่ " + url_enpoint + "/register")
	count_login++
	url := url_api + "/login"
	fmt.Printf("email: ")
	fmt.Scanf("%s \n", &email)
	fmt.Printf("pass(รหัสผ่านจะไม่โชว์ตอนพิมพ์): ")
	// fmt.Scanf("%s \n", &password)
	// password, e := terminal.ReadPassword(int(os.Stdin.Fd()))
	passwd, e := terminal.ReadPassword(int(os.Stdin.Fd()))
	if e != nil {
		log.Fatal(e)
	}
	password = string(passwd)
	fmt.Println()

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

	fmt.Printf("%v\n", resp_body["message"])
	if resp_body["success"] == true {
		// return_token := resp_data["data"]["token"]
		md := (resp_body["data"]).(map[string]interface{})

		// fmt.Printf("%v\n", md)
		// fmt.Printf("%Tn", md)
		// fmt.Println(token)
		mytoken := md["jwt_token"].(string)
		os.WriteFile(main_path+"/auth", []byte(mytoken), 0644)
		os.Exit(0)
	} else {
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
	}
	token = string(readAuth)
}

func SetToken(my_token string) {
	token = my_token

}

func Auth() bool {
	// url := url_api + "/bot/check_exec"
	url := url_api + "/user/me"

	req, _ := http.NewRequest("GET", url, nil)
	get_token := GetToken()
	req.Header.Set("auth-token", get_token)
	client := &http.Client{}

	fmt.Println("checking token...")
	res, err := client.Do(req)
	// res, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	var resp_body map[string]interface{}

	if res.StatusCode == 406 {
		fmt.Println("โทเคนหมดอายุโปรดเข้าสู่ระบบใหม่")
		Login()
	}
	content, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal([]byte(string(content)), &resp_body)
	if resp_body["success"] == true {
		return true
	}

	fmt.Println(resp_body["message"])
	return false
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
