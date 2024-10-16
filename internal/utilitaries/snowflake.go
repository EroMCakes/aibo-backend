package utilitaries

import (
	"hash/fnv"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/bwmarrin/snowflake"
)

var (
	snowflakeNode *snowflake.Node
	nodeOnce      sync.Once
)

// getNodeID determines a unique Snowflake node ID from the current pod name.
//
// If the POD_NAME environment variable is not set, it defaults to 1.
//
// If the POD_NAME environment variable is set, it attempts to parse the
// last part of the pod name as an integer. If that fails, it hashes the pod
// name and takes the result modulo 1024.
func getNodeID() int64 {
	podName := os.Getenv("POD_NAME")
	if podName == "" {
		return 1 // Fallback to a default value
	}

	parts := strings.Split(podName, "-")
	lastPart := parts[len(parts)-1]

	if id, err := strconv.ParseInt(lastPart, 10, 64); err == nil {
		return id % 1024
	}

	h := fnv.New32a()
	h.Write([]byte(podName))
	return int64(h.Sum32()) % 1024
}

// initSnowflakeNode initializes the Snowflake node used for generating IDs.
//
// The node is configured with a unique ID based on the current pod name.
// If the POD_NAME environment variable is not set, it defaults to 1.
// Otherwise, it attempts to parse the last part of the pod name as an
// integer; if that fails, it hashes the pod name and takes the result
// modulo 1024.
//
// The node is created once and then reused for all subsequent calls.
func initSnowflakeNode() {
	nodeID := getNodeID()
	var err error
	snowflakeNode, err = snowflake.NewNode(nodeID)
	if err != nil {
		panic(err) // Or handle this error as appropriate for your application
	}
}

// GenerateSnowflakeID generates a new Snowflake ID
func GenerateSnowflakeID() snowflake.ID {
	nodeOnce.Do(initSnowflakeNode)
	sfId := snowflakeNode.Generate() // Ensures the node is initialized only once
	log.Printf("Generated Snowflake ID: %d", sfId)
	return sfId
}
