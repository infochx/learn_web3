package main

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
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

type Employee struct {
	ID         int
	Name       string
	Department string
	Salary     int
}

type Book struct {
	ID     int
	Title  string
	Author string
	Price  int
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
	db.AutoMigrate(&Employee{})
	db.AutoMigrate(&Book{})

	// SQL语句练习
	// 题目1：基本CRUD操作
	// CRUDTest(db)
	// 题目2：事务语句
	// TransactionTest(db)

	// Sqlx入门
	// 题目1：使用SQL扩展库进行查询
	// SqlxTest1()

	// 题目2：实现类型安全映射
	// SqlxTest2()

	// 进阶gorm
	// 题目1：模型定义
	// GormTest(db)
	// 题目2：关联查询
	// GormTest2(db)

	// 题目3：钩子函数
	GormTest3(db)
}

func GormTest3(db *gorm.DB) {
	// 题目3：钩子函数
	// 继续使用博客系统的模型。
	// 要求 ：为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
	post := Post{Title: "testBeforeCreate", AuthorID: 5}
	if err := db.Debug().Create(&post).Error; err != nil {
		panic(err)
	}
	// 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
	comment := Comment{Content: "content", PostID: 1}
	db.Create(&comment)
	if err := db.Debug().Delete(&Comment{ID: comment.ID}).Error; err != nil {
		panic(err)
	}
}
func (comment *Comment) BeforeDelete(db *gorm.DB) (err error) {
	if err = db.First(&comment).Error; err != nil {
		return err
	}
	post := Post{}
	if err := db.Debug().Where("id = ?", comment.PostID).First(&post).Error; err != nil {
		return err
	}
	post.Count--
	if post.Count <= 0 {
		post.Count = 0
		post.Status = "无评论"
	} else {
		post.Status = "有评论"
	}
	if err := db.Debug().Model(&post).Select("count", "status").
		Updates(&Post{Count: post.Count, Status: post.Status}).Error; err != nil {
		return err
	}
	return
}

func (post *Post) BeforeCreate(db *gorm.DB) (err error) {
	user := User{}
	if err = db.Debug().Where("id =?", post.AuthorID).First(&user).Error; err != nil {
		panic(err)
	}
	user.Count++
	err = db.Debug().Model(&User{}).Where("id =?", user.ID).Update("count", user.Count).Error
	if err != nil {
		return err
	}
	return
}

func GormTest2(db *gorm.DB) {
	// 题目2：关联查询
	user := User{}
	user = User{Name: "zhangsan",
		Post: []Post{
			{Title: "title3",
				Comment: []Comment{
					{Content: "111"}, {Content: "222"}, {Content: "333"}, {Content: "444"},
				},
			},
			{Title: "title4",
				Comment: []Comment{{Content: "2-111"}}},
		},
	}
	db.Save(&user)
	user.ID = 0

	// 要求 ：编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
	if err := db.Preload("Post.Comment").Where("id =?", 4).First(&user).Error; err != nil {
		panic(err)
	}

	// 编写Go代码，使用Gorm查询评论数量最多的文章信息。
	var result struct {
		ID           int
		Title        string
		CommentCount int
	}
	db.Model(&Post{}).Select("posts.id,posts.title,count(comments.id) comment_count").
		Joins("left join comments on comments.post_id=posts.id").
		Group("comments.post_id").Order("comment_count desc").First(&result)
	fmt.Println(result)

}

type User struct {
	ID    int
	Name  string
	Count int    `default:"0"`
	Post  []Post `gorm:"foreignKey:AuthorID"`
}
type Post struct {
	ID       int
	Title    string
	AuthorID int
	Count    int       `default:"0"`
	Comment  []Comment `gorm:"foreignKey:PostID"`
	Status   string
}
type Comment struct {
	ID      int
	Content string
	PostID  int
}

func GormTest(db *gorm.DB) {
	// 题目1：模型定义
	// 假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
	// 要求 ：使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章），
	//  Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。编写Go代码，使用Gorm创建这些模型对应的数据库表。
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Post{})
	db.AutoMigrate(&Comment{})
}

func SqlxTest2() {
	db, err := sqlx.Connect("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True")
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	// 题目2：实现类型安全映射
	// 假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
	// 要求 ：定义一个 Book 结构体，包含与 books 表对应的字段。
	// 编写Go代码，使用Sqlx执行一个复杂的查询，
	// 例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。
	books := []Book{}
	err = db.Select(&books, "select * from books where price >?", 50)
	if err != nil {
		panic(err)
	}
}

func SqlxTest1() {
	// 题目1：使用SQL扩展库进行查询
	// 假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
	db, err := sqlx.Connect("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True")
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	// 要求 ：编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，
	// 并将结果映射到一个自定义的 Employee 结构体切片中。
	employees := []Employee{}
	err = db.Select(&employees, "select * from employees where department=?", "技术部")
	if err != nil {
		panic(err)
	}

	// 编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
	emp := Employee{}
	err = db.Get(&emp, "select * from employees order by salary desc limit 1")
	if err != nil {
		panic(err)
	}
}

func TransactionTest(db *gorm.DB) {
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

	err := db.Transaction(func(tx *gorm.DB) (err error) {
		if err = db.Where("id =?", accountA.ID).First(&accountA).Error; err != nil {
			return err
		}
		if accountA.Balance < 100 {
			return errors.New("balance not enough！")
		}

		if err = db.Model(&accountA).Update("balance", accountA.Balance-100).Error; err != nil {
			return err
		}

		if err = db.Model(&accountB).Update("balance", accountB.Balance+100).Error; err != nil {
			return err
		}

		if err = db.Save(&Transaction{FromAccountId: accountA.ID, ToAccountId: accountB.ID, Amount: 100}).Error; err != nil {
			return err
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
	db.Model(&Student{}).Where("name =?", "张三").Select("grade").Updates(map[string]interface{}{"grade": "四年级"})
	// 编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
	db.Debug().Where("age <=?", 15).Delete(&Student{})
}
