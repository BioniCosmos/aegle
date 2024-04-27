package graphql

import (
	"log"
	"reflect"

	"github.com/bionicosmos/submgr/models"
	"github.com/graphql-go/graphql"
)

var schema graphql.Schema

func init() {
	var err error
	schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootQuery",
			Fields: graphql.Fields{
				"user": &graphql.Field{
					Type: graphql.NewObject(graphql.ObjectConfig{
						Name: "User",
						Fields: graphql.Fields{
							"id":        &graphql.Field{Type: graphql.ID},
							"name":      &graphql.Field{Type: graphql.String},
							"email":     &graphql.Field{Type: graphql.String},
							"level":     &graphql.Field{Type: graphql.Int},
							"startDate": &graphql.Field{Type: graphql.String},
							"cycles":    &graphql.Field{Type: graphql.Int},
							"account": &graphql.Field{Type: graphql.NewObject(graphql.ObjectConfig{
								Name: "Account",
								Fields: graphql.Fields{
									"vless":  &graphql.Field{Type: graphql.String},
									"vmess":  &graphql.Field{Type: graphql.String},
									"trojan": &graphql.Field{Type: graphql.String},
								},
							})},
							"profileIds": &graphql.Field{Type: graphql.NewList(graphql.ID)},
						},
					}),
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{Type: graphql.String},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						id := p.Args["id"].(string)
						user, err := models.FindUser(id)
						if err != nil {
							return nil, err
						}
						userType := reflect.TypeOf(*user)
						userValue := reflect.ValueOf(*user)
						data := make(map[string]any)
						for i := 0; i < userType.NumField(); i++ {
							field := userType.Field(i)
							fieldValue := userValue.Field(i)
							data[field.Tag.Get("json")] = fieldValue.Interface()
						}
						delete(data, "profiles")
						data["profileIds"] = make([]string, 0)
						profileIds := data["profileIds"].([]string)
						for profileId := range user.Profiles {
							profileIds = append(profileIds, profileId.Hex())
						}
						return data, nil
					},
				},
			},
		}),
	})
	if err != nil {
		log.Fatal(err)
	}
}
