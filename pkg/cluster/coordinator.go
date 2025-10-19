package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

type Node struct {
	ID          string            `json:"id"`
	Address     string            `json:"address"`
	Port        int               `json:"port"`
	Status      NodeStatus        `json:"status"`
	Capabilities []string         `json:"capabilities"`
	Load        *NodeLoad         `json:"load"`
	Metadata    map[string]string `json:"metadata"`
	LastSeen    time.Time         `json:"last_seen"`
}

type NodeStatus string

const (
	NodeStatusActive   NodeStatus = "active"
	NodeStatusInactive NodeStatus = "inactive"
	NodeStatusDraining NodeStatus = "draining"
	NodeStatusFailed   NodeStatus = "failed"
)

type NodeLoad struct {
	CPU        float64 `json:"cpu"`
	Memory     float64 `json:"memory"`
	ActiveJobs int     `json:"active_jobs"`
	QueueSize  int     `json:"queue_size"`
}

type Coordinator interface {
	RegisterNode(ctx context.Context, node *Node) error
	UnregisterNode(ctx context.Context, nodeID string) error
	GetNodes(ctx context.Context) ([]*Node, error)
	GetNode(ctx context.Context, nodeID string) (*Node, error)
	UpdateNodeLoad(ctx context.Context, nodeID string, load *NodeLoad) error
	DistributeJob(ctx context.Context, job *Job) (*Node, error)
	ElectLeader(ctx context.Context) (string, error)
	IsLeader(ctx context.Context) (bool, error)
	WatchNodes(ctx context.Context) (<-chan NodeEvent, error)
}

type Job struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"`
	Priority    int               `json:"priority"`
	Requirements []string         `json:"requirements"`
	Payload     interface{}       `json:"payload"`
	CreatedAt   time.Time         `json:"created_at"`
	AssignedTo  string            `json:"assigned_to,omitempty"`
}

type NodeEvent struct {
	Type EventType `json:"type"`
	Node *Node     `json:"node"`
}

type EventType string

const (
	EventNodeJoined  EventType = "node_joined"
	EventNodeLeft    EventType = "node_left"
	EventNodeUpdated EventType = "node_updated"
	EventNodeFailed  EventType = "node_failed"
)

type ConsulCoordinator struct {
	client    *api.Client
	config    *ConsulConfig
	logger    *zap.Logger
	nodeID    string
	leaderKey string
	mu        sync.RWMutex
	nodes     map[string]*Node
}

type ConsulConfig struct {
	Address    string `json:"address"`
	Datacenter string `json:"datacenter"`
	Token      string `json:"token"`
	Prefix     string `json:"prefix"`
}

func NewConsulCoordinator(config *ConsulConfig, nodeID string, logger *zap.Logger) (*ConsulCoordinator, error) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = config.Address
	consulConfig.Datacenter = config.Datacenter
	consulConfig.Token = config.Token

	client, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}

	return &ConsulCoordinator{
		client:    client,
		config:    config,
		logger:    logger,
		nodeID:    nodeID,
		leaderKey: fmt.Sprintf("%s/leader", config.Prefix),
		nodes:     make(map[string]*Node),
	}, nil
}

func (c *ConsulCoordinator) RegisterNode(ctx context.Context, node *Node) error {
	key := fmt.Sprintf("%s/nodes/%s", c.config.Prefix, node.ID)
	
	data, err := json.Marshal(node)
	if err != nil {
		return fmt.Errorf("failed to marshal node: %w", err)
	}

	session := &api.SessionEntry{
		Name:      fmt.Sprintf("node-%s", node.ID),
		TTL:       "30s",
		Behavior:  api.SessionBehaviorDelete,
		LockDelay: time.Second,
	}

	sessionID, _, err := c.client.Session().Create(session, nil)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	kv := &api.KVPair{
		Key:     key,
		Value:   data,
		Session: sessionID,
	}

	_, err = c.client.KV().Put(kv, nil)
	if err != nil {
		return fmt.Errorf("failed to register node: %w", err)
	}

	go c.renewSession(ctx, sessionID)

	c.mu.Lock()
	c.nodes[node.ID] = node
	c.mu.Unlock()

	c.logger.Info("Node registered", zap.String("node_id", node.ID))
	return nil
}

func (c *ConsulCoordinator) UnregisterNode(ctx context.Context, nodeID string) error {
	key := fmt.Sprintf("%s/nodes/%s", c.config.Prefix, nodeID)
	
	_, err := c.client.KV().Delete(key, nil)
	if err != nil {
		return fmt.Errorf("failed to unregister node: %w", err)
	}

	c.mu.Lock()
	delete(c.nodes, nodeID)
	c.mu.Unlock()

	c.logger.Info("Node unregistered", zap.String("node_id", nodeID))
	return nil
}

