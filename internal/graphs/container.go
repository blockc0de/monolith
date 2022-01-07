package graphs

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/blockc0de/engine"
	"github.com/blockc0de/engine/block"
	"github.com/blockc0de/engine/compress"
	"github.com/blockc0de/engine/interop"
	"github.com/blockc0de/monolith/internal/storage"
	"github.com/blockc0de/monolith/internal/types"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stores/redis"
)

type Graph struct {
	Hash  string
	Graph *block.Graph
}

type Container struct {
	redisClient        *redis.Redis
	mutex              sync.Mutex
	graphs             map[string]*engine.Engine
	pendingQueue       chan Graph
	restartingSet      map[string]*block.Graph
	restartingSetMutex sync.Mutex
}

func NewContainer(redisClient *redis.Redis) *Container {
	c := Container{
		graphs:        make(map[string]*engine.Engine),
		restartingSet: make(map[string]*block.Graph),
		redisClient:   redisClient,
		pendingQueue:  make(chan Graph, 1024),
	}
	go c.polling()
	return &c
}

func (c *Container) LoadGraphs() (int, error) {
	graphsManager := storage.GraphsManager{
		RedisClient: c.redisClient,
	}

	var cursor uint64
	var graphs []*block.Graph
	for {
		slice, curr, err := graphsManager.Scan(cursor, 1)
		if err != nil {
			return 0, err
		}

		for _, graph := range slice {
			n, err := strconv.Atoi(graph.State)
			if err != nil {
				continue
			}

			state := GraphStateEnum(n)
			if state != GraphStateEnumStarting && state != GraphStateEnumStarted && state != GraphStateEnumRestarting {
				continue
			}

			data, err := compress.GraphCompression{}.DecompressGraphData(graph.Data)
			if err != nil {
				return 0, err
			}

			instance, err := interop.LoadGraph(data)
			if err != nil {
				return 0, err
			}

			instance.Hash = graph.Hash
			graphs = append(graphs, instance)
		}

		if curr == 0 {
			break
		}
		cursor = curr
	}

	for _, graph := range graphs {
		c.AddNewGraph(graph.Hash, graph)
	}

	return len(graphs), nil
}

func (c *Container) AddNewGraph(hash string, graph *block.Graph) {
	c.mutex.Lock()
	engine, ok := c.graphs[hash]
	c.mutex.Unlock()

	if !ok {
		c.updateGraphState(hash, GraphStateEnumStarting)
		c.pendingQueue <- Graph{Hash: hash, Graph: graph}
		return
	}

	c.restartingSetMutex.Lock()
	c.restartingSet[hash] = graph
	c.restartingSetMutex.Unlock()

	engine.Stop()

	c.updateGraphState(hash, GraphStateEnumRestarting)
}

func (c *Container) StopGraphByHash(hash string) {
	c.mutex.Lock()
	engine, ok := c.graphs[hash]
	c.mutex.Unlock()

	if ok && engine != nil {
		engine.Stop()
	}
}

func (c *Container) polling() {
	for graph := range c.pendingQueue {
		go c.runEngine(graph.Hash, graph.Graph)
	}
}

func (c *Container) runEngine(hash string, graph *block.Graph) {
	ev := engine.Event{
		AppendLog: func(msgType string, message string) {
			c.appendLog(hash, msgType, message)
		},
	}
	engine := engine.NewEngine(graph, ev)
	c.updateGraphState(hash, GraphStateEnumStarted)

	c.mutex.Lock()
	if _, ok := c.graphs[hash]; ok {
		c.mutex.Unlock()
		logx.Errorf("Abandon run graph, the graph is running, hash: %s", hash)
		return
	}
	c.graphs[hash] = engine
	c.mutex.Unlock()

	logx.Infof("Graph hash %s started", hash)

	// Threads are blocked when engine running
	err := engine.Run(context.Background())
	if err == nil {
		logx.Infof("Graph hash %s stopped", hash)
	} else {
		logx.Errorf("Engine encountered an error, reason: %s", err.Error())
	}

	c.mutex.Lock()
	delete(c.graphs, hash)
	c.mutex.Unlock()

	if err != nil {
		c.updateGraphState(hash, GraphStateEnumError)
	} else {
		c.updateGraphState(hash, GraphStateEnumStopped)
	}

	// Restart graph if necessary
	c.restartingSetMutex.Lock()
	graph, ok := c.restartingSet[hash]
	if ok {
		delete(c.restartingSet, hash)
	}
	c.restartingSetMutex.Unlock()

	if ok {
		engine.AppendLog("warn",
			fmt.Sprintf("Graph hash %s stopped successfully, restarting...", hash))
		logx.Infof("Graph hash %s stopped successfully, restarting...", hash)
		c.pendingQueue <- Graph{Hash: hash, Graph: graph}
	}
}

func (c *Container) appendLog(hash string, msgType string, message string) {
	graphsManager := storage.GraphsManager{
		RedisClient: c.redisClient,
	}

	log := types.Log{
		Type:      msgType,
		Message:   message,
		Timestamp: time.Now().UnixMilli(),
	}
	err := graphsManager.AppendLog(hash, log)
	if err != nil {
		logx.Errorf("Failed to append log, reason: %s", err.Error())
	}
}

func (c *Container) updateGraphState(hash string, state GraphStateEnum) {
	graphsManager := storage.GraphsManager{
		RedisClient: c.redisClient,
	}

	err := graphsManager.SetState(hash, strconv.Itoa(int(state)))
	if err != nil {
		logx.Errorf("Failed to update graph state, reason: %s", err.Error())
	}
}
