package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 登录凭据和网站信息
var (
	baseURL = "https://www.bugutv.vip"
	client  *http.Client
)

func init() {
	jar, _ := cookiejar.New(nil)
	client = &http.Client{
		Jar: jar,
	}
}

// 登录网站
func login(username, password string) bool {
	fmt.Println("准备登录...")

	// 预请求主页（获取 cookie）
	_, err := client.Get(baseURL)
	if err != nil {
		fmt.Println("预请求失败:", err)
		return false
	}
	time.Sleep(1 * time.Second)

	// 登录请求
	data := url.Values{}
	data.Set("action", "user_login")
	data.Set("username", username)
	data.Set("password", password)
	data.Set("rememberme", "1")

	resp, err := client.PostForm(baseURL+"/wp-admin/admin-ajax.php", data)
	if err != nil {
		fmt.Println("登录请求失败:", err)
		return false
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if strings.Contains(string(body), "登录成功") || strings.Contains(string(body), "\\u767b\\u5f55\\u6210\\u529f") {
		fmt.Println("登录成功")
		return true
	}

	fmt.Println("登录失败")
	return false
}

// 获取当前积分
func getPoint() (int, error) {
	resp, err := client.Get(baseURL + "/user")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	time.Sleep(1 * time.Second)

	re := regexp.MustCompile(`<span class="badge badge-warning-lighten"><i class="fas fa-coins"></i> (.*?)</span>`)
	match := re.FindStringSubmatch(string(body))
	if len(match) < 2 {
		return 0, fmt.Errorf("未找到积分信息")
	}
	point, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, err
	}
	return point, nil
}

// 签到请求
func check() {
	resp, _ := client.Get(baseURL + "/user")
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	time.Sleep(1 * time.Second)

	nonceRe := regexp.MustCompile(`data-nonce="(.*?)"`)
	nonceMatch := nonceRe.FindStringSubmatch(string(body))
	if len(nonceMatch) < 2 {
		fmt.Println("未获取到 data-nonce")
		return
	}
	nonce := nonceMatch[1]
	fmt.Println("准备签到，data-nonce: " + nonce)

	data := url.Values{}
	data.Set("action", "user_qiandao")
	data.Set("nonce", nonce)

	resp2, _ := client.PostForm(baseURL+"/wp-admin/admin-ajax.php", data)
	defer resp2.Body.Close()
	body2, _ := io.ReadAll(resp2.Body)
	time.Sleep(1 * time.Second)

	content := string(body2)
	if strings.Contains(content, "\\u4eca\\u65e5\\u5df2\\u7b7e\\u5230") {
		fmt.Println("今日已签到，请明日再来")
	} else if strings.Contains(content, "\\u7b7e\\u5230\\u6210\\u529f") {
		fmt.Println("签到成功，奖励已到账：1.0积分")
	} else {
		fmt.Println("签到失败，返回内容:", content)
	}
}

// 获取 wpnonce 并退出
func logout() {
	resp, _ := client.Get(baseURL + "/user")
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	wpnonceRe := regexp.MustCompile(`action=logout&amp;redirect_to=https%3A%2F%2Fwww.bugutv.vip&amp;_wpnonce=(.*?)"`)
	match := wpnonceRe.FindStringSubmatch(string(body))
	if len(match) < 2 {
		fmt.Println("未获取到 wpnonce，无法退出登录")
		return
	}
	wpnonce := match[1]

	logoutURL := fmt.Sprintf(baseURL+"/wp-login.php?action=logout&redirect_to=https%%3A%%2F%%2Fwww.bugutv.vip&_wpnonce=%s", wpnonce)
	_, err := client.Get(logoutURL)
	if err == nil {
		fmt.Println("退出登录成功")
	} else {
		fmt.Println("退出登录失败:", err)
	}
}

func main() {
	fmt.Println("开始运行 bugutv 自动签到脚本")
	uname := os.Getenv("BUGUTV_USERNAME")
	upassword := os.Getenv("BUGUTV_PASSWORD")
	if uname == "" || upassword == "" {
		fmt.Println("请设置环境变量 BUGUTV_USERNAME 和 BUGUTV_PASSWORD")
		return
	}

	for i := 0; i < 3; i++ {
		if i > 0 {
			fmt.Printf("尝试第 %d 次...\n", i+1)
		}
		if login(uname, upassword) {
			pointBefore, err := getPoint()
			if err != nil {
				fmt.Println("获取积分失败:", err)
				continue
			}

			check()

			pointAfter, err := getPoint()
			if err != nil {
				fmt.Println("获取积分失败:", err)
				continue
			}

			earned := pointAfter - pointBefore
			fmt.Println("***************布谷TV签到:结果统计***************")
			fmt.Printf("%s 本次获得积分: %d 个\n累计积分: %d 个\n", uname, earned, pointAfter)
			fmt.Println("**************************************************")

			logout()
			break
		}
		time.Sleep(10 * time.Second)
	}
}
