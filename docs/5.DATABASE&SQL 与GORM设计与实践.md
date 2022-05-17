# 理解database/sql

## 基本用法

```go
package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID   int
	Name string
}

func main() {
	// import driver 实现
	// 使用driver + DSN 初始化 DB 连接
	db, err := sql.Open("mysql", "user:password@tcp(127.0.01:3306)/hello")

	// 执行一条sql，通过rows取回返回的数据处理完毕，需要释放链接
	rows, err := db.Query("select id, name from users where id = ?", 1)
	if err != nil {
		// XXX
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name)

		// 数据、错误处理
		if err != nil {
			// XXX
		}

		users = append(users, user)
	}

	// 错误处理
	if rows.Err() != nil {
		// XXX
	}
}
```

## 设计原理

![image-20220515200403845](images/image-20220515200403845.png)

```go
// DB is a database handle representing a pool of zero or more
// underlying connections. It's safe for concurrent use by multiple
// goroutines.
//
// The sql package creates and frees connections automatically; it
// also maintains a free pool of idle connections. If the database has
// a concept of per-connection state, such state can be reliably observed
// within a transaction (Tx) or connection (Conn). Once DB.Begin is called, the
// returned Tx is bound to a single connection. Once Commit or
// Rollback is called on the transaction, that transaction's
// connection is returned to DB's idle connection pool. The pool size
// can be controlled with SetMaxIdleConns.
type DB struct {
   // Atomic access only. At top of struct to prevent mis-alignment
   // on 32-bit platforms. Of type time.Duration.
   waitDuration int64 // Total time waited for new connections.

   connector driver.Connector
   // numClosed is an atomic counter which represents a total number of
   // closed connections. Stmt.openStmt checks it before cleaning closed
   // connections in Stmt.css.
   numClosed uint64

   mu           sync.Mutex // protects following fields
   freeConn     []*driverConn
   connRequests map[uint64]chan connRequest
   nextRequest  uint64 // Next key to use in connRequests.
   numOpen      int    // number of opened and pending open connections
   // Used to signal the need for new connections
   // a goroutine running connectionOpener() reads on this chan and
   // maybeOpenNewConnections sends on the chan (one send per needed connection)
   // It is closed during db.Close(). The close tells the connectionOpener
   // goroutine to exit.
   openerCh          chan struct{}
   closed            bool
   dep               map[finalCloser]depSet
   lastPut           map[*driverConn]string // stacktrace of last conn's put; debug only
   maxIdleCount      int                    // zero means defaultMaxIdleConns; negative means 0
   maxOpen           int                    // <= 0 means unlimited
   maxLifetime       time.Duration          // maximum amount of time a connection may be reused
   maxIdleTime       time.Duration          // maximum amount of time a connection may be idle before being closed
   cleanerCh         chan struct{}
   waitCount         int64 // Total number of connections waited for.
   maxIdleClosed     int64 // Total number of connections closed due to idle count.
   maxIdleTimeClosed int64 // Total number of connections closed due to idle time.
   maxLifetimeClosed int64 // Total number of connections closed due to max connection lifetime limit.

   stop func() // stop cancels the connection opener.
}
```

- 操作过程实现

```go
// QueryContext executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
   var rows *Rows
   var err error
   for i := 0; i < maxBadConnRetries; i++ {
      rows, err = db.query(ctx, query, args, cachedOrNewConn)
      if err != driver.ErrBadConn {
         break
      }
   }
   if err == driver.ErrBadConn {
      return db.query(ctx, query, args, alwaysNewConn)
   }
   return rows, err
}
```

- 使用Driver连接

```go
// Driver is the interface that must be implemented by a database
// driver.
//
// Database drivers may implement DriverContext for access
// to contexts and to parse the name only once for a pool of connections,
// instead of once per connection.
type Driver interface {
   // Open returns a new connection to the database.
   // The name is a string in a driver-specific format.
   //
   // Open may return a cached connection (one previously
   // closed), but doing so is unnecessary; the sql package
   // maintains a pool of idle connections for efficient re-use.
   //
   // The returned connection is only used by one goroutine at a
   // time.
   Open(name string) (Conn, error)
}
```

- 自定义注册驱动，遇见重名的驱动会panic

```go
// Register makes a database driver available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, driver driver.Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("sql: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("sql: Register called twice for driver " + name)
	}
	drivers[name] = driver
}
```

