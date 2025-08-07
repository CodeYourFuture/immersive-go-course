// Test harness for writing tests for Raft.
//
// Eli Bendersky [https://eli.thegreenplace.net]
// This code is in the public domain.
package raft

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func init() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	seed := time.Now().UnixNano()
	fmt.Println("seed", seed)
	rand.Seed(seed)
}

type Harness struct {
	mu sync.Mutex

	// cluster is a list of all the raft servers participating in a cluster.
	cluster map[string]*Server
	storage map[string]*MapStorage

	// commitChans has a channel per server in cluster with the commi channel for
	// that server.
	commitChans map[string]chan CommitEntry

	// commits at index i holds the sequence of commits made by server i so far.
	// It is populated by goroutines that listen on the corresponding commitChans
	// channel.
	commits map[string][]CommitEntry

	// connected has a bool per server in cluster, specifying whether this server
	// is currently connected to peers (if false, it's partitioned and no messages
	// will pass to or from it).
	connected map[string]bool

	// alive has a bool per server in cluster, specifying whether this server is
	// currently alive (false means it has crashed and wasn't restarted yet).
	// connected implies alive.
	alive map[string]bool

	n int
	t *testing.T
}

// NewHarness creates a new test Harness, initialized with n servers connected
// to each other.
func NewHarness(t *testing.T, n int) *Harness {
	ns := make(map[string]*Server)
	connected := make(map[string]bool, n)
	alive := make(map[string]bool, n)
	commitChans := make(map[string]chan CommitEntry)
	commits := make(map[string][]CommitEntry)
	ready := make(chan interface{})
	storage := make(map[string]*MapStorage)

	// Create all Servers in this cluster, assign ids and peer ids.
	for i := 0; i < n; i++ {
		serverId := fmt.Sprintf("%d", i)
		storage[serverId] = NewMapStorage()
		commitChans[serverId] = make(chan CommitEntry)
		ns[serverId] = NewServer(serverId, "127.0.0.1", storage[serverId], ready, commitChans[serverId], 7601+i)
		ns[serverId].Serve(nil)
		alive[serverId] = true
	}

	// Connect all peers to each other.
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j {
				ns[fmt.Sprintf("%d", i)].ConnectToPeer(fmt.Sprintf("%d", j), fmt.Sprintf("127.0.0.1:%d", 7601+j))
			}
		}
		connected[fmt.Sprintf("%d", i)] = true
	}
	close(ready)

	h := &Harness{
		cluster:     ns,
		storage:     storage,
		commitChans: commitChans,
		commits:     commits,
		connected:   connected,
		alive:       alive,
		n:           n,
		t:           t,
	}
	for i := 0; i < n; i++ {
		go h.collectCommits(i)
	}
	return h
}

// Shutdown shuts down all the servers in the harness and waits for them to
// stop running.
func (h *Harness) Shutdown() {
	for i := 0; i < h.n; i++ {
		h.cluster[str(i)].DisconnectAll()
		h.connected[str(i)] = false
	}
	for i := 0; i < h.n; i++ {
		if h.alive[str(i)] {
			h.alive[str(i)] = false
			h.cluster[str(i)].Shutdown()
		}
	}
	for i := 0; i < h.n; i++ {
		close(h.commitChans[str(i)])
	}
}

// DisconnectPeer disconnects a server from all other servers in the cluster.
func (h *Harness) DisconnectPeer(id string) {
	tlog("Disconnect %s", id)
	h.cluster[id].DisconnectAll()
	for peer, server := range h.cluster {
		if peer != id {
			server.DisconnectPeer(id)
		}
	}
	h.connected[id] = false
}

// ReconnectPeer connects a server to all other servers in the cluster.
func (h *Harness) ReconnectPeer(id string) {
	tlog("Reconnect %s", id)
	for peerId, peerServer := range h.cluster {
		if peerId != id && h.alive[peerId] {
			if err := peerServer.ConnectToPeer(id, fmt.Sprintf("127.0.0.1:%d", 7601+toInt(id))); err != nil {
				h.t.Fatal(err)
			}
			if err := h.cluster[id].ConnectToPeer(peerId, fmt.Sprintf("127.0.0.1:%d", 7601+toInt(peerId))); err != nil {
				h.t.Fatal(err)
			}
		}
	}
	h.connected[id] = true
}

