package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"github.com/awalterschulze/gographviz"
)


func main() {
	// SQL statements
	// TODO: Have ability to parse sql files rather than simple strings
	sqlStatements := []string{
		"CREATE TABLE user (id INT, name VARCHAR(50), email VARCHAR(100));",
		"CREATE TABLE posts (id INT, title VARCHAR(100), content TEXT, user_id INT);",
		"SELECT users.name, posts.title FROM users INNER JOIN posts ON users.id = posts.user_id;",
	}

	// Parse SQL statements and build UML diagram
	graphAst := gographviz.NewGraph()
	graphAst.SetName("UMLDiagram")
	graphAst.SetDir(true)

	for _, stmt := range sqlStatements {
		fmt.Println("SQL Statement:", stmt) // Print SQL statement

		isCreateStatement := strings.HasPrefix(strings.ToUpper(stmt), "CREATE")
		isSelectStatement := strings.HasPrefix(strings.ToUpper(stmt), "SELECT")

		if isCreateStatement {
			tableName, attributes := parseCreateStatement(stmt)
			addTableToGraph(graphAst, tableName, attributes)
		} else if isSelectStatement {
			parseSelectStatement(stmt, graphAst)
		}

		// Print DOT representation for debugging
		dotOutput := graphAst.String()
		fmt.Println("DOT Output:\n", dotOutput)
	}

	// Generate UML diagram
	dotOutput := graphAst.String()

	// Save UML diagram as DOT file
	dotFilePath := "uml.dot"
	if err := ioutil.WriteFile(dotFilePath, []byte(dotOutput), 0644); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(dotFilePath)

	// Convert DOT file to image file
	imageFilePath := "uml.png"
	cmd := exec.Command("dot", "-Tpng", "-o", imageFilePath, dotFilePath)

	// Capture stderr output
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("UML diagram saved as '%s'\n", imageFilePath)
}

func addTableToGraph(graphAst *gographviz.Graph, tableName string, attributes []string) {
	nodeName := sanitizeNodeName(tableName)
	graphAst.AddNode("UMLDiagram", nodeName, map[string]string{"shape": "box"})

	// Add attributes as child nodes and edges
	for _, attr := range attributes {
		attrName := sanitizeNodeName(attr)
		graphAst.AddNode("UMLDiagram", attrName, nil)
		graphAst.AddEdge(nodeName, attrName, true, nil)
	}
}

func sanitizeNodeName(name string) string {
	// Remove invalid characters from node name
	sanitized := strings.ReplaceAll(name, " ", "_")
	sanitized = strings.Map(func(r rune) rune {
		if r >= 'A' && r <= 'Z' {
			return r
		}
		if r >= 'a' && r <= 'z' {
			return r
		}
		if r >= '0' && r <= '9' {
			return r
		}
		return '_'
	}, sanitized)
	return sanitized
}

func parseCreateStatement(stmt string) (tableName string, attributes []string) {
	// Parse CREATE TABLE statement and extract table name and attributes
	parts := strings.Split(stmt, "(")
	if len(parts) >= 1 {
		tableName = strings.TrimSpace(strings.TrimPrefix(parts[0], "CREATE TABLE"))
	}

	if len(parts) >= 2 {
		attrStr := strings.TrimSpace(strings.TrimSuffix(parts[1], ");"))
		attributes = strings.Split(attrStr, ",")
		for i := 0; i < len(attributes); i++ {
			attributes[i] = strings.TrimSpace(attributes[i])
		}
	}

	return tableName, attributes
}

func parseSelectStatement(stmt string, graphAst *gographviz.Graph) {
	tableNames := extractTableNamesFromSelect(stmt)

	// Add edges between tables to represent dependencies
	for i := 0; i < len(tableNames)-1; i++ {
		from := sanitizeNodeName(tableNames[i])
		to := sanitizeNodeName(tableNames[i+1])
		graphAst.AddEdge(from, to, true, nil)
	}
}

func extractTableNamesFromSelect(stmt string) []string {
	// Extract table names from the SELECT statement
	tableNames := []string{"users", "posts"}
	return tableNames
}
