package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/skip2/go-qrcode"
)

type UserInfo struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("======================================")
		fmt.Println("📋 用户二维码生成工具")
		fmt.Println("======================================")
		fmt.Println("")
		fmt.Println("用法:")
		fmt.Println("  go run generate_user_qr.go <用户名> <姓名> [手机号]")
		fmt.Println("")
		fmt.Println("示例:")
		fmt.Println("  go run generate_user_qr.go zhangsan 张三")
		fmt.Println("  go run generate_user_qr.go lisi 李四 13800138000")
		fmt.Println("")
		fmt.Println("======================================")
		return
	}

	username := os.Args[1]
	name := os.Args[2]
	phone := ""
	if len(os.Args) > 3 {
		phone = os.Args[3]
	}

	// 创建用户信息
	userInfo := UserInfo{
		Username: username,
		Name:     name,
		Phone:    phone,
	}

	// 转换为JSON
	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		fmt.Printf("❌ JSON编码失败: %v\n", err)
		return
	}

	// 生成二维码
	qrCode, err := qrcode.Encode(string(jsonData), qrcode.Medium)
	if err != nil {
		fmt.Printf("❌ 二维码生成失败: %v\n", err)
		return
	}

	// 打印二维码（ASCII版本）
	qrCodeString := string(qrCode)
	fmt.Println("\n" + qrCodeString)

	// 保存到文件
	filename := fmt.Sprintf("%s_qrcode.png", username)
	err = qrcode.WriteFile(filename, qrcode.Medium, []byte(string(jsonData)))
	if err != nil {
		fmt.Printf("❌ 保存二维码失败: %v\n", err)
		return
	}

	fmt.Println("\n======================================")
	fmt.Printf("✅ 二维码生成成功！\n")
	fmt.Printf("📝 用户信息:\n")
	fmt.Printf("   用户名: %s\n", username)
	fmt.Printf("   姓名: %s\n", name)
	if phone != "" {
		fmt.Printf("   手机: %s\n", phone)
	}
	fmt.Printf("\n📁 文件保存: %s\n", filename)
	fmt.Printf("\n📱 二维码内容(JSON):\n")
	fmt.Printf("   %s\n", string(jsonData))
	fmt.Println("\n💡 使用方法:")
	fmt.Println("   1. 将二维码图片发送给用户")
	fmt.Println("   2. 在抽奖页面点击'扫码添加用户'")
	fmt.Println("   3. 扫描二维码即可添加用户到抽奖池")
	fmt.Println("======================================")
}
