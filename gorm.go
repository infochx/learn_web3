package main

import (
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Student struct {
	ID    int
	Name  string
	Age   int
	Grade string
}

type Account struct {
	ID      int
	Balance int
}

type Transaction struct {
	ID            int
	FromAccountId int
	ToAccountId   int
	Amount        int
}

func InitDB(dst ...interface{}) *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(dst...)

	return db
}

func main() {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		fmt.Println("connect db error!")
		panic(err)
	}

	// SQL语句练习
	// 题目1：基本CRUD操作
	// CRUDTest(db)
	// 题目2：事务语句
	txTest(db)

	// lesson02.Run(db)
	// lesson03.Run(db)
	// lesson03_02.Run(db)
	// lesson03_03.Run(db)
	// lesson03_04.Run(db)
	// lesson04.Run(db)
}

func txTest(db *gorm.DB) {
	// 假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和
	// transactions 表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
	// 要求 ：编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。
	// 在事务中，需要先检查账户 A 的余额是否足够，
	// 如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。
	// 如果余额不足，则回滚事务。
	db.AutoMigrate(&Account{})
	db.AutoMigrate(&Transaction{})
	accountA := Account{}
	accountB := Account{}
	db.First(&accountA, 1)
	db.First(&accountB, 2)

	err := db.Transaction(func(tx *gorm.DB) error {

		if re := db.Where("id =?", accountA.ID).Find(&accountA); re.Error != nil {
			return re.Error
		}
		if accountA.Balance < 100 {
			return errors.New("balance not enough！")
		}

		if er := db.Model(&accountA).Update("balance", accountA.Balance-100).Error; er != nil {
			return er
		}

		if er := db.Model(&accountB).Update("balance", accountB.Balance+100).Error; er != nil {
			return er
		}

		if er := db.Save(&Transaction{FromAccountId: accountA.ID, ToAccountId: accountB.ID, Amount: 100}).Error; er != nil {
			return er
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

}

func CRUDTest(db *gorm.DB) {
	db.AutoMigrate(&Student{})
	// 要求 ：编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
	student := Student{Name: "张三", Age: 20, Grade: "三年级"}
	db.Create(&student)
	// 编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
	students := []Student{}
	db.Where("age > ?", 18).Find(&students)
	// 编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"
	db.Debug().Model(&Student{}).Where("name =?", "张三").Select("grade").Updates(map[string]interface{}{"grade": "四年级"})
	// 编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
	db.Debug().Where("age <=?", 15).Delete(&Student{})
}
