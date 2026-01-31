package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"lottery-system/config"
	"lottery-system/models"
	"lottery-system/utils"
)

func main() {
	config.LoadConfig()
	config.InitDB()

	// 检查命令行参数
	if len(os.Args) < 2 {
		fmt.Println("======================================")
		fmt.Println("📋 批量创建用户工具")
		fmt.Println("======================================")
		fmt.Println("")
		fmt.Println("⚠️  重要提示:")
		fmt.Println("   - 创建的用户仅用于抽奖池")
		fmt.Println("   - 用户无法登录系统（仅管理员可登录）")
		fmt.Println("   - 管理员代为用户执行抽奖操作")
		fmt.Println("")
		fmt.Println("用法:")
		fmt.Println("  go run batch_create_users.go <用户数据文件>")
		fmt.Println("")
		fmt.Println("文件格式（每行一个用户）:")
		fmt.Println("  用户名,密码,姓名")
		fmt.Println("")
		fmt.Println("示例:")
		fmt.Println("  zhangsan,123456,张三")
		fmt.Println("  lisi,123456,李四")
		fmt.Println("  wangwu,123456,王五")
		fmt.Println("")
		fmt.Println("包含手机号（可选）:")
		fmt.Println("  zhangsan,123456,张三,13800138001")
		fmt.Println("  lisi,123456,李四,13800138002")
		fmt.Println("")
		fmt.Println("======================================")
		return
	}

	filename := os.Args[1]

	// 检查文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Fatalf("❌ 文件不存在: %s", filename)
	}

	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("❌ 无法打开文件: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// 默认公司ID为1（默认公司）
	companyID := 1
	successCount := 0
	failCount := 0

	fmt.Printf("🔄 开始批量创建用户...\n")
	fmt.Printf("📂 公司ID: %d\n", companyID)
	fmt.Println("")

	lineNum := 0
	var users []models.User

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行
		if line == "" {
			continue
		}

		// 跳过注释行（以#开头）
		if strings.HasPrefix(line, "#") {
			fmt.Printf("⏭️  第%d行: 跳过注释\n", lineNum)
			continue
		}

		// 解析数据
		parts := strings.Split(line, ",")

		if len(parts) < 3 {
			fmt.Printf("⚠️  第%d行: 格式错误（需要: 用户名,密码,姓名[,手机号]）\n", lineNum)
			failCount++
			continue
		}

		username := strings.TrimSpace(parts[0])
		password := strings.TrimSpace(parts[1])
		name := strings.TrimSpace(parts[2])

		// 可选：手机号
		phone := ""
		if len(parts) > 3 {
			phone = strings.TrimSpace(parts[3])
		}

		// 验证数据
		if username == "" || password == "" || name == "" {
			fmt.Printf("⚠️  第%d行: 数据不完整\n", lineNum)
			failCount++
			continue
		}

		// 检查密码长度
		if len(password) < 6 {
			fmt.Printf("⚠️  第%d行: %s - 密码太短（至少6位）\n", lineNum, username)
			failCount++
			continue
		}

		// 检查用户名是否已存在，允许重名
		// 如果有手机号，用 (username, phone) 判断；如果没有手机号，允许重名
		finalUsername := username
		var existingUsers []models.User
		query := config.DB.Where("company_id = ?", companyID).Where("username = ?", username)

		if phone != "" {
			query = query.Where("phone = ?", phone)
		}

		if err := query.Find(&existingUsers).Error; err == nil && len(existingUsers) > 0 {
			// 有手机号的用户：认为已存在
			if phone != "" {
				fmt.Printf("⚠️  第%d行: %s (%s) - 该用户名和手机号的用户已存在\n", lineNum, username, name)
				failCount++
				continue
			}

			// 没有手机号但有重名用户：自动添加序号
			var count int64
			config.DB.Model(&models.User{}).Where("username = ? AND company_id = ?", username, companyID).Count(&count)
			finalUsername = fmt.Sprintf("%s_%d", username, count+1)
			fmt.Printf("ℹ️  第%d行: %s - 检测到重名，自动修改为 %s\n", lineNum, username, finalUsername)
		}

		// 加密密码
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			fmt.Printf("❌ 第%d行: %s - 密码加密失败\n", lineNum, username)
			failCount++
			continue
		}

		// 创建用户
		user := models.User{
			CompanyID: companyID,
			Username:  username,
			Password:  hashedPassword,
			Role:      models.RoleUser,
			Name:      name,
			Phone:     phone,
			HasDrawn:  false,
		}

		users = append(users, user)
		fmt.Printf("✅ 第%d行: %s (%s) - 准备创建\n", lineNum, username, name)
	}

	// 所有数据验证通过后，批量插入数据库
	fmt.Printf("\n💾 开始保存到数据库...\n")

	for _, user := range users {
		if err := config.DB.Create(&user).Error; err != nil {
			fmt.Printf("❌ %s - 创建失败: %v\n", user.Username, err)
			failCount++
		} else {
			fmt.Printf("✅ %s - 创建成功 (ID:%d)\n", user.Username, user.ID)
			successCount++
		}
	}

	fmt.Println("")
	fmt.Println("======================================")
	fmt.Printf("✅ 批量创建完成\n")
	fmt.Printf("📊 统计:\n")
	fmt.Printf("   成功: %d 个\n", successCount)
	fmt.Printf("   失败: %d 个\n", failCount)
	fmt.Printf("   总计: %d 个\n", successCount+failCount)
	fmt.Println("======================================")

	if successCount > 0 {
		fmt.Println("\n📋 用户信息:")
		fmt.Println("   ⚠️  用户已创建，但无法登录系统")
		fmt.Println("   ℹ️  仅管理员可以登录并代为抽奖")
		fmt.Println("   💡  请在管理后台查看用户列表")
	}
}
