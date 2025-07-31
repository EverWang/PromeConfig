package main

import (
	"fmt"
	"d/GITVIEW/PromeConfig/backend/pkg/password"
)

func main() {
	// 测试密码哈希和验证
	plainPassword := "password123"
	storedhash := "$2a$10$9XSH7PiSX2iFTIjHbBG1seqHj3Y.d/VwjNtB7FMnTZTsErJEEv0g2"
	
	fmt.Printf("Plain password: %s\n", plainPassword)
	fmt.Printf("Stored hash: %s\n", storedhash)
	
	// 验证密码
	isValid := password.Verify(storedhash, plainPassword)
	fmt.Printf("Password verification result: %t\n", isValid)
	
	// 生成新的哈希进行对比
	newHash, err := password.Hash(plainPassword)
	if err != nil {
		fmt.Printf("Error generating hash: %v\n", err)
		return
	}
	fmt.Printf("New hash: %s\n", newHash)
	
	// 验证新哈希
	isNewValid := password.Verify(newHash, plainPassword)
	fmt.Printf("New hash verification result: %t\n", isNewValid)
}