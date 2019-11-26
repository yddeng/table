package pgsql

import (
	"fmt"
	"strings"
)

/*
 * 插入数据
 * tableName:表名 fields:键值对
 */
func Set(tableName string, fields map[string]interface{}) error {
	sqlStr := `
INSERT INTO "%s" (%s)
VALUES (%s);`

	keys, values := []string{}, []string{}
	args := []interface{}{}
	var i = 1
	for k, v := range fields {
		keys = append(keys, k)
		values = append(values, fmt.Sprintf("$%d", i))
		i++
		args = append(args, v)
	}

	sqlStatement := fmt.Sprintf(sqlStr, tableName, strings.Join(keys, ","), strings.Join(values, ","))
	smt, err := dbConn.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	_, err = smt.Exec(args...)
	return err
}

/*
 * 更新数据
 * tableName:表名 whereStr:选择规则 fields:键值对
 */
func Update(tableName, whereStr string, fields map[string]interface{}) error {
	sqlStr := `
UPDATE "%s" 
SET %s
WHERE %s;`

	keys := []string{}
	args := []interface{}{}
	var i = 1
	for k, v := range fields {
		keys = append(keys, fmt.Sprintf(`%s = $%d`, k, i))
		i++
		args = append(args, v)
	}

	sqlStatement := fmt.Sprintf(sqlStr, tableName, strings.Join(keys, ","), whereStr)
	//fmt.Println(sqlStatement)
	smt, err := dbConn.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	_, err = smt.Exec(args...)
	return err

}

/*
 * 没有数据插入，有则添加。
 * tableName:表名 key:主键名 fields:键值对，包含主键
 */
func SetNx(tableName, key string, fields map[string]interface{}) error {
	sqlStr := `
INSERT INTO "%s" (%s)
VALUES(%s) 
ON conflict(%s) DO 
UPDATE SET %s;`

	keys, values, sets := []string{}, []string{}, []string{}
	args := []interface{}{}
	var i = 1
	for k, v := range fields {
		keys = append(keys, k)
		values = append(values, fmt.Sprintf("$%d", i))
		if key != k {
			sets = append(sets, fmt.Sprintf(`%s = $%d`, k, i))
		}
		i++
		args = append(args, v)
	}

	sqlStatement := fmt.Sprintf(sqlStr, tableName, strings.Join(keys, ","), strings.Join(values, ","), key, strings.Join(sets, ","))
	//fmt.Println(sqlStatement)
	smt, err := dbConn.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	_, err = smt.Exec(args...)
	return err
}

/*
 * 读取数据。
 * tableName:表名 whereStr:选择规则 fields:要查询的键名
 * ret 返回键值对
 */
func Get(tableName, whereStr string, fields []string) (ret map[string]interface{}, err error) {
	sqlStr := `
SELECT %s FROM "%s" 
WHERE %s;`

	keys := []string{}
	values := []interface{}{}
	for _, k := range fields {
		keys = append(keys, k)
		values = append(values, new(interface{}))
	}

	sqlStatement := fmt.Sprintf(sqlStr, strings.Join(keys, ","), tableName, whereStr)
	//fmt.Println(sqlStatement)
	row := dbConn.QueryRow(sqlStatement)
	err = row.Scan(values...)
	if err != nil {
		return nil, err
	}

	ret = map[string]interface{}{}
	for i, k := range fields {
		ret[k] = *(values[i].(*interface{}))
	}

	return ret, nil
}

/*
 * 读取所有数据。
 * tableName:表名 fields:要查询的键名
 * ret 返回键值对的slice
 */
func GetAll(tableName string, fields []string) (ret []map[string]interface{}, err error) {
	keys := []string{}
	values := []interface{}{}
	for _, k := range fields {
		keys = append(keys, k)
		values = append(values, new(interface{}))
	}

	sqlStatement := fmt.Sprintf(`SELECT %s FROM "%s";`, strings.Join(keys, ","), tableName)
	rows, err := dbConn.Query(sqlStatement)
	if err != nil {
		return nil, err
	}

	ret = []map[string]interface{}{}
	for rows.Next() {
		err := rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		mid := map[string]interface{}{}
		for i, k := range fields {
			mid[k] = *(values[i].(*interface{}))
		}
		ret = append(ret, mid)
	}

	return ret, nil
}