func (c *ConsulCoordinator) GetNodes(ctx context.Context) ([]*Node, error) {
	prefix := fmt.Sprintf("%s/nodes/", c.config.Prefix)
	
	pairs, _, err := c.client.KV().List(prefix, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	var nodes []*Node
	for _, pair := range pairs {
		var node Node
		if err := json.Unmarshal(pair.Value, &node); err != nil {
			c.logger.Warn("Failed to unmarshal node", zap.Error(err))
			continue
		}
		nodes = append(nodes, &node)
	}

	return nodes, nil
}

func (c *ConsulCoordinator) GetNode(ctx context.Context, nodeID string) (*Node, error) {
	key := fmt.Sprintf("%s/nodes/%s", c.config.Prefix, nodeID)
	
	pair, _, err := c.client.KV().Get(key, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	if pair == nil {
		return nil, fmt.Errorf("node not found: %s", nodeID)
	}

	var node Node
	if err := json.Unmarshal(pair.Value, &node); err != nil {
		return nil, fmt.Errorf("failed to unmarshal node: %w", err)
	}

	return &node, nil
}

func (c *ConsulCoordinator) UpdateNodeLoad(ctx context.Context, nodeID string, load *NodeLoad) error {
	node, err := c.GetNode(ctx, nodeID)
	if err != nil {
		return err
	}

	node.Load = load
	node.LastSeen = time.Now()

	return c.RegisterNode(ctx, node)
}

func (c *ConsulCoordinator) DistributeJob(ctx context.Context, job *Job) (*Node, error) {
	nodes, err := c.GetNodes(ctx)
	if err != nil {
		return nil, err
	}

	var bestNode *Node
	var bestScore float64

	for _, node := range nodes {
		if node.Status != NodeStatusActive {
			continue
		}

		if !c.nodeSupportsJob(node, job) {
			continue
		}

		score := c.calculateNodeScore(node, job)
		if bestNode == nil || score > bestScore {
			bestNode = node
			bestScore = score
		}
	}

	if bestNode == nil {
		return nil, fmt.Errorf("no suitable node found for job")
	}

	return bestNode, nil
}

func (c *ConsulCoordinator) ElectLeader(ctx context.Context) (string, error) {
	session := &api.SessionEntry{
		Name:      fmt.Sprintf("leader-%s", c.nodeID),
		TTL:       "30s",
		Behavior:  api.SessionBehaviorRelease,
		LockDelay: time.Second,
	}

	sessionID, _, err := c.client.Session().Create(session, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	kv := &api.KVPair{
		Key:     c.leaderKey,
		Value:   []byte(c.nodeID),
		Session: sessionID,
	}

	acquired, _, err := c.client.KV().Acquire(kv, nil)
	if err != nil {
		return "", fmt.Errorf("failed to acquire leader lock: %w", err)
	}

	if acquired {
		go c.renewSession(ctx, sessionID)
		return c.nodeID, nil
	}

	pair, _, err := c.client.KV().Get(c.leaderKey, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get current leader: %w", err)
	}

	if pair != nil {
		return string(pair.Value), nil
	}

	return "", fmt.Errorf("no leader found")
}

func (c *ConsulCoordinator) IsLeader(ctx context.Context) (bool, error) {
	pair, _, err := c.client.KV().Get(c.leaderKey, nil)
	if err != nil {
		return false, err
	}

	if pair != nil && string(pair.Value) == c.nodeID {
		return true, nil
	}

	return false, nil
}

func (c *ConsulCoordinator) WatchNodes(ctx context.Context) (<-chan NodeEvent, error) {
	eventCh := make(chan NodeEvent, 100)
	
	go func() {
		defer close(eventCh)
		
		prefix := fmt.Sprintf("%s/nodes/", c.config.Prefix)
		var lastIndex uint64
		
		for {
			select {
			case <-ctx.Done():
				return
			default:
				pairs, meta, err := c.client.KV().List(prefix, &api.QueryOptions{
					WaitIndex: lastIndex,
					WaitTime:  30 * time.Second,
				})
				
				if err != nil {
					c.logger.Error("Failed to watch nodes", zap.Error(err))
					time.Sleep(5 * time.Second)
					continue
				}
				
				lastIndex = meta.LastIndex
				
				c.processNodeChanges(pairs, eventCh)
			}
		}
	}()
	
	return eventCh, nil
}

func (c *ConsulCoordinator) renewSession(ctx context.Context, sessionID string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, _, err := c.client.Session().Renew(sessionID, nil)
			if err != nil {
				c.logger.Error("Failed to renew session", zap.Error(err))
				return
			}
		}
	}
}

func (c *ConsulCoordinator) nodeSupportsJob(node *Node, job *Job) bool {
	for _, req := range job.Requirements {
		found := false
		for _, cap := range node.Capabilities {
			if cap == req {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (c *ConsulCoordinator) calculateNodeScore(node *Node, job *Job) float64 {
	if node.Load == nil {
		return 0
	}

	cpuScore := 1.0 - node.Load.CPU
	memoryScore := 1.0 - node.Load.Memory
	jobScore := 1.0 / (float64(node.Load.ActiveJobs) + 1)

	priorityWeight := float64(job.Priority) / 10.0

	return (cpuScore + memoryScore + jobScore) * (1.0 + priorityWeight)
}

func (c *ConsulCoordinator) processNodeChanges(pairs api.KVPairs, eventCh chan<- NodeEvent) {
	currentNodes := make(map[string]*Node)
	
	for _, pair := range pairs {
		var node Node
		if err := json.Unmarshal(pair.Value, &node); err != nil {
			continue
		}
		currentNodes[node.ID] = &node
	}
	
	c.mu.Lock()
	defer c.mu.Unlock()
	
	for id, node := range currentNodes {
		if _, exists := c.nodes[id]; !exists {
			eventCh <- NodeEvent{
				Type: EventNodeJoined,
				Node: node,
			}
		}
	}
	
	for id, node := range c.nodes {
		if _, exists := currentNodes[id]; !exists {
			eventCh <- NodeEvent{
				Type: EventNodeLeft,
				Node: node,
			}
		}
	}
	
	c.nodes = currentNodes
}