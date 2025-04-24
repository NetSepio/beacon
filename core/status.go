package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/NetSepio/beacon/util"
	"github.com/NetSepio/beacon/web3"
	log "github.com/sirupsen/logrus"
)

// NodeStatus represents the current status of a node
type NodeStatus struct {
	ID        string
	Name      string
	Spec      string
	Config    string
	IPAddress string
	Region    string
}

// GetNodeStatus returns the current status of the node
func GetNodeStatus() (*web3.NodeStatus, error) {
	// Get node status from blockchain
	status, err := web3.GetNodeStatus()
	if err != nil {
		return nil, fmt.Errorf("failed to get node status from blockchain: %v", err)
	}

	return status, nil
} 