```go
package main

import (
	"database/sql"
	"database/sql/driver"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLDriver struct {
}

func (m MySQLDriver) Open(name string) (driver.Conn, error) {
	//TODO implement me
	panic("implement me")
}

func main() {
	sql.Open("mysql", "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&&parseTime=True&loc=Local")
}

// 注册 Driver
func init() {
	sql.Register("mysql", &MySQLDriver{})
}
```

- Driver 连接接口2

```go
// A Connector represents a driver in a fixed configuration
// and can create any number of equivalent Conns for use
// by multiple goroutines.
//
// A Connector can be passed to sql.OpenDB, to allow drivers
// to implement their own sql.DB constructors, or returned by
// DriverContext's OpenConnector method, to allow drivers
// access to context and to avoid repeated parsing of driver
// configuration.
//
// If a Connector implements io.Closer, the sql package's DB.Close
// method will call Close and return error (if any).
type Connector interface {
	// Connect returns a connection to the database.
	// Connect may return a cached connection (one previously
	// closed), but doing so is unnecessary; the sql package
	// maintains a pool of idle connections for efficient re-use.
	//
	// The provided context.Context is for dialing purposes only
	// (see net.DialContext) and should not be stored or used for
	// other purposes. A default timeout should still be used
	// when dialing as a connection pool may call Connect
	// asynchronously to any query.
	//
	// The returned connection is only used by one goroutine at a
	// time.
	Connect(context.Context) (Conn, error)

	// Driver returns the underlying Driver of the Connector,
	// mainly to maintain compatibility with the Driver method
	// on sql.DB.
	Driver() Driver
}
```

```go
// OpenDB opens a database using a Connector, allowing drivers to
// bypass a string based data source name.
//
// Most users will open a database via a driver-specific connection
// helper function that returns a *DB. No database drivers are included
// in the Go standard library. See https://golang.org/s/sqldrivers for
// a list of third-party drivers.
//
// OpenDB may just validate its arguments without creating a connection
// to the database. To verify that the data source name is valid, call
// Ping.
//
// The returned DB is safe for concurrent use by multiple goroutines
// and maintains its own pool of idle connections. Thus, the OpenDB
// function should be called just once. It is rarely necessary to
// close a DB.
func OpenDB(c driver.Connector) *DB {
	ctx, cancel := context.WithCancel(context.Background())
	db := &DB{
		connector:    c,
		openerCh:     make(chan struct{}, connectionRequestQueueSize),
		lastPut:      make(map[*driverConn]string),
		connRequests: make(map[uint64]chan connRequest),
		stop:         cancel,
	}

	go db.connectionOpener(ctx)

	return db
}
```

```go
package main

import (
   "database/sql"
   "github.com/go-sql-driver/mysql"
)

func main() {
   connector, err := mysql.NewConnector(&mysql.Config{
      User:      "gorm",
      Passwd:    "gorm",
      Net:       "tcp",
      Addr:      "127.0.0.1:3306",
      DBName:    "gorm",
      ParseTime: true,
   })

   db := sql.OpenDB(connector)
}
```

- DB连接的几种类型

![image-20220515205628063](images/image-20220515205628063.png)

- 处理返回数据的几种方式

![image-20220515212207847](images/image-20220515212207847.png)

- 实现这些接口解析数据

```go
// Rows is an iterator over an executed query's results.
type Rows interface {
   // Columns returns the names of the columns. The number of
   // columns of the result is inferred from the length of the
   // slice. If a particular column name isn't known, an empty
   // string should be returned for that entry.
   Columns() []string

   // Close closes the rows iterator.
   Close() error

   // Next is called to populate the next row of data into
   // the provided slice. The provided slice will be the same
   // size as the Columns() are wide.
   //
   // Next should return io.EOF when there are no more rows.
   //
   // The dest should not be written to outside of Next. Care
   // should be taken when closing Rows not to modify
   // a buffer held in dest.
   Next(dest []Value) error
}
```

# GORM使用简介

## 背景知识

“设计简洁，功能强大，自由扩展的全功能ORM”

## 基本用法

```go
package main

import (
   "gorm.io/driver/mysql"
   "gorm.io/gorm"
)

type User struct {
   ID   int
   Name string
}

func main() {
   db, _ := gorm.Open(mysql.Open("user:password@tcp(127.0.0.1:3306/hello"))
   
   var users []User
   _ = db.Select("id", "name").Find(&users, 1).Error
}
```

