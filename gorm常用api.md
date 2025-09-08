# GORM 常用 API 参考文档

## 目录
- [数据库连接](#数据库连接)
- [模型定义](#模型定义)
- [CRUD 操作](#crud-操作)
- [查询操作](#查询操作)
- [更新操作](#更新操作)
- [删除操作](#删除操作)
- [关联操作](#关联操作)
- [事务操作](#事务操作)
- [钩子函数](#钩子函数)
- [迁移操作](#迁移操作)

## 数据库连接

### 连接数据库
```go
import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

// MySQL
dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// PostgreSQL
import "gorm.io/driver/postgres"
dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

// SQLite
import "gorm.io/driver/sqlite"
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
```

### 配置选项
```go
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    SkipDefaultTransaction: true,  // 跳过默认事务
    NamingStrategy: schema.NamingStrategy{
        TablePrefix: "t_",   // 表名前缀
        SingularTable: true, // 使用单数表名
    },
    Logger: logger.Default.LogMode(logger.Info), // 日志级别
})
```

## 模型定义

### 基础模型
```go
type User struct {
    ID        uint           `gorm:"primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
    Name      string         `gorm:"size:255;not null"`
    Email     string         `gorm:"uniqueIndex"`
    Age       int
}

// 使用 gorm.Model
type User struct {
    gorm.Model
    Name  string
    Email string
}
```

### 字段标签
```go
type User struct {
    ID       uint   `gorm:"primaryKey;autoIncrement"`
    Name     string `gorm:"size:255;not null;index"`
    Email    string `gorm:"uniqueIndex;size:100"`
    Password string `gorm:"size:255;not null"`
    Status   int    `gorm:"default:1;comment:用户状态"`
    Profile  string `gorm:"type:text"`
}
```

## CRUD 操作

### 创建 (Create)
```go
// 创建单条记录
user := User{Name: "张三", Email: "zhangsan@example.com"}
result := db.Create(&user)
// user.ID 返回插入数据的主键
// result.Error 返回 error
// result.RowsAffected 返回插入记录的条数

// 批量创建
users := []User{
    {Name: "张三", Email: "zhangsan@example.com"},
    {Name: "李四", Email: "lisi@example.com"},
}
db.Create(&users)

// 创建时指定字段
db.Select("name", "email").Create(&user)

// 创建时忽略字段
db.Omit("password").Create(&user)
```

### 读取 (Read)
```go
// 获取第一条记录（主键升序）
var user User
db.First(&user)
// SELECT * FROM users ORDER BY id LIMIT 1;

// 获取最后一条记录（主键降序）
db.Last(&user)

// 随机获取一条记录
db.Take(&user)

// 根据主键检索
db.First(&user, 10)
// SELECT * FROM users WHERE id = 10;

db.First(&user, "10")
// SELECT * FROM users WHERE id = 10;

// 检索全部对象
var users []User
db.Find(&users)
// SELECT * FROM users;

db.Find(&users, []int{1,2,3})
// SELECT * FROM users WHERE id IN (1,2,3);
```

## 查询操作

### Where 条件
```go
// String 条件
db.Where("name = ?", "张三").First(&user)
db.Where("name <> ?", "张三").Find(&users)
db.Where("name IN ?", []string{"张三", "李四"}).Find(&users)
db.Where("name LIKE ?", "%三%").Find(&users)
db.Where("name = ? AND age >= ?", "张三", "22").Find(&users)

// Struct 条件
db.Where(&User{Name: "张三", Age: 20}).First(&user)
// SELECT * FROM users WHERE name = "张三" AND age = 20 ORDER BY id LIMIT 1;

// Map 条件
db.Where(map[string]interface{}{"name": "张三", "age": 20}).Find(&users)

// 主键切片条件
db.Where([]int64{20, 21, 22}).Find(&users)
// SELECT * FROM users WHERE id IN (20, 21, 22);
```

### Not 条件
```go
db.Not("name = ?", "张三").First(&user)
db.Not(map[string]interface{}{"name": []string{"张三", "李四"}}).Find(&users)
db.Not([]int64{1,2,3}).First(&user)
```

### Or 条件
```go
db.Where("role = ?", "admin").Or("role = ?", "super_admin").Find(&users)
// SELECT * FROM users WHERE role = 'admin' OR role = 'super_admin';

db.Where("name = 'jinzhu'").Or(User{Name: "jinzhu 2", Age: 18}).Find(&users)
```

### 选择特定字段
```go
db.Select("name", "age").Find(&users)
db.Select([]string{"name", "age"}).Find(&users)

// 排除字段
db.Omit("created_at", "updated_at").Find(&users)
```

### 排序
```go
db.Order("age desc, name").Find(&users)
db.Order("age desc").Order("name").Find(&users)
```

### 限制和偏移
```go
db.Limit(3).Find(&users)
// SELECT * FROM users LIMIT 3;

db.Limit(10).Offset(5).Find(&users)
// SELECT * FROM users OFFSET 5 LIMIT 10;

// 取消 limit 条件，用 -1
db.Limit(10).Find(&users1).Limit(-1).Find(&users2)
```

### 分组和Having
```go
type result struct {
    Date  time.Time
    Total int
}

db.Model(&User{}).Select("date(created_at) as date, sum(age) as total").Group("date(created_at)").Having("sum(age) > ?", 100).Find(&results)
```

### 连接查询
```go
type result struct {
    Name  string
    Email string
}

db.Model(&User{}).Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Find(&results)

// 带参数的 Joins
db.Joins("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "jinzhu@example.org").Find(&users)

// 多个 Joins
db.Joins("JOIN emails ON emails.user_id = users.id").Joins("JOIN credit_cards ON credit_cards.user_id = users.id").Find(&users)
```

## 更新操作

### Save 保存所有字段
```go
db.First(&user)
user.Name = "李四"
user.Age = 100
db.Save(&user)
// UPDATE users SET name='李四', age=100, updated_at = '2013-11-17 21:34:10' WHERE id=111;
```

### Update 更新单个列
```go
// 条件更新
db.Model(&user).Update("name", "hello")
// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111;

db.Model(&user).Where("active = ?", true).Update("name", "hello")
// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111 AND active=true;
```

### Updates 更新多个列
```go
// 根据 struct 更新属性，只会更新非零值的字段
db.Model(&user).Updates(User{Name: "hello", Age: 18, Active: false})
// UPDATE users SET name='hello', age=18, updated_at = '2013-11-17 21:34:10' WHERE id = 111;

// 根据 map 更新属性
db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET name='hello', age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;
```

### 更新选定字段
```go
// Select 指定字段
db.Model(&user).Select("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET name='hello' WHERE id=111;

// Omit 忽略字段
db.Model(&user).Omit("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;
```

### 批量更新
```go
db.Model(User{}).Where("role = ?", "admin").Updates(User{Name: "hello", Age: 18})
// UPDATE users SET name='hello', age=18 WHERE role = 'admin';

db.Table("users").Where("id IN ?", []int{10, 11}).Updates(map[string]interface{}{"name": "hello", "age": 18})
// UPDATE users SET name='hello', age=18 WHERE id IN (10, 11);
```

## 删除操作

### 删除记录
```go
// 删除一条记录
db.Delete(&user)
// DELETE from users where id=10;

// 带额外条件的删除
db.Where("name = ?", "jinzhu").Delete(&user)
// DELETE from users where id=10 AND name = "jinzhu";

// 根据主键删除
db.Delete(&User{}, 10)
// DELETE FROM users WHERE id = 10;

db.Delete(&User{}, "10")
// DELETE FROM users WHERE id = 10;

db.Delete(&users, []int{1,2,3})
// DELETE FROM users WHERE id IN (1,2,3);
```

### 批量删除
```go
db.Where("email LIKE ?", "%jinzhu%").Delete(&User{})
// DELETE from users where email LIKE "%jinzhu%";

db.Delete(&User{}, "email LIKE ?", "%jinzhu%")
// DELETE from users where email LIKE "%jinzhu%";
```

### 软删除
```go
// 如果模型包含了 DeletedAt 字段，它将自动获得软删除的能力
type User struct {
    ID      uint
    Deleted gorm.DeletedAt
    Name    string
}

db.Delete(&user)
// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE id = 111;

// 查询时会忽略被软删除的记录
db.Where("age = 20").Find(&user)
// SELECT * FROM users WHERE age = 20 AND deleted_at IS NULL;

// Unscoped 方法可以查询被软删除的记录
db.Unscoped().Where("age = 20").Find(&users)
// SELECT * FROM users WHERE age = 20;

// 永久删除
db.Unscoped().Delete(&user)
// DELETE FROM users WHERE id=10;
```

## 关联操作

### 一对一 (Has One)
```go
type User struct {
    gorm.Model
    Name      string
    CreditCard CreditCard
}

type CreditCard struct {
    gorm.Model
    Number string
    UserID uint
}

// 查询
db.Preload("CreditCard").Find(&users)

// 创建关联
db.Model(&user).Association("CreditCard").Append(&CreditCard{Number: "411111111111"})
```

### 一对多 (Has Many)
```go
type User struct {
    gorm.Model
    Name       string
    CreditCards []CreditCard
}

// 预加载
db.Preload("CreditCards").Find(&users)

// 嵌套预加载
db.Preload("CreditCards.Bank").Find(&users)
```

### 多对多 (Many To Many)
```go
type User struct {
    gorm.Model
    Languages []Language `gorm:"many2many:user_languages;"`
}

type Language struct {
    gorm.Model
    Name string
}

// 预加载
db.Preload("Languages").Find(&users)

// 添加关联
db.Model(&user).Association("Languages").Append([]Language{languageZH, languageEN})
```

## 事务操作

### 手动事务
```go
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

if err := tx.Error; err != nil {
    return err
}

if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
    tx.Rollback()
    return err
}

if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
    tx.Rollback()
    return err
}

return tx.Commit().Error
```

### 事务函数
```go
db.Transaction(func(tx *gorm.DB) error {
    // 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
    if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
        // 返回任何错误都会回滚事务
        return err
    }

    if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
        return err
    }

    // 返回 nil 提交事务
    return nil
})
```

### 嵌套事务
```go
db.Transaction(func(tx *gorm.DB) error {
    tx.Create(&user1)

    tx.Transaction(func(tx2 *gorm.DB) error {
        tx2.Create(&user2)
        return errors.New("rollback user2") // 回滚 user2 但不会回滚 user1
    })

    tx.Transaction(func(tx2 *gorm.DB) error {
        tx2.Create(&user3)
        return nil
    })

    return nil
})
```

## 钩子函数

### 创建钩子
```go
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
    u.UUID = uuid.New()
    if u.Role == "admin" {
        return errors.New("invalid role")
    }
    return
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
    if u.ID == 1 {
        tx.Model(u).Update("role", "admin")
    }
    return
}
```

### 更新钩子
```go
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
    if u.readonly() {
        err = errors.New("read only user")
    }
    return
}

func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
    if u.Confirmed {
        tx.Model(&u).Update("confirmed_at", time.Now())
    }
    return
}
```

### 删除钩子
```go
func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
    if u.Role == "admin" {
        return errors.New("admin user not allowed to delete")
    }
    return
}

func (u *User) AfterDelete(tx *gorm.DB) (err error) {
    logger.Info("User deleted", u.ID)
    return
}
```

## 迁移操作

### 自动迁移
```go
db.AutoMigrate(&User{})
db.AutoMigrate(&User{}, &Product{}, &Order{})

// 创建表时添加表后缀
db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{})
```

### 检查表是否存在
```go
// 检查模型 `User` 对应的表是否存在
db.Migrator().HasTable(&User{})
db.Migrator().HasTable("users")

// 为当前数据库创建表
db.Migrator().CreateTable(&User{})

// 删除模型 `User` 对应的表
db.Migrator().DropTable(&User{})
db.Migrator().DropTable("users")

// 重命名表
db.Migrator().RenameTable(&User{}, &UserInfo{})
db.Migrator().RenameTable("users", "user_infos")
```

### 列操作
```go
type User struct {
    Name string
}

// 添加列
db.Migrator().AddColumn(&User{}, "Name")
// 删除列
db.Migrator().DropColumn(&User{}, "Name")
// 修改列
db.Migrator().AlterColumn(&User{}, "Name")
// 检查列是否存在
db.Migrator().HasColumn(&User{}, "Name")
// 重命名列
db.Migrator().RenameColumn(&User{}, "Name", "Name2")
```

### 索引操作
```go
type User struct {
    gorm.Model
    Name  string `gorm:"size:255;index:idx_name,unique"`
    Email string `gorm:"size:255;index:idx_email"`
}

// 创建索引
db.Migrator().CreateIndex(&User{}, "Name")
db.Migrator().CreateIndex(&User{}, "idx_user_name")

// 删除索引
db.Migrator().DropIndex(&User{}, "Name")
db.Migrator().DropIndex(&User{}, "idx_user_name")

// 检查索引是否存在
db.Migrator().HasIndex(&User{}, "Name")
db.Migrator().HasIndex(&User{}, "idx_user_name")

// 重命名索引
db.Migrator().RenameIndex(&User{}, "Name", "Name2")
```

## 常用查询示例

### 分页查询
```go
func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        page, _ := strconv.Atoi(r.URL.Query().Get("page"))
        if page == 0 {
            page = 1
        }

        pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
        switch {
        case pageSize > 100:
            pageSize = 100
        case pageSize <= 0:
            pageSize = 10
        }

        offset := (page - 1) * pageSize
        return db.Offset(offset).Limit(pageSize)
    }
}

db.Scopes(Paginate(r)).Find(&users)
```

### 计数查询
```go
var count int64
db.Model(&User{}).Where("name = ?", "jinzhu").Count(&count)

var users []User
db.Model(&User{}).Where("name = ?", "jinzhu").Find(&users).Count(&count)
```

### 原生SQL
```go
type Result struct {
    ID   int
    Name string
    Age  int
}

var result Result
db.Raw("SELECT id, name, age FROM users WHERE name = ?", "jinzhu").Scan(&result)

var results []Result
db.Raw("SELECT id, name, age FROM users WHERE name = ?", "jinzhu").Scan(&results)

// 执行原生SQL
db.Exec("DROP TABLE users")
db.Exec("UPDATE orders SET shipped_at=? WHERE id IN ?", time.Now(), []int64{1,2,3})
```

## 错误处理

### 常见错误
```go
import "errors"

// ErrRecordNotFound 没有找到记录错误
errors.Is(err, gorm.ErrRecordNotFound)

// 检查错误
if err := db.Where("name = ?", "jinzhu").First(&user).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        // 没有找到记录
    } else {
        // 其他错误
    }
}

// 获取受影响的行数
result := db.Create(&user)
result.RowsAffected // 返回插入记录的条数
result.Error        // 返回 error
```

## 性能优化

### 预加载
```go
// 预加载 Orders
db.Preload("Orders").Find(&users)

// 预加载 Orders 和 Orders 的 OrderItems
db.Preload("Orders.OrderItems").Find(&users)

// 条件预加载
db.Preload("Orders", "state NOT IN (?)", "cancelled").Find(&users)

// 自定义预加载 SQL
db.Preload("Orders", func(db *gorm.DB) *gorm.DB {
    return db.Order("orders.amount DESC")
}).Find(&users)
```

### 选择字段
```go
// 只选择需要的字段
db.Select("name", "age").Find(&users)

// 排除不需要的字段
db.Omit("password").Find(&users)
```

### 批量操作
```go
// 批量插入
var users = []User{{Name: "jinzhu1"}, {Name: "jinzhu2"}, {Name: "jinzhu3"}}
db.Create(&users)

// 批量插入时指定批次大小
db.CreateInBatches(users, 100)
```

这个文档涵盖了GORM的大部分常用API，可以作为日常开发的参考手册。