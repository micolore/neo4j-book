package main

import (
	"log"
	"strconv"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

var (
	neo4jURL = "bolt://localhost:17687"
)

func CreateDriver(uri, username, password string) (neo4j.Driver, error) {
	return neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
}

func CloseDriver(driver neo4j.Driver) error {
	return driver.Close()
}

func CypherWrite(driver neo4j.Driver, Cypher string, DB string) error {

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(Cypher, nil)
		if err != nil {
			log.Println("write to DB with error:", err)
			return nil, err
		}
		return result.Consume()
	})

	return err
}

func NodeQuery(driver neo4j.Driver, Cypher string, DB string) ([]neo4j.Node, error) {

	var list []neo4j.Node
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	_, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(Cypher, nil)
		if err != nil {
			return nil, err
		}

		for result.Next() {
			record := result.Record()
			if value, ok := record.Get("p"); ok {
				node := value.(neo4j.Node)
				list = append(list, node)
			}
		}
		if err = result.Err(); err != nil {
			return nil, err
		}

		return list, result.Err()
	})

	if err != nil {
		log.Println("Read error:", err)
	}
	return list, err
}

func RelationshipQuery(driver neo4j.Driver, Cypher string, DB string) ([]neo4j.Relationship, error) {

	var list []neo4j.Relationship
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	_, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(Cypher, nil)
		if err != nil {
			log.Println("RelationshipQuery failed: ", err)
			return nil, err
		}
		for result.Next() {
			record := result.Record()
			if value, ok := record.Get("r"); ok {
				relationship := value.(neo4j.Relationship)
				list = append(list, relationship)
			}
		}
		if err = result.Err(); err != nil {
			return nil, err
		}
		return list, result.Err()
	})

	if err != nil {
		log.Println("Read error:", err)
	}
	return list, err
}

type Node struct {
	NodeId   string `json:"id"`
	ObjId    string `json:"objId" `
	NodeName string `json:"node_name" `
	Desc     string `json:"desc"`
}

type Relation struct {
	RelationId   string `json:"id" `
	RelationName string `json:"relation_name" `
	StartId      string `json:"source" `
	EndId        string `json:"target" `
}

type NodeData struct {
	Data Node `json:"data" form:"data"`
}
type RelationData struct {
	Data Relation `json:"data" form:"data"`
}

// Neo4j-测试获取结果集
func TestGetNode(t *testing.T) {
	nodes := make([]NodeData, 0)
	driver, err := CreateDriver(neo4jURL, "neo4j", "123456")
	defer func(driver neo4j.Driver) {
		err = CloseDriver(driver)
		if err != nil {
			log.Println("neo4j close error:", err)
		}
	}(driver)
	if err != nil {
		log.Println("error connecting to neo4j:", err)
	}

	data, err := NodeQuery(driver, "match (p:Person)-[:LOVES]->(d:Dog) return p,d", "")

	for i := 0; i < len(data); i++ {
		var node NodeData
		node.Data.NodeId = strconv.FormatInt(data[i].Id, 10)
		node.Data.NodeName = data[i].Props["name"].(string)
		node.Data.Desc = data[i].Props["desc"].(string)
		nodes = append(nodes, node)
	}

}
