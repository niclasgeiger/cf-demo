package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"encoding/json"

	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

// Neo4jHandler returns a simple message
func Neo4jHandler(credentials *Credentials) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		numResult, err := createNeo4JRows(credentials.Username, credentials.Password, credentials.BoltUrl)
		if err != nil {
			fmt.Fprint(w, "Could not create rows: ", err)
		} else {
			fmt.Fprintf(w, "CREATED ROWS: %d\n", numResult)
		}
	}
}

type VcapServices struct {
	Neo4j []Neo4j `json:"Neo4j"`
}

type Neo4j struct {
	Credentials *Credentials `json:"credentials"`
}

type Credentials struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	BrowserUrl string `json:"browser_url"`
	BoltUrl    string `json:"bolt_url"`
}

func loadEnvConfig() (*Credentials, error) {
	vcapServices := os.Getenv("VCAP_SERVICES")
	var variables VcapServices
	if err := json.Unmarshal([]byte(vcapServices), &variables); err != nil {
		return nil, err
	}
	if len(variables.Neo4j) == 0 {
		return nil, errors.New("no neo4j service binding available")
	}
	return variables.Neo4j[0].Credentials, nil
}

func createNeo4JRows(username, password, boltUrl string) (int64, error) {
	driver := bolt.NewDriver()
	conn, err := driver.OpenNeo(fmt.Sprintf("bolt://%s:%s@%s:%s", username, password, boltUrl, "443"))
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	// Start by creating a node
	result, err := conn.ExecNeo("CREATE (n:NODE {foo: {foo}, bar: {bar}})", map[string]interface{}{"foo": 1, "bar": 2.2})
	if err != nil {
		return 0, err
	}
	numResult, _ := result.RowsAffected()
	return numResult, nil
}

func main() {
	credentials, err := loadEnvConfig()
	if err != nil {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "could not load service binding: ", err)
		})
	} else {
		http.HandleFunc("/", Neo4jHandler(credentials))

	}

	log.Fatal(http.ListenAndServe(":80", nil))
}