// CrashPeer "crashes" a server by disconnecting it from all peers and then
// asking it to shut down. We're not going to use the same server instance
// again, but its storage is retained.
func (h *Harness) CrashPeer(id string) {
	tlog("Crash %s", id)
	h.DisconnectPeer(id)
	h.alive[id] = false
	h.cluster[id].Shutdown()

	// Clear out the commits slice for the crashed server; Raft assumes the client
	// has no persistent state. Once this server comes back online it will replay
	// the whole log to us.
	h.mu.Lock()
	h.commits[id] = h.commits[id][:0]
	h.mu.Unlock()
}

// RestartPeer "restarts" a server by creating a new Server instance and giving
// it the appropriate storage, reconnecting it to peers.
func (h *Harness) RestartPeer(id string) {
	if h.alive[id] {
		log.Fatalf("id=%s is alive in RestartPeer", id)
	}
	tlog("Restart %s", id)

	ready := make(chan interface{})
	portOffset, _ := strconv.Atoi(id) // should not discard err if not in a test

	h.cluster[id] = NewServer(id, "127.0.0.1", h.storage[id], ready, h.commitChans[id], 7601+portOffset)
	h.cluster[id].Serve(nil)
	h.ReconnectPeer(id)
	close(ready)
	h.alive[id] = true
	sleepMs(20)
}

// CheckSingleLeader checks that only a single server thinks it's the leader.
// Returns the leader's id and term. It retries several times if no leader is
// identified yet.
func (h *Harness) CheckSingleLeader() (string, int) {
	for r := 0; r < 8; r++ {
		leaderId := -1
		leaderTerm := -1
		for i := 0; i < h.n; i++ {
			if h.connected[str(i)] {
				_, term, isLeader := h.cluster[str(i)].cm.Report()
				if isLeader {
					if leaderId < 0 {
						leaderId = i
						leaderTerm = term
					} else {
						h.t.Fatalf("both %d and %d think they're leaders", leaderId, i)
					}
				}
			}
		}
		if leaderId >= 0 {
			return str(leaderId), leaderTerm
		}
		time.Sleep(150 * time.Millisecond)
	}

	h.t.Fatalf("leader not found")
	return "", -1
}

// CheckNoLeader checks that no connected server considers itself the leader.
func (h *Harness) CheckNoLeader() {
	for i := 0; i < h.n; i++ {
		if h.connected[str(i)] {
			_, _, isLeader := h.cluster[str(i)].cm.Report()
			if isLeader {
				h.t.Fatalf("server %d leader; want none", i)
			}
		}
	}
}

func (h *Harness) printCommits() {
	fmt.Printf("Printing commits for connected servers\n")
	for i := 0; i < h.n; i++ {
		if h.connected[str(i)] {

			fmt.Printf("[%s] is connected, commits are\n", str(i))
			for _, commit := range h.commits[str(i)] {
				fmt.Printf("Commit: %+v\n", commit)
			}
		}
	}

}

