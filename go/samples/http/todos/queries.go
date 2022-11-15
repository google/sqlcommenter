package todos

type PGQueries struct{}
type MySQLQueries struct{}

type TodosQueries interface {
	CreateTodosTableIfNotExists() string
	ListTodos(search string) string
	InsertTodo(task string) string
	UpdateTodo(id int, task string) string
	DeleteTodo(id int) string
}

// PGQueries impl
func (sql PGQueries) CreateTodosTableIfNotExists() string {
	return ""
}

func (sql PGQueries) ListTodos(search string) string {
	return ""
}

func (sql PGQueries) InsertTodo(task string) string {
	return ""
}

func (sql PGQueries) UpdateTodo(id int, task string) string {
	return ""
}

func (sql PGQueries) DeleteTodo(id int) string {
	return ""
}

// MySQLQueries impl
func (sql MySQLQueries) CreateTodosTableIfNotExists() string {
	return ""
}

func (sql MySQLQueries) ListTodos(search string) string {
	return ""
}

func (sql MySQLQueries) InsertTodo(task string) string {
	return ""
}

func (sql MySQLQueries) UpdateTodo(id int, task string) string {
	return ""
}

func (sql MySQLQueries) DeleteTodo(id int) string {
	return ""
}
