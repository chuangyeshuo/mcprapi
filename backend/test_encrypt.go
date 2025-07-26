/*
- @Author: lidi10@staff.weibo.com

  - @Date: 2025-07-26 17:36:11
    "fmt"

    "mcprapi/backend/internal/pkg/encrypt"

- Copyright (c) 2023 by Weibo, All Rights Reserved.
*/
package main

import (
	"fmt"

	"mcprapi/backend/internal/pkg/encrypt"
)

func main() {
	password := "123456"
	hash := encrypt.GenerateHash(password)
	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Hash: %s\n", hash)

	// 验证密码
	isValid := encrypt.VerifyPassword(password, hash)
	fmt.Printf("Verify result: %t\n", isValid)

	// 验证数据库中的哈希值
	dbHash := "jGl25bVBBBW96Qi9Te4V37Fnqchz/Eu4qB2JKbpbGKw="
	isValidDB := encrypt.VerifyPassword(password, dbHash)
	fmt.Printf("Verify DB hash result: %t\n", isValidDB)
}