## CURD

```go
// 创建
user := User{Name: "wuyou", Age: 18, Birthday: time.Now()}
// 通过数据的指针创建
result := db.Create(&user)

// 返回主键 last insert id
print(user.ID)
// 返回 error
print(result.Error)
// 返回影响的行数
print(result.RowsAffected)

// 批量创建
users = []User{{Name: "wuyou"}, {Name: "wuyou"}, {Name: "wuyou"}}
db.Create(&users)
db.CreateInBatches(users, 100)
for _, user := range users {
   print(user.ID)
}
```

```go
// 读取
var product Product
// 查询id为1
db.First(&product, 1)
// 查询code为L1212的product
db.First(&product, "code = ?", "L1212")

result = db.Find(&users, []int{1, 2, 3})
// 找到的记录数
print(result.RowsAffected)
// First, Last, Take 查不到数据
errors.Is(result.Error, gorm.ErrRecordNotFound)

// 更新某个字段
db.Model(&product).Update("Price", 2000)
db.Model(&product).UpdateColumn("Price", 2000)

// 更新多个字段
db.Model(&product).Updates(Product{Price: 2000, Code: "L1212"})
db.Model(&product).Updates(map[string]interface{}{"Price": 2000, "Code": "L1212"})

// 批量更新
db.Model(&Product{}).Where("price < ?", 2000).Updates(map[string]interface{}{"Price": 2000})

// 删除 product
db.Delete(&product)
```

## 模型定义

```go
// User 模型定义
type User struct {
   ID           uint
   Name         string
   Email        *string
   Age          uint8
   Birthday     *time.Time
   MemberNumber sql.NullTime
   ActivateAt   sql.NullTime
   CreateAt     time.Time
   UpdateAt     time.Time
   DeletedAt    gorm.DeletedAt `gorm:"index"`
}
```

- 等价的model定义

```go
// User1 模型定义
type User1 struct {
   gorm.Model
   ID           uint
   Name         string
   Email        *string
   Age          uint8
   Birthday     *time.Time
   MemberNumber sql.NullTime
   ActivateAt   sql.NullTime
}

// Model a basic GoLang struct which includes the following fields: ID, CreatedAt, UpdatedAt, DeletedAt
// It may be embedded into your model or you may build your own model without it
//    type User struct {
//      gorm.Model
//    }
type Model struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt DeletedAt `gorm:"index"`
}
```

## 惯例约定

> 约定优于配置

- 表名为 struct name 的 snake_cases 复数格式
- 字段名为 field name  的snake_case 单数格式
- ID / Id 字段为主键，如果为数字，则为自增主键
- CreateAt 字段，创建时，保存当前时间
- UpdateAt字段，创建、更新时，保存当前时间
- gorm.DeletedAt字段，默认开启 soft delete 模式

## 关联介绍

```go
import (
   "database/sql"
   "errors"
   "gorm.io/driver/mysql"
   "gorm.io/gorm"
   "time"
)

// User 模型定义
type User struct {
   ID           uint
   Name         string
   Email        *string
   Age          uint8
   Birthday     time.Time
   MemberNumber sql.NullTime
   ActivateAt   sql.NullTime
   CreateAt     time.Time
   UpdateAt     time.Time
   DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// User1 模型定义
type User1 struct {
   gorm.Model
   ID           uint
   Name         string
   Email        *string
   Age          uint8
   Birthday     *time.Time
   MemberNumber sql.NullTime
   ActivateAt   sql.NullTime
}

type Product struct {
   Price int
   Code  string
}
```

```go
func main() {
   db, _ := gorm.Open(mysql.Open("user:password@tcp(127.0.0.1:3306/hello"))

   // 保存用户及其关联 （Upsert）
   db.Save(&User{
      Name:      "wuyou",
      Languages: []Language{{Name: "zh-CN"}, {Name: "en-US"}},
   })

   user := User{}
   languages := Language{}
   // 关联模式
   langAssociation := db.Model(&user).Association("Languages")
   // 查询关联
   langAssociation.Find(&languages)
   // 将汉语，英语添加到用户掌握的语言中
   langAssociation.Append([]Language{languageZH, languageEN})
   // 把用户掌握的语言替换为汉语，德语
   langAssociation.Replace([]Language{languageZH, languageDE})
   // 删除用户掌握的两个语言
   langAssociation.Delete(languageZH, languageEN)
   // 删除用户掌握的语言
   langAssociation.Clear()
   // 返回用户所掌握的语言的数量
   langAssociation.Count()

   var user1 User
   var user2 User
   var user3 User
   users := []User{user1, user2, user3}
   // 批量模式 Append, Replace
   langAssociation = db.Model(&users).Association("Languages")

   db.Model(&users).Association("Team").Append(&user1, &user2, &[]User{user1, user2, user3})
}
```

