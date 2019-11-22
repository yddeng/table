package pgsql

import (
	"fmt"
	"strings"
)

//
func InsertUser(uerName, pwd, per string) error {
	sqlStr := `
INSERT INTO "user" (user_name, password,permission,)
VALUES ('%s','%s','%s');`

	sqlStatement := fmt.Sprintf(sqlStr, uerName, pwd, per)
	smt, err := dbConn.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	_, err = smt.Exec()
	return err
}

func UpdateUser(key string, fields map[string]interface{}) error {
	sqlStr := `
UPDATE "user" 
SET %s
WHERE user_name = '%s';`

	keys := []string{}
	args := []interface{}{}
	var i = 1
	for k, v := range fields {
		keys = append(keys, fmt.Sprintf(`%s = $%d`, k, i))
		i++
		args = append(args, v)
	}

	sqlStatement := fmt.Sprintf(sqlStr, strings.Join(keys, ","), key)
	smt, err := dbConn.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	_, err = smt.Exec(args...)
	return err
}

//
func LoadUser(userName string) (string, string, error) {
	sqlStatement := `
SELECT password,permission FROM "user" 
WHERE user_name = $1;`

	var pwd, per string
	row := dbConn.QueryRow(sqlStatement, userName)
	err := row.Scan(&pwd, &per)
	if err != nil {
		return "", "", err
	}
	return pwd, per, nil
}

func AllUser() ([]map[string]interface{}, error) {
	rows, err := dbConn.Query(`SELECT * FROM "user";`)
	if err != nil {
		return nil, err
	}

	ret := []map[string]interface{}{}
	for rows.Next() {
		var name, pwd, permission string
		err := rows.Scan(&name, &pwd, &permission)
		if err != nil {
			return nil, err
		}
		ret = append(ret, map[string]interface{}{
			"userName":   name,
			"password":   pwd,
			"permission": permission,
		})
	}

	return ret, nil
}
