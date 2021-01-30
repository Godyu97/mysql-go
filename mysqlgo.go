package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

type user struct {
	id   int
	name string
	age  int
}

func dbInit() error {
	dsn := "root:yu970928@tcp(104.233.130.15:3306)/godyu97centos?charset=utf8mb4"
	//检验dsn格式并尝试连接，无法确认是否连接成功
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	//检验dsn数据源是否正确，是否连接成功
	err = db.Ping()
	if err != nil {
		return err
	}

	//设置数据库最大连接数
	db.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	db.SetMaxIdleConns(10)

	return nil
}
func main() {
	err = dbInit()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("数据库连接成功！")

	// useDatabase("godyu97centos")

	// creatTestDB("testDB2")

	// fmt.Print("请输入查询ID:")
	// var n1 int
	// fmt.Scanf("%d\n", &n1)
	// queryRow(n1)

	// fmt.Print("请输入查询ID范围>=?:")
	// var n2 int
	// fmt.Scanf("%d\n", &n2)
	// queryRows(n2)

	// insertRow()

	// updateRow()

	// deleteROW()

	// prepareDelete()

	// transactionAge()

}
func queryRow(n int) {
	sqlStr := `select user_ID,user_name,user_age from usert where user_ID=?;`
	var u user

	//QueryROw和Scan成对使用；sql.Row的Scan方法中内置了defer Close语句
	err = db.QueryRow(sqlStr, n).Scan(&u.id, &u.name, &u.age)
	if err != nil {
		fmt.Println("查询Row失败")
		return
	}
	fmt.Println(u)
}

func queryRows(n int) {
	sqlStr := `select user_ID,user_name,user_age from usert where user_ID>=?;`
	rows, err := db.Query(sqlStr, n)
	if err != nil {
		fmt.Println("查询Rows失败")
		return
	}
	//sql.Rows需要手动Close
	defer rows.Close()
	for rows.Next() {
		var u user
		err = rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			fmt.Println("Rows Scan失败")
			return
		}
		fmt.Println(u)
	}
}

func insertRow() {

	fmt.Println("请输入插入的数据")
	var u user
	fmt.Print("ID:")
	fmt.Scanf("%d\n", &u.id)
	fmt.Print("name:")
	fmt.Scanf("%s\n", &u.name)
	fmt.Print("age:")
	fmt.Scanf("%d\n", &u.age)
	//mysql insert操作
	sqlStr := `insert into usert(user_ID,user_name,user_age) values (?,?,?);`
	ret, err := db.Exec(sqlStr, u.id, u.name, u.age)
	if err != nil {
		fmt.Println("插入失败")
		return
	}
	theID, err := ret.LastInsertId()
	if err != nil {
		fmt.Println("获取新ID失败")
	}
	theRA, err := ret.RowsAffected()
	if err != nil {
		fmt.Println("获取受影响行失败")
	}
	fmt.Printf("插入成功,新ID=%d,受影响行=%d\n", theID, theRA)
}

func updateRow() {

	fmt.Println("请输入要修改的数据")
	fmt.Print("ID:")
	var id int
	fmt.Scanf("%d\n", &id)
	fmt.Print("修改age:")
	var age int
	fmt.Scanf("%d\n", &age)

	//mysql update操作
	sqlStr := `UPDATE usert set user_age=? where user_ID=?;`
	ret, err := db.Exec(sqlStr, age, id)
	if err != nil {
		fmt.Println("修改数据失败")
		return
	}
	theRA, err := ret.RowsAffected()
	if err != nil {
		fmt.Println("获取RA失败")
		return
	}
	fmt.Printf("修改数据ID:%d成功，%d行受影响\n", id, theRA)

}

func deleteROW() {
	fmt.Print("请输入删除数据的ID:")
	var id int
	fmt.Scanf("%d\n", &id)
	//mysql delete操作
	sqlStr := `DELETE from usert where user_ID=?;`
	ret, err := db.Exec(sqlStr, id)
	if err != nil {
		fmt.Println("删除数据失败")
		return
	}
	theRA, err := ret.RowsAffected()
	if err != nil {
		fmt.Println("获取RA失败")
		return
	}
	fmt.Printf("删除数据ID:%d成功，%d行受影响\n", id, theRA)
}

func prepareDelete() {
	fmt.Print("请输入删除数据的ID:")
	var id int
	fmt.Scanf("%d\n", &id)
	//mysql 预处理操作
	sqlStr := `DELETE from usert where user_ID=?;`
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Println("预处理失败")
		return
	}
	defer stmt.Close()
	ret, err := stmt.Exec(id)
	if err != nil {
		fmt.Println("删除数据失败")
		return
	}
	theRA, err := ret.RowsAffected()
	if err != nil {
		fmt.Println("获取RA失败")
		return
	}
	fmt.Printf("删除数据ID:%d成功，%d行受影响\n", id, theRA)
}

func creatTestDB(newDB string) {
	sqlStr := fmt.Sprintf("CREATE DATABASE %s;", newDB)
	_, err := db.Exec(sqlStr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("创建数据库%s成功\n", newDB)
}

func useDatabase(dbName string) {
	if db != nil {
		db.Close()
	}
	dsn := "root:yu970928@tcp(104.233.130.15:3306)/" + dbName
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	//设置数据库最大连接数
	db.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	db.SetMaxIdleConns(10)
	fmt.Printf("use 数据库:%s成功\n", dbName)
}

func transactionAge() {
	tx, err := db.Begin()
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		fmt.Println("事务开始失败，已回滚")
		return
	}

	sqlStr1 := `UPDATE usert SET user_age=user_age-1 WHERE user_ID=?;`
	ret1, err := tx.Exec(sqlStr1, 0)
	if err != nil {
		tx.Rollback()
		fmt.Println("exec sqlStr1 失败，已回滚")
		return
	}
	RA1, err := ret1.RowsAffected()
	if err != nil {
		tx.Rollback()
		fmt.Println("获取RA1失败")
		return
	}

	sqlStr2 := `UPDATE usert SET user_age=user_age+1 WHERE user_ID=?;`
	ret2, err := tx.Exec(sqlStr2, 1)
	if err != nil {
		tx.Rollback()
		fmt.Println("exec sqlStr2 失败，已回滚")
		return
	}
	RA2, err := ret2.RowsAffected()
	if err != nil {
		tx.Rollback()
		fmt.Println("获取RA2失败")
		return
	}

	if RA1 == 1 && RA2 == 1 {
		tx.Commit()
		fmt.Println("事务提交")
	} else {
		tx.Rollback()
		fmt.Println("事务回滚")
	}
}
