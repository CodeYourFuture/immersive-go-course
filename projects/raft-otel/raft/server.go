// Server container for a Raft Consensus Module. Exposes Raft to the network
// and enables RPCs between Raft peers.
//
// Based on code by Eli Bendersky [https://eli.thegreenplace.net], modified somewhat by Laura Nolan.
// This code is in the public domain.
package raft

import (
	"context"
	"fmt"
	"log"
	"net"
	"raft/raft_proto"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// Server wraps a raft.ConsensusModule along with a rpc.Server that exposes its
// methods as gRPC endpoints. It also manages the peers of the Raft server. The
// main goal of this type is to simplify the code of raft.Server for
// presentation purposes. raft.ConsensusModule has a *Server to do its peer
// communication and doesn't have to worry about the specifics of running an
// RPC server.
type Server struct {
	mu sync.Mutex

	serverId string
	ip       string

	cm      *ConsensusModule
	storage Storage

	listener   net.Listener
	listenPort int
	grpcServer *grpc.Server

	commitChan  chan CommitEntry
	peerClients map[string]raft_proto.RaftServiceClient

	ready <-chan interface{}
	quit  chan interface{}

	raft_proto.UnimplementedRaftServiceServer
	raft_proto.UnimplementedRaftKVServiceServer

	fsm *KV
}

type KV struct {
	vals map[string]string
}

func NewKV() *KV {
	return &KV{vals: make(map[string]string)}
}

func NewServer(serverId string, ip string, storage Storage, ready <-chan interface{}, commitChan chan CommitEntry, listenPort int) *Server {
	s := new(Server)
	s.serverId = serverId
	s.ip = ip
	s.peerClients = make(map[string]raft_proto.RaftServiceClient)
	s.storage = storage
	s.ready = ready
	s.commitChan = commitChan
	s.quit = make(chan interface{})
	s.listenPort = listenPort
	return s
}

// If fsm is set then commitChan is read by the FSM, otherwise commitChan can be read by tests
// Bit icky. Oh well.
func (s *Server) Serve(fsm *KV) {
	s.mu.Lock()
	s.cm = NewConsensusModule(s.serverId, s, s.storage, s.ready, s.commitChan)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.ip, s.listenPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s.listener = lis

	gs := grpc.NewServer()
	s.grpcServer = gs

	raft_proto.RegisterRaftServiceServer(gs, s)
	raft_proto.RegisterRaftKVServiceServer(gs, s)

	go gs.Serve(s.listener)
	if fsm != nil {
		s.fsm = fsm
		// TODO should close this goroutine on shutdown
		go s.fsm.readCommits(s.commitChan)
	}

	log.Printf("[%v] listening at %s", s.serverId, s.listener.Addr())

	s.mu.Unlock()
}

// DisconnectAll closes all the client connections to peers for this server.
func (s *Server) DisconnectAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id := range s.peerClients {
		if s.peerClients[id] != nil {
			s.peerClients[id] = nil
		}
	}
}

// Shutdown closes the server and waits for it to shut down properly.
func (s *Server) Shutdown() {
	s.grpcServer.GracefulStop()
	s.cm.Stop()
	close(s.quit)
	s.listener.Close()
}

func (s *Server) GetListenAddr() net.Addr {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.listener.Addr()
}

func (s *Server) ConnectToPeer(peerId string, peerAddr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.peerClients[peerId] == nil {
		conn, err := grpc.Dial(peerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}
		client := raft_proto.NewRaftServiceClient(conn)
		s.peerClients[peerId] = client
	}
	s.cm.AddPeerID(peerId)
	return nil
}

// DisconnectPeer disconnects this server from the peer identified by peerId.
func (s *Server) DisconnectPeer(peerId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.peerClients[peerId] != nil {
		s.peerClients[peerId] = nil
		return nil
	}
	return nil
}

func (s *Server) CallRequestVote(id string, args RequestVoteArgs, reply *RequestVoteReply) error {
	s.mu.Lock()
	peer := s.peerClients[id]
	s.mu.Unlock()

	// If this is called after shutdown (where client.Close is called), it will
	// return an error.
	if peer == nil {
		return fmt.Errorf("call client %s after it's closed", id)
	} else {
		req := raft_proto.RequestVoteRequest{
			Term:         int64(args.Term),
			CandidateId:  args.CandidateId,
			LastLogIndex: int64(args.LastLogIndex),
			LastLogTerm:  int64(args.LastLogTerm),
		}
		resp, err := peer.RequestVote(context.TODO(), &req)
		if err != nil {
			return err
		}

		reply.Term = int(resp.GetTerm())
		reply.VoteGranted = resp.GetVoteGranted()
	}
	return nil
}