## 关联操作 - Preload / Joins 预加载

```go
func main() {
   db, _ := gorm.Open(mysql.Open("user:password@tcp(127.0.0.1:3306/hello"))
   var users []User
   var user User
   // 查询用户的时候并找出其订单，个人信息（1 + 1条SQL）
   db.Preload("Orders").Preload("Profile").Find(&users)
   // SELECT * FROM users
   // SELECT * FROM orders WHERE user_id IN (1,2,3,4); // 一对多
   // SELECT * FROM profiles WHERE user_id IN (1,2,3,4); // 一对一

   // 使用 Join SQL 加载 （单条JOIN SQL）
   db.Joins("Company").Joins("Manager").Find(&user, 1)
   db.Joins("Company", db.Where(&Company{Alive: true})).Find(&users)

   // 预加载关联全部 （只加载一级关联）
   db.Preload(clause.Associations).Find(&users)

   // 多级预加载
   db.Preload("Orders.OrderItems.Product").Find(&users)
   // 多级预加载 + 预加载全部一级关联
   db.Preload("Orders.OrderItems.Product").Preload(clause.Associations).Find(&users)

   // 查询用户的时候找出其未取消的订单
   db.Preload("Orders", "state NOT IN (?)", "cancelled").Find(&users)
   db.Preload("Orders", "state =  ?", "paid").Preload("Orders.OrderItems").Find(&users)

   db.Preload("Orders", func(db *gorm.DB) *gorm.DB {
      return db.Order("orders.amount DESC")
   }).Find(&users)
}
```

## 级联删除

```go
import (
   "gorm.io/driver/mysql"
   "gorm.io/gorm"
   "gorm.io/gorm/clause"
)

type Order struct {
}

type Account struct {
}

type CreditCard struct {
   BillingAddress interface{}
}

type User struct {
   ID          uint
   Name        string
   Account     Account      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
   CreditCards []CreditCard `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
   Orders      []Order      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func main() {
   db, _ := gorm.Open(mysql.Open("user:password@tcp(127.0.0.1:3306/hello"))
   var user User
   // 如果需要使用GORM Migrate 数据库迁移数据库外键才行
   db.AutoMigrate(&User{})
   // 入宫未启用软删除，在删除 User 时会自动删除其依赖
   db.Delete(&User{})

   // 方法2：使用 Select 实现级联删除，不依赖数据库的约束以及软删除
   // 删除 user 时候，也删除 user 的 Orders、CreditCards 记录
   db.Select("Account").Delete(&user)

   // 删除 user时候，也删除 user 的Orders、CreditCards 记录
   db.Select("Orders", "CreditCards").Delete(&user)

   // 删除 user 时候， 也删除 user 的 Orders、CreditCards 记录，也删除订单的BillingAddress
   db.Select("Orders", "Orders.BillingAddress", "CreditCards").Delete(&user)

   // 删除 user 时候，也删除用户及其依赖的所有 has one/many、many2many 记录
   db.Select(clause.Associations).Delete(&user)
}
```

## 总结

- 基本用法
- Model 定义
- 惯例约束
- 关联

# GORM设计原理

![image-20220516085457201](images/image-20220516085457201.png)

## SQL 生成

```sql
select `name`, `age`, `employee_number`
FROM `users`
where role <> "manager"
  AND age > 35
ORDER BY age DESC LIMIT 10
OFFSET 10 FOR
UPDATE
```

```go
db.Where("role <> ?", "manager").Where("age > ?", 35).
   Limit(100).Order("age desc").Find(&user)
```

![image-20220516090711740](images/image-20220516090711740.png)

> GORM API 方法添加 Clauses 至 GORM Statement

```go
// Where add conditions
func (db *DB) Where(query interface{}, args ...interface{}) (tx *DB) {
	tx = db.getInstance()
	if conds := tx.Statement.BuildCondition(query, args...); len(conds) > 0 {
		tx.Statement.AddClause(clause.Where{Exprs: conds})
	}
	return
}

