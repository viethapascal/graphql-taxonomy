package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"graphql-taxonomy/dao"
	"log"
	"strconv"
)

type User struct {
	Name string `json:"name"`
	ID string `json:"id"`
}

var UserType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "User",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.String,

				},
				"name": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

var QueryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"user": &graphql.Field{
					Type: UserType,
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						idQuery, isOk := p.Args["id"].(string)
						res, err := FindUser(idQuery)
						//log.Println("query result:", res)
						if isOk && err !=nil {
							return res, nil
						}
						return nil, err
					},
				},
				"list_users": &graphql.Field{
					Type: graphql.NewList(UserType),
					Description: "Get User List",
					Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
						res, err := ListUsers()
						return res, err
					},
				},
			},
		},
	)

var mutationType = graphql.NewObject(graphql.ObjectConfig{

})
var schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query: QueryType,
		},
	)

func FindUser(id string) (User, error){
	p, _ := strconv.Atoi(id)
	d, err := dao.Db.Query("select id, name from students where id=$1", p)
	if err != nil {
		panic(err)
	}
	res := User{}
	if d.Next() {
		err := d.Scan(&res.ID, &res.Name)
		if err != nil {
			return res, err
		}
	}
	return res, nil
}
func ListUsers() ([]User, error) {
	d, err := dao.Db.Query("select id, name from students")
	if err != nil {
		panic(err)
	}
	res := []User{}
	for d.Next() {
		var user = User{}
		err := d.Scan(&user.ID, &user.Name)
		if err != nil {
			return res, err
		}

		res = append(res, user)
	}
	return res, nil
}
func executeQuery(query string, schema graphql.Schema) (*graphql.Result, error) {
	result := graphql.Do(graphql.Params{
		Schema: schema,
		RequestString: query,
	})
	if result.HasErrors() {
		log.Println(result.Errors[0].Message)
		return nil, errors.New(result.Errors[0].Message)
	}
	return  result, nil
}
func GetUserHandler(c *gin.Context) {
	var (
		query = c.Query("query")
	)

	result, err := executeQuery(query, schema)
	if err == nil {
		c.JSON(200, result)
	} else {

		c.JSON(500, "Internal Error")
	}
}
