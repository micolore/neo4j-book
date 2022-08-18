package main

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"testing"
)

func TestBasic(t *testing.T) {
	ConnectionNeo4j()
}

func ConnectionNeo4j() {
	dbUri := "neo4j://localhost:17687"
	driver, err := neo4j.NewDriver(dbUri, neo4j.BasicAuth("neo4j", "123456", ""))
	if err != nil {
		panic(err)
	}
	defer driver.Close()
	item, err := insertItem(driver)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", item)
}

func insertItem(driver neo4j.Driver) (*Item, error) {
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()
	result, err := session.WriteTransaction(createItemFn)
	if err != nil {
		return nil, err
	}
	return result.(*Item), nil
}

func createItemFn(tx neo4j.Transaction) (interface{}, error) {
	records, err := tx.Run("CREATE (n:Item { id: $id, name: $name }) RETURN n.id, n.name", map[string]interface{}{
		"id":   1,
		"name": "Item 1",
	})
	if err != nil {
		return nil, err
	}
	record, err := records.Single()
	if err != nil {
		return nil, err
	}
	return &Item{
		Id:   record.Values[0].(int64),
		Name: record.Values[1].(string),
	}, nil
}

type Item struct {
	Id   int64
	Name string
}
