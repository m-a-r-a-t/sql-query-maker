package main

import (
	"fmt"
	sql_query_maker "github.com/m-a-r-a-t/sql-query-maker"
)

type Object struct {
	Id   string
	Name string
}

func main() {
	maker := sql_query_maker.NewQueryMaker(10)

	maker.Add("INSERT INTO object(id,name)")

	values := []Object{
		{"1", "ffd"},
		{"2", "aaa"},
		{"2", "aaa"},
		{"2", "aaa"},
		{"2", "aaa"},
		{"2", "aaa"},
		{"2", "aaa"},
	}

	for i := 0; i < len(values); i++ {
		maker.Values(values[i].Id, values[i].Name)
	}

	query, args := maker.Make()

	fmt.Println(query)
	fmt.Println(args)

	maker.Clear()

	maker.Add("UPDATE project SET")

	maker.Add("name = ?,", "ezhe")
	maker.Add("email = ?,", "abc@mail.ru")

	maker.Where("id = ?", "12345")

	query, args = maker.Make()

	fmt.Println(query)
	fmt.Println(args)

}
