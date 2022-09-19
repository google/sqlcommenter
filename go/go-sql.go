package sqlcommenterGo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
	commenterOptions map[string]bool
}

var goSQLCommenterTags = map[string]string{"dbDriver": "go/sql"}

func Open(driverName, dataSourceName string, commenterOptions map[string]bool) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	return &DB{DB: db, commenterOptions: commenterOptions}, err
}

func (db *DB) QueryRow(query string, args ...interface{}) any {
	return db.DB.QueryRow(db.withComment(context.Background(), query), args...)
}

func (db *DB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.DB.QueryContext(ctx, db.withComment(ctx, query), args...)
}

func (db *DB) withComment(ctx context.Context, query string) string {

	var finalCommentsMap = map[string]string{}
	var finalCommentsStr string = ""
	query = strings.TrimSpace(query)

	for key, element := range db.commenterOptions {
		if element { // Checks if option is set as true
			if _, ok := goSQLCommenterTags[key]; ok { // Check if static value is assigned and if true append
				finalCommentsMap[key] = goSQLCommenterTags[key]
			} else if ctx.Value(key) != nil { // Append if key is avaliable in context
				finalCommentsMap[key] = ctx.Value(key).(string)
			}
		}

	}
	if len(finalCommentsMap) > 0 { // Converts comments map to string and appends it to query
		jsonStr, err := json.Marshal(finalCommentsMap)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		} else {
			finalCommentsStr = strings.Replace(string(jsonStr), "{", "/*", 1)
			finalCommentsStr = strings.Replace(string(finalCommentsStr), "}", "*/", 1)
		}
	}

	if strings.Contains(query, ";") {
		query = strings.Replace(string(query), ";", "", 1)
		return fmt.Sprintf("%s%s;", query, finalCommentsStr)
	}
	return fmt.Sprintf("%s%s", query, finalCommentsStr)

}

// TODO:
func encodeValue(v string) string {
	urlEscape := strings.ReplaceAll(url.PathEscape(string(v)), "+", "%20")
	return fmt.Sprintf("'%s'", urlEscape)
}

func encodeKey(k string) string {
	return url.QueryEscape(string(k))
}
