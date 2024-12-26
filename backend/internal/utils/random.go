package utils

import (
	"math/rand"

	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
	"github.com/mozillazg/go-pinyin"
	"golang.org/x/crypto/bcrypt"
)

func GenerateRandomPassword(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*")
	random_password := make([]rune, length)
	for i := range random_password {
		random_password[i] = letters[rand.Intn(len(letters))]
	}
	return string(random_password)
}

func generateRandomChineseName() string {
	var commonSurnames = []string{
		"王", "李", "张", "刘", "陈", "杨", "赵", "黄", "周", "吴",
		"徐", "孙", "胡", "朱", "高", "林", "何", "郭", "马", "罗",
	}
	var commonNameCharacters = []string{
		"伟", "强", "芳", "敏", "静", "丽", "刚", "杰", "娟", "勇",
		"艳", "涛", "明", "军", "磊", "洋", "勇", "霞", "飞", "玲",
		"超", "华", "平", "辉", "梅", "鑫", "龙", "鹏", "玉", "斌",
		"庆", "建", "丹", "彬", "凤", "旭", "宁", "乐", "成", "欣",
	}

	surname := commonSurnames[rand.Intn(len(commonSurnames))]
	nameLength := rand.Intn(2) + 1
	name := ""

	for i := 0; i < nameLength; i++ {
		name += commonNameCharacters[rand.Intn(len(commonNameCharacters))]
	}

	return surname + name
}

func generateRandomDigit() string {
	var digits = "0123456789"
	return string(digits[rand.Intn(len(digits))])
}

func generateRandomRole() string {
	var roles = []string{"普通助理", "资深助理", "黑心"}
	return roles[rand.Intn(len(roles))]
}

func GenerateRandomUser() (*models.User, error) {
	user := &models.User{}

	// generate full name
	user.FullName = generateRandomChineseName()

	// generate username
	pinyinArray := pinyin.LazyConvert(user.FullName, nil)

	username := ""
	for _, pinyin := range pinyinArray {
		subLength := rand.Intn(len(pinyin)) + 1
		username += pinyin[:subLength]
	}

	digitsLength := rand.Intn(3) + 1
	for i := 0; i < digitsLength; i++ {
		username += generateRandomDigit()
	}

	user.Username = username

	// generate email
	user.Email = username + "@mail2.sysu.edu.cn"

	// generate role
	user.Role = generateRandomRole()

	// assign a password
	password := "ecnc@8403"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = string(passwordHash)

	return user, nil
}