// CheckCommitted verifies that all connected servers have cmd committed with
// the same index. It also verifies that all commands *before* cmd in
// the commit sequence match. For this to work properly, all commands submitted
// to Raft should be unique positive ints.
// Returns the number of servers that have this command committed, and its
// log index.
// TODO: this check may be too strict. Consider tha a server can commit
// something and crash before notifying the channel. It's a valid commit but
// this checker will fail because it may not match other servers. This scenario
// is described in the paper...
func (h *Harness) CheckCommitted(cmdNum int) (nc int, index int) {
	cmd := cmdFor(cmdNum)

	h.mu.Lock()
	defer h.mu.Unlock()

	// Find the length of the commits slice for connected servers.
	commitsLen := -1
	for i := 0; i < h.n; i++ {
		if h.connected[str(i)] {
			if commitsLen >= 0 {
				// If this was set already, expect the new length to be the same.
				if len(h.commits[str(i)]) != commitsLen {
					h.t.Fatalf("Fatal error: commits[%d] = %v, want commitsLen = %d", i, h.commits[str(i)], commitsLen)
				}
			} else {
				commitsLen = len(h.commits[str(i)])
			}
		}
	}

	// Check consistency of commits from the start and to the command we're asked
	// about. This loop will return once a command=cmd is found.
	for c := 0; c < commitsLen; c++ {
		cmdAtC := cmdFor(-1)
		for i := 0; i < h.n; i++ {
			if h.connected[str(i)] {
				cmdOfN := h.commits[str(i)][c].Command
				if cmdInt(cmdAtC) >= 0 {
					if !equalsCmd(cmdOfN, cmdAtC) {
						h.t.Errorf("got %v, want %v at h.commits[%d][%d]", cmdOfN, cmdAtC, i, c)
					}
				} else {
					cmdAtC = cmdOfN
				}
			}
		}
		if equalsCmd(cmdAtC, cmd) {
			// Check consistency of Index.
			index := -1
			nc := 0
			for i := 0; i < h.n; i++ {
				if h.connected[str(i)] {
					if index >= 0 && h.commits[str(i)][c].Index != index {
						h.t.Errorf("got Index=%d, want %d at h.commits[%d][%d]", h.commits[str(i)][c].Index, index, i, c)
					} else {
						index = h.commits[str(i)][c].Index
					}
					nc++
				}
			}
			return nc, index
		}
	}

	// If there's no early return, we haven't found the command we were looking
	// for.
	h.t.Errorf("cmd=%v not found in commits", cmd)
	return -1, -1
}

// CheckCommittedN verifies that cmd was committed by exactly n connected
// servers.
func (h *Harness) CheckCommittedN(cmd int, n int) {
	nc, _ := h.CheckCommitted(cmd)
	if nc != n {
		h.t.Errorf("CheckCommittedN got nc=%d, want %d", nc, n)
	}
}

// CheckNotCommitted verifies that no command equal to cmd has been committed
// by any of the active servers yet.
func (h *Harness) CheckNotCommitted(cmdNum int) {
	cmd := cmdFor(cmdNum)
	h.mu.Lock()
	defer h.mu.Unlock()

	for i := 0; i < h.n; i++ {
		if h.connected[str(i)] {
			for c := 0; c < len(h.commits[str(i)]); c++ {
				gotCmd := h.commits[str(i)][c].Command
				if equalsCmd(gotCmd, cmd) {
					h.t.Errorf("found %v at commits[%d][%d], expected none", cmd, i, c)
				}
			}
		}
	}
}

// SubmitToServer submits the command to serverId.
func (h *Harness) SubmitToServer(serverId string, cmdNum int) bool {
	cmd := cmdFor(cmdNum)
	return h.cluster[serverId].cm.Submit(cmd)
}

func tlog(format string, a ...interface{}) {
	format = "[TEST] " + format
	log.Printf(format, a...)
}

func sleepMs(n int) {
	time.Sleep(time.Duration(n) * time.Millisecond)
}

// collectCommits reads channel commitChans[i] and adds all received entries
// to the corresponding commits[i]. It's blocking and should be run in a
// separate goroutine. It returns when commitChans[i] is closed.
func (h *Harness) collectCommits(i int) {
	for c := range h.commitChans[str(i)] {
		h.mu.Lock()
		tlog("collectCommits(%d) got %+v", i, c)
		h.commits[str(i)] = append(h.commits[str(i)], c)
		h.mu.Unlock()
	}
}

func str(i int) string {
	return fmt.Sprintf("%d", i)
}

func equalsCmd(a CommandImpl, b CommandImpl) bool {
	if a.Command != b.Command {
		return false
	}

	if len(a.Args) != len(b.Args) {
		return false
	}

	for i := 0; i < len(a.Args); i++ {
		if a.Args[i] != b.Args[i] {
			return false
		}
	}

	return true
}

func cmdFor(n int) CommandImpl {
	c := CommandImpl{Command: "set", Args: make([]string, 0)}
	c.Args = append(c.Args, "testkey")
	c.Args = append(c.Args, str(n))

	return c
}

func cmdInt(c CommandImpl) int {
	ns := c.Args[1]
	v, err := strconv.Atoi(ns)
	if err != nil {
		log.Fatalf("can't parse %s as int", ns)
	}
	return v
}

func toInt(str string) int {
	i, err := strconv.Atoi(str)

	if err != nil {
		log.Fatalf("can't convert %s to int", str)
	}
	return i
}
