package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"
	"todo/model"

	"github.com/stretchr/testify/assert"
)

func TestTodos(t *testing.T) {
	os.Remove("./test.db")

	assert := assert.New(t)
	ah := MakeHandler("./test.db")
	defer ah.Close()

	ts := httptest.NewServer(ah)
	defer ts.Close()

	res, err := http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test todo"}})
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	var todo model.Todo
	err = json.NewDecoder(res.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name, "Test todo")

	id1 := todo.ID

	res2, err := http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test todo2"}})
	assert.NoError(err)
	assert.Equal(http.StatusOK, res2.StatusCode)

	var todo2 model.Todo
	err = json.NewDecoder(res2.Body).Decode(&todo2)
	assert.NoError(err)

	id2 := todo2.ID

	resp, err := http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	todos := []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)

	for _, t := range todos {
		if t.ID == id1 {
			assert.Equal("Test todo", t.Name)
		} else if t.ID == id2 {
			assert.Equal("Test todo2", t.Name)
		} else {
			assert.Error(fmt.Errorf("testID should be id1 or id2"))
		}
	}

	resp, err = http.Get(ts.URL + "/complete-todo/" + strconv.Itoa(id1) + "?complete=true")
	assert.NoError(err)
	assert.Equal(resp.StatusCode, http.StatusOK)

	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	todos = []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)

	for _, t := range todos {
		if t.ID == id1 {
			assert.True(t.Completed)
		}
	}

	req, _ := http.NewRequest("DELETE", ts.URL+"/todos/"+strconv.Itoa(id1), nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(resp.StatusCode, http.StatusOK)

	resp, err = http.Get(ts.URL + "/complete-todo/" + strconv.Itoa(id1) + "?complete=true")
	assert.NoError(err)
	assert.Equal(resp.StatusCode, http.StatusOK)

	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	todos = []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 1)

	for _, t := range todos {
		assert.Equal(t.ID, id2)
	}
}
