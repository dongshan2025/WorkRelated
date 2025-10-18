// https://github.com/couchbase/gocb
// go get github.com/couchbase/gocb/v2@latest   最新
// go get github.com/couchbase/gocb@v1.6.7 低版本
package main

import (
	"fmt"
	"log"

	"github.com/couchbase/gocb"
)

func main() {
	// cluster, err := gocb.Connect("couchbase://localhost")
	cluster, err := gocb.Connect("http://localhost:8091")
	if err != nil {
		log.Fatal(err)
		return
	}

	err = cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "123456",
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	bucket, err := cluster.OpenBucket("User", "")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Example operation: Upsert a document
	_, err = bucket.Upsert("my_document_id", map[string]string{"name": "John Doe2"}, 0)
	if err != nil {
		fmt.Printf("Error upserting document: %s\n", err)
		return
	}

	fmt.Println("Document upserted successfully!")

	// Close the bucket and cluster connection
	bucket.Close()
	cluster.Close()
}
