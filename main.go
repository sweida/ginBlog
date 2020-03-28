package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"math/rand"
	"time"
)

type User struct {
	gorm.Model
	Name string `gorm:"type:varchar(20);nor null"`
	Telephone string `gorm:"varchar(110):not null;unique"`
	Password string `gorm:"size:255;not null"`
}

func main() {
	db := InitDB()
	defer db.Close()

	r := gin.Default()
	r.POST("api/auth/register", func(ctx *gin.Context) {
		name := ctx.PostForm("name")
		telephone := ctx.PostForm("telephone")
		password := ctx.PostForm("password")

		//数据验证
		if len(telephone) != 11 {
			ctx.JSON(200, gin.H{
				"status": "error",
				"code": 422,
				"msg": "手机号必须为11位",
			})
			return
		}

		if len(password) < 6 {
			ctx.JSON(200, gin.H{
				"status": "error",
				"code": 422,
				"msg": "密码不能少于6位",
			})
			return
		}

		//自动添加用户名
		if len(name) == 0 {
			name = RandomString(10)
		}

		if isTelephoneExist(db, telephone) {
			ctx.JSON(200, gin.H{
				"status": "error",
				"code": 422,
				"msg": "用户已经存在",
			})
			return
		}

		//新建用户
		newUser := User {
			Name : name,
			Telephone : telephone,
			Password : password,
		}
		db.Create(&newUser)

		//创建用户
		ctx.JSON(200, gin.H{
			"status": "success",
			"msg": "注册成功",
		})
	})
	panic(r.Run())
}

//生成随机字符串
func RandomString(n int) string{
	var letters = []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

func InitDB() *gorm.DB {
	db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/ginBlog?parseTime=true")
	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}

	//driverName := "mysql"
	//host := "localhost"
	//port := "3306"
	//database := "ginBlog"
	//username := "root"
	//password := "root"
	//charser := "utf8"
	//args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%S&parseTime=true",
	//	username,
	//	password,
	//	host,
	//	port,
	//	database,
	//	charser)
	//
	//db, err := gorm.Open(driverName, args)
	//if err != nil {
	//	panic("failed to connect database, err: " + err.Error())
	//}
	//自动创建数据表
	db.AutoMigrate(&User{})

	return db
}