// Where add conditions
func (db *DB) Where(query interface{}, args ...interface{}) (tx *DB) {
	tx = db.getInstance()
	if conds := tx.Statement.BuildCondition(query, args...); len(conds) > 0 {
		tx.Statement.AddClause(clause.Where{Exprs: conds})
	}
	return
}
```

> GORM Finisher 方法执行 GORM Statemnt

```go
// Find find records that match given conditions
func (db *DB) Find(dest interface{}, conds ...interface{}) (tx *DB) {
   tx = db.getInstance()
   if len(conds) > 0 {
      if exprs := tx.Statement.BuildCondition(conds[0], conds[1:]...); len(exprs) > 0 {
         tx.Statement.AddClause(clause.Where{Exprs: exprs})
      }
   }
   tx.Statement.Dest = dest
   return tx.callbacks.Query().Execute(tx)
}
```

```sql
func (p *processor) Execute(db *DB) *DB {
   // call scopes
   for len(db.Statement.scopes) > 0 {
      scopes := db.Statement.scopes
      db.Statement.scopes = nil
      for _, scope := range scopes {
         db = scope(db)
      }
   }

   var (
      curTime           = time.Now()
      stmt              = db.Statement
      resetBuildClauses bool
   )

   if len(stmt.BuildClauses) == 0 {
      stmt.BuildClauses = p.Clauses
      resetBuildClauses = true
   }

   // assign model values
   if stmt.Model == nil {
      stmt.Model = stmt.Dest
   } else if stmt.Dest == nil {
      stmt.Dest = stmt.Model
   }

   // parse model values
   if stmt.Model != nil {
      if err := stmt.Parse(stmt.Model); err != nil && (!errors.Is(err, schema.ErrUnsupportedDataType) || (stmt.Table == "" && stmt.TableExpr == nil && stmt.SQL.Len() == 0)) {
         if errors.Is(err, schema.ErrUnsupportedDataType) && stmt.Table == "" && stmt.TableExpr == nil {
            db.AddError(fmt.Errorf("%w: Table not set, please set it like: db.Model(&user) or db.Table(\"users\")", err))
         } else {
            db.AddError(err)
         }
      }
   }

   // assign stmt.ReflectValue
   if stmt.Dest != nil {
      stmt.ReflectValue = reflect.ValueOf(stmt.Dest)
      for stmt.ReflectValue.Kind() == reflect.Ptr {
         if stmt.ReflectValue.IsNil() && stmt.ReflectValue.CanAddr() {
            stmt.ReflectValue.Set(reflect.New(stmt.ReflectValue.Type().Elem()))
         }

         stmt.ReflectValue = stmt.ReflectValue.Elem()
      }
      if !stmt.ReflectValue.IsValid() {
         db.AddError(ErrInvalidValue)
      }
   }

   for _, f := range p.fns {
      f(db)
   }

   if stmt.SQL.Len() > 0 {
      db.Logger.Trace(stmt.Context, curTime, func() (string, int64) {
         return db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...), db.RowsAffected
      }, db.Error)
   }

   if !stmt.DB.DryRun {
      stmt.SQL.Reset()
      stmt.Vars = nil
   }

   if resetBuildClauses {
      stmt.BuildClauses = nil
   }

   return db
}
```

> 为什么这样设计？

- 自定义Clause Builder
- 方便扩展 Clause
- 自由选择 Clauses

### 自定义 Builder

```go
// 不同数据库甚至不同版本的数据库支持的 SQL 不同
// SELECT * FROM `users` LOCK IN SHARE MODE // MySQL < 8,MariaDB
// SELECT * FROM `users` FOR SHARE OF `users` // MySQL 8
db.Clauses(clause.Locking{
   Strength: "SHARE",
   Table:    clause.Table{Name: clause.CurrentTable},
}).Find(&users)
```

```go
func (dialector Dialector) Initialize(db *gorm.DB) (err error) {
   ctx := context.Background()

   // register callbacks
   callbacks.RegisterDefaultCallbacks(db, &callbacks.Config{
      CreateClauses: CreateClauses,
      QueryClauses:  QueryClauses,
      UpdateClauses: UpdateClauses,
      DeleteClauses: DeleteClauses,
   })

   if dialector.DriverName == "" {
      dialector.DriverName = "mysql"
   }

   if dialector.DefaultDatetimePrecision == nil {
      dialector.DefaultDatetimePrecision = &defaultDatetimePrecision
   }

   if dialector.Conn != nil {
      db.ConnPool = dialector.Conn
   } else {
      db.ConnPool, err = sql.Open(dialector.DriverName, dialector.DSN)
      if err != nil {
         return err
      }
   }

   if !dialector.Config.SkipInitializeWithVersion {
      err = db.ConnPool.QueryRowContext(ctx, "SELECT VERSION()").Scan(&dialector.ServerVersion)
      if err != nil {
         return err
      }

      if strings.Contains(dialector.ServerVersion, "MariaDB") {
         dialector.Config.DontSupportRenameIndex = true
         dialector.Config.DontSupportRenameColumn = true
         dialector.Config.DontSupportForShareClause = true
      } else if strings.HasPrefix(dialector.ServerVersion, "5.6.") {
         dialector.Config.DontSupportRenameIndex = true
         dialector.Config.DontSupportRenameColumn = true
         dialector.Config.DontSupportForShareClause = true
      } else if strings.HasPrefix(dialector.ServerVersion, "5.7.") {
         dialector.Config.DontSupportRenameColumn = true
         dialector.Config.DontSupportForShareClause = true
      } else if strings.HasPrefix(dialector.ServerVersion, "5.") {
         dialector.Config.DisableDatetimePrecision = true
         dialector.Config.DontSupportRenameIndex = true
         dialector.Config.DontSupportRenameColumn = true
         dialector.Config.DontSupportForShareClause = true
      }
   }

   for k, v := range dialector.ClauseBuilders() {
      db.ClauseBuilders[k] = v
   }
   return
}
```

### 扩展子句

![image-20220516093718541](images/image-20220516093718541.png)

![image-20220516093734169](images/image-20220516093734169.png)

## 插件扩展

### 插件是怎么工作的？

![image-20220516093949555](images/image-20220516093949555.png)

![image-20220516094116585](images/image-20220516094116585.png)

```go
import (
   "gorm.io/driver/mysql"
   "gorm.io/gorm"
)

