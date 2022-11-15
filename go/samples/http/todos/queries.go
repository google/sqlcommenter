package todos

type PGQueries struct{}
type MySQLQueries struct{}

type TodosQueries interface {
	CreateTodosTableIfNotExists() string
	ListTodos() string
	InsertTodo() string
	UpdateTodo() string
	DeleteTodo() string
	TodoById() string
}

// PGQueries impl
func (sql PGQueries) CreateTodosTableIfNotExists() string {
	return `CREATE TABLE IF NOT EXISTS todos(
		id SERIAL NOT NULL,
		task VARCHAR(1000) NOT NULL,

		PRIMARY KEY(id)
)`
}

func (sql PGQueries) ListTodos() string {
	return "SELECT id, task FROM todos WHERE LOWER(task) LIKE $1 ORDER BY id DESC"
}

func (sql PGQueries) InsertTodo() string {
	return "INSERT INTO todos (task) VALUES ($1) RETURNING id, task"
}

func (sql PGQueries) UpdateTodo() string {
	return "UPDATE todos SET task=$1 WHERE id=$2"
}

func (sql PGQueries) DeleteTodo() string {
	return "DELETE FROM todos WHERE id=$1"
}

func (sql PGQueries) TodoById() string {
	return "SELECT id, task FROM todos WHERE id=$1"
}

// MySQLQueries impl
func (sql MySQLQueries) CreateTodosTableIfNotExists() string {
	return `CREATE TABLE IF NOT EXISTS todos(
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		task VARCHAR(1000) NOT NULL
)`
}

func (sql MySQLQueries) ListTodos() string {
	return "SELECT id, task FROM todos WHERE task LIKE ? ORDER BY id DESC"
}

func (sql MySQLQueries) InsertTodo() string {
	return "INSERT INTO todos (task) VALUES (?)"
}

func (sql MySQLQueries) UpdateTodo() string {
	return "UPDATE todos SET task=? WHERE id=?"
}

func (sql MySQLQueries) DeleteTodo() string {
	return "DELETE FROM todos WHERE id=?"
}

func (sql MySQLQueries) TodoById() string {
	return "SELECT id, task FROM todos WHERE id=?"
}
