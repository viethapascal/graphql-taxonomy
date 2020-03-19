package graph_type

import (
	"errors"
	"fmt"
	"github.com/graphql-go/graphql"
	"graphql-taxonomy/dao"
	"log"
	"reflect"
)

type FacebookAudienceSegment struct {
	ID                    string                 `json:"id" multipart:"fbSegmentId"`
	InstanceID            string                 `json:"instance_id"`
	AdvertiserID          string                 `json:"advertiser_id"`
	FbAccountID           string                 `json:"fb_account_id" multipart:"accountId"`
	Type                  string                 `json:"type" multipart:"audience_type"`
	EventSourceID         string                 `json:"event_source_group"`
	CreatedTime           int64                  `json:"created_time,omitempty"`
	ModifiedTime          int64                  `json:"modified_time,omitempty"`
	DmpAudienceSegmentRef string                 `json:"dmp_audience_segment_ref,omitempty"`
	ServiceSegment        string                 `json:"service_segment,omitempty"`
	Archived              bool                   `json:"archived,omitempty"`
	AudienceName          string                 `json:"name,omitempty" multipart:"name"`
	ConditionSets         map[string]interface{} `json:"condition_sets" multipart:"conditionSets"`
	LookBackWindow        int                    `json:"lookback_window" multipart:"loopbackWindow"`

	Size         string                 `json:"size,omitempty"`
	AudienceData map[string]interface{} `json:"fb_audience_data,omitempty"`
	SourceType   string                 `json:"source_type,omitempty"`
}

var (
	graphqlStringField = &graphql.Field{
		Type:graphql.String,
	}
	grapqhqlIntField = &graphql.Field{
		Type: graphql.Float,
	}
	graphqlBoolField = &graphql.Field{
		Type: graphql.Boolean,
	}
)


var FbAudienceType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "FacebookAudienceSegment",
		Fields: graphql.Fields{
			"id":graphqlStringField,
			"instance_id": graphqlStringField,
			"advertiser_id": graphqlStringField,
			"fb_account_id": graphqlStringField,
			"type": graphqlStringField,
			"event_source_group": graphqlStringField,
			"dmp_audience_segment_ref": graphqlStringField,
			"service_segment": graphqlStringField,
			"created_time": grapqhqlIntField,
			"modified_time": grapqhqlIntField,
			"archived": graphqlBoolField,
			"source_type": graphqlStringField,
			"size": grapqhqlIntField,
			"fb_audience_data": &graphql.Field{
				Type: graphql.String,
			},
		},

	})

var QueryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"audiences": &graphql.Field{
				Type: graphql.NewList(FbAudienceType),
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"type": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
					log.Println(p.Info.Fragments)
					whereSt := ""
					audienceID, isOk := p.Args["id"].(string)
					if isOk {
						whereSt  = whereSt + fmt.Sprintf(" id='%s'", audienceID)
					}
					audienceType, isOk := p.Args["type"].(string)
					if isOk {
						whereSt = whereSt + fmt.Sprintf(" type='%s'", audienceType)
					}
					querySt := "SELECT * FROM facebook_audience_segment where archived=false AND " + whereSt
					d, err := dao.Db.Query(querySt)
					if err != nil {
						return nil, err
					}
					defer d.Close()
					columns, _ := d.Columns()
					res := make([]interface{}, 0)

					dest := make([]interface{}, len(columns)) // A temporary interface{} slice
					result := make([]interface{}, len(columns))
					for i := range columns {
						dest[i] = &result[i]// Put pointers to each string in the interface slice
					}
					for d.Next() {
						err = d.Scan(dest...)
						if err != nil {
							return res, err
						}
						obj := make(map[string]interface{}, len(columns))
						for idx := range columns {
							obj[columns[idx]] = result[idx]
							if columns[idx] == "created_time" {
								log.Println(reflect.TypeOf(result[idx]))
							}
						}
						log.Println(obj)
						res = append(res, obj)
					}
					return res, nil
				},

			},
		},
	})

var FbAudienceSchema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: QueryType,
	})
func ExecuteQuery(query string, schema graphql.Schema) (*graphql.Result, error) {
	result := graphql.Do(graphql.Params{
		Schema: schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		log.Println()
		return nil, errors.New(result.Errors[0].Message)
	}
	return  result, nil
}