func main() {
   db, _ := gorm.Open(mysql.Open("user:password@tcp(127.0.0.1:3306/hello"))

   // 注册新的 Callback
   db.Callback().Create().Register("MyPlugin", func(db *gorm.DB) {})

   // 删除 Callback
   db.Callback().Create().Remove("gorm:begin_transaction")

   // 替换 Callback
   db.Callback().Create().Replace("gorm:before_create", func(db *gorm.DB) {})

   // 查询注册的 Callback
   db.Callback().Create().Get("gorm:create")

   // 指定 Callback 顺序
   db.Callback().Create().
      Before("gorm:create").
      After("MyPlugin").
      Register("MyPlugin2", func(db *gorm.DB) {})

   // 注册到所有服务之前
   db.Callback().Create().Before("*").Register("MyPlugin:newCallBack", func(db *gorm.DB) {})

   // 注册时检查条件
   enableTransaction := func(db *gorm.DB) bool { return !db.SkipDefaultTransaction }
   db.Callback().Create().Match(enableTransaction).Register("gorm:begin_transaction",  func(db *gorm.DB) {})
}
```

![image-20220516124059462](images/image-20220516124059462.png)

### 多租户

```go
import (
   "context"
   "gorm.io/driver/mysql"
   "gorm.io/gorm"
)

func getTenantID(ctx context.Context) (uint, error) {
   return 0, nil
}

func main() {
   db, _ := gorm.Open(mysql.Open("user:password@tcp(127.0.0.1:3306/hello"))
   // 根据 TenantID 过滤
   var setTenantScope = func(db *gorm.DB) {
      if tenantID, err := getTenantID(db.Statement.Context); err != nil {
         db.Where("tenant_id = ?", tenantID)
      } else {
         db.AddError(err)
      }
   }

   db.Callback().Query().Before("gorm:query").Register("set_tenant_scope", setTenantScope)
   db.Callback().Delete().Before("gorm:delete").Register("set_tenant_scope", setTenantScope)
   db.Callback().Update().Before("gorm:update").Register("set_tenant_scope", setTenantScope)

   // 设置 TenantID
   var setTenantID = func(db *gorm.DB) {
      tenantID, _ := getTenantID(db.Statement.Context)
      db.Statement.SetColumn("tenant_id", tenantID)
      // ...
   }

   db.Callback().Update().Before("gorm:create").Register("set_tenant_id", setTenantID)
}
```

### 多数据库、读写分离

![image-20220516143730069](images/image-20220516143730069.png)

## ConnPool

![image-20220516143954766](images/image-20220516143954766.png)

```go
import (
   "gorm.io/driver/sqlite"
   "gorm.io/gorm"
)

