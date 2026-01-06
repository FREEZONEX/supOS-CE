package dbutil

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// EscapeForLike escapes special characters for SQL LIKE queries.
func EscapeForLike(keyword string) string {
	if keyword == "" {
		return keyword
	}
	keyword = strings.ReplaceAll(keyword, `\`, `\\`)
	keyword = strings.ReplaceAll(keyword, "%", `\%`)
	keyword = strings.ReplaceAll(keyword, "_", `\_`)
	keyword = strings.ReplaceAll(keyword, "'", "''")
	return keyword
}

// BuildLikePattern builds a LIKE pattern with wildcards.
func BuildLikePattern(keyword string) string {
	if keyword == "" {
		return "%"
	}
	return "%" + EscapeForLike(keyword) + "%"
}

// GetFullTableName ensures the table name is properly quoted.
func GetFullTableName(tableName string) string {
	var sb strings.Builder
	GetFullTableNameBuilder(tableName, &sb)
	return sb.String()
}

// GetFullTableNameBuilder is a builder version of GetFullTableName.
func GetFullTableNameBuilder(tableName string, sb *strings.Builder) {
	if len(tableName) > 1 && tableName[0] == '"' && tableName[len(tableName)-1] == '"' {
		sb.WriteString(tableName)
		return
	}
	dot := strings.Index(tableName, ".")
	if dot > 0 {
		db := tableName[:dot]
		table := tableName[dot+1:]
		addQuotation(sb, db)
		sb.WriteByte('.')
		addQuotation(sb, table)
	} else {
		addQuotation(sb, tableName)
	}
}

// GetCleanTableName removes quotes and schema from a table name.
func GetCleanTableName(tableName string) string {
	st := strings.LastIndex(tableName, ".") + 1
	ed := len(tableName)
	if st < ed && tableName[st] == '"' {
		st++
	}
	if st < ed && tableName[ed-1] == '"' {
		ed--
	}
	return tableName[st:ed]
}

func addQuotation(sb *strings.Builder, name string) {
	if name[0] != '"' {
		sb.WriteByte('"')
	}
	sb.WriteString(name)
	if name[len(name)-1] != '"' {
		sb.WriteByte('"')
	}
}

// ExecuteSQLScript executes an SQL script from a file within a single transaction.
// The script is split into individual statements, correctly handling semicolons
// within strings and comments. If any statement fails, the entire transaction is rolled back.
func ExecuteSQLScript(db *sql.DB, scriptPath string) error {
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("script file not found: %s", scriptPath)
	}

	content, err := ioutil.ReadFile(scriptPath)
	if err != nil {
		return fmt.Errorf("failed to read script file %s: %w", scriptPath, err)
	}

	statements, err := splitSQLStatements(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse script file %s: %w", scriptPath, err)
	}

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Execute each statement within the transaction
	for _, stmt := range statements {
		if _, err := tx.Exec(stmt); err != nil {
			// If an error occurs, roll back the transaction and return the error
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("transaction rollback failed: %v after execution error: %w", rbErr, err)
			}
			return fmt.Errorf("failed to execute statement: %s\nError: %w", stmt, err)
		}
	}

	// If all statements were successful, commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// splitSQLStatements intelligently splits a SQL script into individual statements.
// It handles semicolons inside strings and comments.
func splitSQLStatements(script string) ([]string, error) {
	var statements []string
	var currentStatement strings.Builder
	var inString, inLineComment, inBlockComment bool
	var stringDelimiter rune

	runes := []rune(script)
	for i := 0; i < len(runes); i++ {
		char := runes[i]

		if inString {
			if char == stringDelimiter {
				if i+1 < len(runes) && runes[i+1] == stringDelimiter { // Handle escaped quotes (e.g., '')
					currentStatement.WriteRune(char)
					currentStatement.WriteRune(runes[i+1])
					i++
					continue
				}
				inString = false
			}
		} else if inLineComment {
			if char == '\n' {
				inLineComment = false
			}
		} else if inBlockComment {
			if char == '*' && i+1 < len(runes) && runes[i+1] == '/' {
				inBlockComment = false
				i++ // Skip the closing '/'
				continue
			}
		} else {
			switch char {
			case '\'', '"':
				inString = true
				stringDelimiter = char
			case '-':
				if i+1 < len(runes) && runes[i+1] == '-' {
					inLineComment = true
					i++ // Skip the second '-'
					continue
				}
			case '/':
				if i+1 < len(runes) && runes[i+1] == '*' {
					inBlockComment = true
					i++ // Skip the '*'
					continue
				}
			case ';':
				stmt := strings.TrimSpace(currentStatement.String())
				if stmt != "" {
					statements = append(statements, stmt)
				}
				currentStatement.Reset()
				continue
			}
		}
		currentStatement.WriteRune(char)
	}

	// Add the last statement if it's not empty
	lastStmt := strings.TrimSpace(currentStatement.String())
	if lastStmt != "" {
		statements = append(statements, lastStmt)
	}

	return statements, nil
}
