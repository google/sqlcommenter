package todos

import (
	"net/http"

	gosql "github.com/google/sqlcommenter/go/database/sql"
	"github.com/julienschmidt/httprouter"
)

type TodosController struct {
	DB  *gosql.DB
	SQL TodosQueries
}

func (c *TodosController) CreateTodosTableIfNotExists() error {
	return nil
}

func (c *TodosController) ActionList(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func (c *TodosController) ActionInsert(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func (c *TodosController) ActionUpdate(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func (c *TodosController) ActionDelete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}
