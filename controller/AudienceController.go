package controller

import (
	"github.com/gin-gonic/gin"
	"graphql-taxonomy/model/graph-type"
	"log"
	"strings"
)

func GetAudience(c *gin.Context)  {
	var (
		q,_ = c.GetRawData()

	)
	//fbAudience := graph_type.FacebookAudienceSegment{}
	responseData := graph_type.FbAudienceType.Fields()
	log.Println("response schema:")
	log.Println(responseData)
	//fields := reflect.TypeOf(fbAudience)
	jsonFieldsName := make([]string, 0)

	for k:= range responseData {
		jsonFieldsName = append(jsonFieldsName, responseData[k].Name)
	}
	fragment := string(" fragment facebookAudienceFields on FacebookAudienceSegment {" + strings.Join(jsonFieldsName, "\n") +
		"}")
	query := string(q)
	query = query + fragment
	log.Println(query)
	result, err := executeQuery(query, graph_type.FbAudienceSchema)
	response := make(map[string]interface{})
	response["data"] = result.Data
	if err == nil {
		response["success"] = true
		c.JSON(200, response)
	} else {
		response["success"] = false
		response["message"] = err.Error()
		c.JSON(500, response)
	}

}