type User struct {
}

type Session struct {
   PrepareStmt bool
}

func main() {
   // 全局模式，所有 DB 操作都会预编译并缓存（缓存不包含参数部分）
   db, _ := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{PrepareStmt: true})
   var user User
   var users []User
   db.First(&user)

   // 会话模式，后续会话的操作都会预售并缓存
   tx := db.Session(&gorm.Session{PrepareStmt: true})
   tx.Find(&users)
   tx.Model(&user).Update("Age", 18)
   // 全局缓存的语句可能会被会话使用
   tx.First(&user, 2)

   stmtManager, _ := tx.ConnPool.(*gorm.PreparedStmtDB)
   // 关闭当前会话的预编译语句
   stmtManager.Close()
}
```

![image-20220516145439461](images/image-20220516145439461.png)

通过自定义插件，先把数据放在某国数据库，再把数据放在原数据库，相对于业务层透明

![image-20220516145604277](images/image-20220516145604277.png)

![image-20220516145815777](images/image-20220516145815777.png)

### 遇到的问题

最开始设计，使用预编译的sql，速度不但没有提升还下降了

原因：预编译的sql使用完就丢弃了，并没有缓存起来

![image-20220516150020205](images/image-20220516150020205.png)

```go
import (
   "code.byted.org/gorm/bytedgorm"
   "gorm.io/gorm"
)

func main() {
   DB, err := gorm.Open(
      // psm 的格式为 p.s.m 无需 _write. _read 等后缀，dbname 为数据库名字
      bytedgorm.MySql("p.s.m" /* 数据库 PSM */, "dbname" /*数据库名*/).WithReadReplicas(),
      bytedgorm.WithDefaults(),
   )
}
```

## Dialector

### Dialector是什么？

```go
import (
   "gorm.io/driver/clickhouse"
   "gorm.io/driver/mysql"
   "gorm.io/driver/postgres"
   "gorm.io/gorm"
)

func main() {
   dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=true&loc=Local"
   _, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
   _, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{})
   _, _ = gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
   // import "xxx.io/caches"
   //_, _ = gorm.Open(caches.New(caches.Config{
   // Fallback: mysql.Open(dsn),
   // Store:    lru.New(lru.Config{}),
   //}))
}
```

![image-20220516152722780](images/image-20220516152722780.png)

# GORM最佳实践

## 数据序列化与 SQL 表达式

> SQL表达式更新创建

![image-20220516153811261](images/image-20220516153811261.png)

> SQL 表达式查询

![image-20220516153835948](images/image-20220516153835948.png)

> 数据序列化

![image-20220516153903778](images/image-20220516153903778.png)

## 批量数据操作

> 批量创建 / 查询

![image-20220516153951169](images/image-20220516153951169.png)

> 批量更新

![image-20220516154131861](images/image-20220516154131861.png)

> 批量数据加速操作

![image-20220516154615131](images/image-20220516154615131.png)

## 代码复用、分库分表、Sharding

> 代码复用

![image-20220516154650748](images/image-20220516154650748.png)

> 分库分表

![image-20220516154725452](images/image-20220516154725452.png)

## Sharding

![image-20220516154801682](images/image-20220516154801682.png)

## 混沌工程 / 压测

![image-20220516154833229](images/image-20220516154833229.png)

![image-20220516154851141](images/image-20220516154851141.png)

## Logger / Trace

![image-20220516154911561](images/image-20220516154911561.png)

## Migrator

![image-20220516154931195](images/image-20220516154931195.png)

![image-20220516154948513](images/image-20220516154948513.png)

## Gen 代码生成 / Raw SQL

> Row SQL

![image-20220516155007761](images/image-20220516155007761.png)

> Gen

![image-20220516155029474](images/image-20220516155029474.png)

## 安全

![image-20220516155119028](images/image-20220516155119028.png)