func (s *Server) RequestVote(ctx context.Context, req *raft_proto.RequestVoteRequest) (*raft_proto.RequestVoteResponse, error) {
	fmt.Printf("[%s] received RequestVote %+v\n", s.serverId, req)

	rva := RequestVoteArgs{
		Term:         int(req.GetTerm()),
		CandidateId:  req.GetCandidateId(),
		LastLogIndex: int(req.GetLastLogIndex()),
		LastLogTerm:  int(req.GetLastLogTerm()),
	}

	rvr := RequestVoteReply{}
	err := s.cm.RequestVote(rva, &rvr)
	if err != nil {
		return nil, err
	}

	resp := raft_proto.RequestVoteResponse{
		Term:        int64(rvr.Term),
		VoteGranted: rvr.VoteGranted,
	}
	return &resp, nil
}

func (s *Server) CallAppendEntries(id string, args AppendEntriesArgs, reply *AppendEntriesReply) error {
	s.mu.Lock()
	peer := s.peerClients[id]
	s.mu.Unlock()

	// If this is called after shutdown (where client.Close is called), it will
	// return an error.
	if peer == nil {
		return fmt.Errorf("call client %s after it's closed", id)
	} else {
		req := raft_proto.AppendEntriesRequest{
			Term:         int64(args.Term),
			Leader:       args.LeaderId,
			PrevLogIndex: int64(args.PrevLogIndex),
			PrevLogTerm:  int64(args.PrevLogTerm),
			LeaderCommit: int64(args.LeaderCommit),
			Entries:      make([]*raft_proto.LogEntry, 0),
		}

		for _, e := range args.Entries {
			en := raft_proto.LogEntry{
				Term:    int64(e.Term),
				Command: &raft_proto.Command{Command: e.Command.Command, Args: e.Command.Args},
			}
			req.Entries = append(req.Entries, &en)
		}

		resp, err := peer.AppendEntries(context.TODO(), &req)
		if err != nil {
			return err
		}

		reply.ConflictIndex = int(resp.GetConflictIndex())
		reply.ConflictTerm = int(resp.GetConflictTerm())
		reply.Success = resp.GetSuccess()
		reply.Term = int(resp.GetTerm())
	}
	return nil
}

func (s *Server) AppendEntries(ctx context.Context, req *raft_proto.AppendEntriesRequest) (*raft_proto.AppendEntriesResponse, error) {
	aea := AppendEntriesArgs{
		Term:         int(req.GetTerm()),
		LeaderId:     req.GetLeader(),
		PrevLogIndex: int(req.GetPrevLogIndex()),
		PrevLogTerm:  int(req.GetPrevLogTerm()),
		Entries:      make([]LogEntry, 0),
		LeaderCommit: int(req.GetLeaderCommit()),
	}
	for _, e := range req.Entries {
		en := LogEntry{
			Term:    int(e.Term),
			Command: CommandImpl{Command: e.Command.Command, Args: e.Command.Args},
		}
		aea.Entries = append(aea.Entries, en)
	}

	aer := AppendEntriesReply{}
	err := s.cm.AppendEntries(aea, &aer)
	if err != nil {
		return nil, err
	}

	resp := raft_proto.AppendEntriesResponse{
		Term:          int64(aer.Term),
		Success:       aer.Success,
		ConflictIndex: int64(aer.ConflictIndex),
		ConflictTerm:  int64(aer.ConflictTerm),
	}

	return &resp, nil
}

// TODO proxy to leader if not leader
func (s *Server) Set(ctx context.Context, req *raft_proto.SetRequest) (*raft_proto.SetResponse, error) {
	cmd := CommandImpl{
		Command: "set",
		Args:    []string{req.Keyname, req.Value},
	}

	res := s.cm.Submit(cmd)
	if !res {
		return nil, status.Error(codes.Unavailable, "not the leader")
	}

	return &raft_proto.SetResponse{}, nil
}

func (s *Server) Get(ctx context.Context, req *raft_proto.GetRequest) (*raft_proto.GetResponse, error) {
	// TODO allow gets from non-leader if the query specified
	if s.cm.state != Leader {
		return &raft_proto.GetResponse{}, status.Error(codes.Unavailable, "not the leader")
	}

	res := s.fsm.get(req.Keyname)
	return &raft_proto.GetResponse{Value: res}, nil
}

func (kv *KV) readCommits(ch chan CommitEntry) {
	for {
		entry := <-ch
		if entry.Command.Command == "set" {
			if len(entry.Command.Args) != 2 {
				log.Printf("Can't parse this set command %+v", entry.Command)
			}
			kn := entry.Command.Args[0]
			val := entry.Command.Args[1]
			kv.set(kn, val)
		}
	}
}

func (kv *KV) set(k string, v string) {
	kv.vals[k] = v
}

func (kv *KV) get(k string) string {
	return kv.vals[k]
}
