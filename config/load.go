package config

import (
	"encoding/json"
	"os"
)
type User struct {
	Listen  string `json:"listen"`
}
func ReadConfig(filePath string) (User, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 文件不存在，创建一个新的JSON文件
		user := User{Listen: ":9090"}
		data, err := json.Marshal(user)
		if err != nil {
			return User{}, err
		}
		err = os.WriteFile(filePath, data, 0644)
		if err != nil {
			return User{}, err
		}
	}
	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return User{}, err
	}
	// 解析JSON数据
	var user User
	err = json.Unmarshal(data, &user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
	//如何调用
	//---------------------------------------
	// filePath := "user.json"
	// // 读取配置文件
	// user, err := ReadConfig(filePath)
	// if err != nil {
	// 	fmt.Println("读取配置文件失败：", err)
	// }
