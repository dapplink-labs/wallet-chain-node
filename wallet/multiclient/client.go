package multiclient

import (
	"sync"
	"time"

	"go.uber.org/atomic"
)

type Client interface {
	GetLatestBlockHeight() (int64, error)
}

type MultiClient struct {
	clients   []Client
	bestIndex atomic.Int32
}

func New(clients []Client) *MultiClient {
	m := &MultiClient{
		clients: clients,
	}
	if len(clients) > 1 {
		go m.sniffLoop()
	}
	return m
}

func (m *MultiClient) BestClient() Client {
	return m.clients[m.bestIndex.Load()]
}

func (m *MultiClient) sniffLoop() {
	t := time.NewTimer(0)
	for {
		select {
		case <-t.C:
			m.sniff()
			t.Reset(time.Second)
		}
	}
}

func (m *MultiClient) sniff() {
	var (
		heights = make([]int64, len(m.clients))
		times   = make([]int64, len(m.clients))
		l       sync.Mutex
		wg      sync.WaitGroup
	)
	wg.Add(len(m.clients))
	for i, client := range m.clients {
		i, client := i, client
		go func() {
			defer wg.Done()
			start := time.Now().UnixNano()
			height, _ := client.GetLatestBlockHeight()
			l.Lock()
			heights[i] = height
			times[i] = time.Now().UnixNano() - start
			l.Unlock()
		}()
	}
	wg.Wait()

	var (
		maxHeight  = heights[0]
		minTime    = times[0]
		bestClient = 0
	)
	for i := 1; i < len(m.clients); i++ {
		if heights[i] > maxHeight {
			maxHeight = heights[i]
			minTime = times[i]
			bestClient = i
		} else if heights[i] == maxHeight {
			if times[i] < minTime {
				minTime = times[i]
				bestClient = i
			}
		}
	}
	m.bestIndex.Store(int32(bestClient))
}
