// Eli Bendersky [https://eli.thegreenplace.net]
// This code is in the public domain.
package raft

import (
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
)

func TestElectionBasic(t *testing.T) {
	h := NewHarness(t, 3)
	defer h.Shutdown()

	h.CheckSingleLeader()
}

func TestElectionLeaderDisconnect(t *testing.T) {
	h := NewHarness(t, 3)
	defer h.Shutdown()

	origLeaderId, origTerm := h.CheckSingleLeader()

	h.DisconnectPeer(origLeaderId)
	sleepMs(350)

	newLeaderId, newTerm := h.CheckSingleLeader()
	if newLeaderId == origLeaderId {
		t.Errorf("want new leader to be different from orig leader")
	}
	if newTerm <= origTerm {
		t.Errorf("want newTerm <= origTerm, got %d and %d", newTerm, origTerm)
	}
}

func TestElectionLeaderAndAnotherDisconnect(t *testing.T) {
	h := NewHarness(t, 3)
	defer h.Shutdown()

	origLeaderId, _ := h.CheckSingleLeader()

	h.DisconnectPeer(origLeaderId)
	otherId := (toInt(origLeaderId) + 1) % 3
	h.DisconnectPeer(str(otherId))

	// No quorum.
	sleepMs(450)
	h.CheckNoLeader()

	// Reconnect one other server; now we'll have quorum.
	h.ReconnectPeer(str(otherId))
	h.CheckSingleLeader()
}

func TestDisconnectAllThenRestore(t *testing.T) {
	h := NewHarness(t, 3)
	defer h.Shutdown()

	sleepMs(100)
	//	Disconnect all servers from the start. There will be no leader.
	for i := 0; i < 3; i++ {
		h.DisconnectPeer(str(i))
	}
	sleepMs(450)
	h.CheckNoLeader()

	// Reconnect all servers. A leader will be found.
	for i := 0; i < 3; i++ {
		h.ReconnectPeer(str(i))
	}
	h.CheckSingleLeader()
}

func TestElectionLeaderDisconnectThenReconnect(t *testing.T) {
	h := NewHarness(t, 3)
	defer h.Shutdown()
	origLeaderId, _ := h.CheckSingleLeader()

	h.DisconnectPeer(origLeaderId)

	sleepMs(350)
	newLeaderId, newTerm := h.CheckSingleLeader()

	h.ReconnectPeer(origLeaderId)
	sleepMs(150)

	againLeaderId, againTerm := h.CheckSingleLeader()

	if newLeaderId != againLeaderId {
		t.Errorf("again leader id got %s; want %s", againLeaderId, newLeaderId)
	}
	if againTerm != newTerm {
		t.Errorf("again term got %d; want %d", againTerm, newTerm)
	}
}

func TestElectionLeaderDisconnectThenReconnect5(t *testing.T) {
	defer leaktest.CheckTimeout(t, 100*time.Millisecond)()

	h := NewHarness(t, 5)
	defer h.Shutdown()

	origLeaderId, _ := h.CheckSingleLeader()

	h.DisconnectPeer(origLeaderId)
	sleepMs(150)
	newLeaderId, newTerm := h.CheckSingleLeader()

	h.ReconnectPeer(origLeaderId)
	sleepMs(150)

	againLeaderId, againTerm := h.CheckSingleLeader()

	if newLeaderId != againLeaderId {
		t.Errorf("again leader id got %s; want %s", againLeaderId, newLeaderId)
	}
	if againTerm != newTerm {
		t.Errorf("again term got %d; want %d", againTerm, newTerm)
	}
}

func TestElectionFollowerComesBack(t *testing.T) {
	defer leaktest.CheckTimeout(t, 100*time.Millisecond)()

	h := NewHarness(t, 3)
	defer h.Shutdown()

	origLeaderId, origTerm := h.CheckSingleLeader()

	otherId := (toInt(origLeaderId) + 1) % 3
	h.DisconnectPeer(str(otherId))
	time.Sleep(650 * time.Millisecond)
	h.ReconnectPeer(str(otherId))
	sleepMs(150)

	// We can't have an assertion on the new leader id here because it depends
	// on the relative election timeouts. We can assert that the term changed,
	// however, which implies that re-election has occurred.
	_, newTerm := h.CheckSingleLeader()
	if newTerm <= origTerm {
		t.Errorf("newTerm=%d, origTerm=%d", newTerm, origTerm)
	}
}

func TestElectionDisconnectLoop(t *testing.T) {
	defer leaktest.CheckTimeout(t, 100*time.Millisecond)()

	h := NewHarness(t, 3)
	defer h.Shutdown()

	for cycle := 0; cycle < 5; cycle++ {
		leaderId, _ := h.CheckSingleLeader()

		h.DisconnectPeer(leaderId)
		otherId := (toInt(leaderId) + 1) % 3
		h.DisconnectPeer(str(otherId))
		sleepMs(310)
		h.CheckNoLeader()

		// Reconnect both.
		h.ReconnectPeer(str(otherId))
		h.ReconnectPeer(leaderId)

		// Give it time to settle
		sleepMs(150)
	}
}

func TestCommitOneCommand(t *testing.T) {
	h := NewHarness(t, 3)
	defer h.Shutdown()

	sleepMs(500)

	origLeaderId, _ := h.CheckSingleLeader()

	tlog("submitting 42 to %d", origLeaderId)
	isLeader := h.SubmitToServer(origLeaderId, 42)
	if !isLeader {
		t.Errorf("want id=%s leader, but it's not", origLeaderId)
	}

	sleepMs(500)
	h.CheckCommittedN(42, 3)
}

func TestSubmitNonLeaderFails(t *testing.T) {
	h := NewHarness(t, 3)
	defer h.Shutdown()

	origLeaderId, _ := h.CheckSingleLeader()
	sid := (toInt(origLeaderId) + 1) % 3
	tlog("submitting 42 to %d", sid)
	isLeader := h.SubmitToServer(str(sid), 42)
	if isLeader {
		t.Errorf("want id=%d !leader, but it is", sid)
	}
	sleepMs(10)
}

func TestCommitMultipleCommands(t *testing.T) {
	h := NewHarness(t, 3)
	defer h.Shutdown()

	origLeaderId, _ := h.CheckSingleLeader()

	values := []int{42, 55, 81}
	for _, v := range values {
		tlog("submitting %d to %d", v, origLeaderId)
		isLeader := h.SubmitToServer(origLeaderId, v)
		if !isLeader {
			t.Errorf("want id=%s leader, but it's not", origLeaderId)
		}
		sleepMs(100)
	}

	sleepMs(250)
	nc, i1 := h.CheckCommitted(42)
	_, i2 := h.CheckCommitted(55)
	if nc != 3 {
		t.Errorf("want nc=3, got %d", nc)
	}
	if i1 >= i2 {
		t.Errorf("want i1<i2, got i1=%d i2=%d", i1, i2)
	}

	_, i3 := h.CheckCommitted(81)
	if i2 >= i3 {
		t.Errorf("want i2<i3, got i2=%d i3=%d", i2, i3)
	}
}

func TestCommitWithDisconnectionAndRecover(t *testing.T) {
	h := NewHarness(t, 3)
	defer h.Shutdown()

	// Submit a couple of values to a fully connected cluster.
	origLeaderId, _ := h.CheckSingleLeader()
	h.SubmitToServer(origLeaderId, 5)
	h.SubmitToServer(origLeaderId, 6)

	sleepMs(250)
	h.CheckCommittedN(6, 3)

	dPeerId := (toInt(origLeaderId) + 1) % 3
	h.DisconnectPeer(str(dPeerId))
	sleepMs(250)

	// Submit a new command; it will be committed but only to two servers.
	h.SubmitToServer(origLeaderId, 7)
	sleepMs(250)
	h.CheckCommittedN(7, 2)

	// Now reconnect dPeerId and wait a bit; it should find the new command too.
	h.ReconnectPeer(str(dPeerId))
	sleepMs(250)
	h.CheckSingleLeader()

	sleepMs(150)
	h.CheckCommittedN(7, 3)
}

func TestNoCommitWithNoQuorum(t *testing.T) {
	h := NewHarness(t, 3)
	defer h.Shutdown()

	// Submit a couple of values to a fully connected cluster.
	origLeaderId, origTerm := h.CheckSingleLeader()
	h.SubmitToServer(origLeaderId, 5)
	h.SubmitToServer(origLeaderId, 6)

	sleepMs(250)
	h.CheckCommittedN(6, 3)

	// Disconnect both followers.
	dPeer1 := (toInt(origLeaderId) + 1) % 3
	dPeer2 := (toInt(origLeaderId) + 2) % 3
	h.DisconnectPeer(str(dPeer1))
	h.DisconnectPeer(str(dPeer2))
	sleepMs(250)

	h.SubmitToServer(origLeaderId, 8)
	sleepMs(250)
	h.CheckNotCommitted(8)

	// Reconnect both other servers, we'll have quorum now.
	h.ReconnectPeer(str(dPeer1))
	h.ReconnectPeer(str(dPeer2))
	sleepMs(600)

	// 8 is still not committed because the term has changed.
	h.CheckNotCommitted(8)

	// A new leader will be elected. It could be a different leader, even though
	// the original's log is longer, because the two reconnected peers can elect
	// each other.
	newLeaderId, againTerm := h.CheckSingleLeader()
	if origTerm == againTerm {
		t.Errorf("got origTerm==againTerm==%d; want them different", origTerm)
	}

	// But new values will be committed for sure...
	h.SubmitToServer(newLeaderId, 9)
	h.SubmitToServer(newLeaderId, 10)
	h.SubmitToServer(newLeaderId, 11)
	sleepMs(350)

	for _, v := range []int{9, 10, 11} {
		h.CheckCommittedN(v, 3)
	}
}

func TestDisconnectLeaderBriefly(t *testing.T) {
	h := NewHarness(t, 3)
	defer h.Shutdown()

	// Submit a couple of values to a fully connected cluster.
	origLeaderId, _ := h.CheckSingleLeader()
	h.SubmitToServer(origLeaderId, 5)
	h.SubmitToServer(origLeaderId, 6)
	sleepMs(250)
	h.CheckCommittedN(6, 3)

	// Disconnect leader for a short time (less than election timeout in peers).
	h.DisconnectPeer(origLeaderId)
	sleepMs(90)
	h.ReconnectPeer(origLeaderId)
	sleepMs(200)

	h.SubmitToServer(origLeaderId, 7)
	sleepMs(250)
	h.CheckCommittedN(7, 3)
}

func TestCommitsWithLeaderDisconnects(t *testing.T) {
	h := NewHarness(t, 5)
	defer h.Shutdown()

	// Submit a couple of values to a fully connected cluster.
	origLeaderId, _ := h.CheckSingleLeader()
	h.SubmitToServer(origLeaderId, 5)
	h.SubmitToServer(origLeaderId, 6)

	sleepMs(250)
	h.CheckCommittedN(6, 5)

	// Leader disconnected...
	h.DisconnectPeer(origLeaderId)
	sleepMs(10)

	// Submit 7 to original leader, even though it's disconnected.
	h.SubmitToServer(origLeaderId, 7)

	sleepMs(250)
	h.CheckNotCommitted(7)

	newLeaderId, _ := h.CheckSingleLeader()

	// Submit 8 to new leader.
	h.SubmitToServer(newLeaderId, 8)
	sleepMs(250)
	h.CheckCommittedN(8, 4)

	// Reconnect old leader and let it settle. The old leader shouldn't be the one
	// winning.
	h.ReconnectPeer(origLeaderId)
	sleepMs(600)

	finalLeaderId, _ := h.CheckSingleLeader()
	if finalLeaderId == origLeaderId {
		t.Errorf("got finalLeaderId==origLeaderId==%s, want them different", finalLeaderId)
	}

	// Submit 9 and check it's fully committed.
	h.SubmitToServer(newLeaderId, 9)
	sleepMs(250)
	h.CheckCommittedN(9, 5)
	h.CheckCommittedN(8, 5)

	// But 7 is not committed...
	h.CheckNotCommitted(7)
}

func TestCrashFollower(t *testing.T) {
	h := NewHarness(t, 3)
	defer h.Shutdown()

	origLeaderId, _ := h.CheckSingleLeader()
	h.SubmitToServer(origLeaderId, 5)

	sleepMs(350)
	h.CheckCommittedN(5, 3)

	h.CrashPeer(str((toInt(origLeaderId) + 1) % 3))
	sleepMs(350)
	h.CheckCommittedN(5, 2)
}

func TestCrashThenRestartFollower(t *testing.T) {
	h := NewHarness(t, 3)
	defer h.Shutdown()

	origLeaderId, _ := h.CheckSingleLeader()
	h.SubmitToServer(origLeaderId, 5)
	h.SubmitToServer(origLeaderId, 6)
	h.SubmitToServer(origLeaderId, 7)

	vals := []int{5, 6, 7}

	sleepMs(350)
	for _, v := range vals {
		h.CheckCommittedN(v, 3)
	}

	peerId := toInt(origLeaderId) + 1%3
	h.CrashPeer(str(peerId))

	sleepMs(350)
	for _, v := range vals {
		h.CheckCommittedN(v, 2)
	}

	// Restart the crashed follower and give it some time to come up-to-date.
	h.RestartPeer(str((toInt(origLeaderId) + 1) % 3))
	sleepMs(1000)
	for _, v := range vals {
		h.CheckCommittedN(v, 3)
	}
}

func TestCrashThenRestartLeader(t *testing.T) {
	h := NewHarness(t, 3)
	defer h.Shutdown()

	origLeaderId, _ := h.CheckSingleLeader()
	h.SubmitToServer(origLeaderId, 5)
	h.SubmitToServer(origLeaderId, 6)
	h.SubmitToServer(origLeaderId, 7)

	vals := []int{5, 6, 7}

	sleepMs(350)
	for _, v := range vals {
		h.CheckCommittedN(v, 3)
	}

	h.CrashPeer(origLeaderId)
	sleepMs(350)
	for _, v := range vals {
		h.CheckCommittedN(v, 2)
	}

	h.RestartPeer(origLeaderId)
	sleepMs(550)
	for _, v := range vals {
		h.CheckCommittedN(v, 3)
	}
}

func TestCrashThenRestartAll(t *testing.T) {
	h := NewHarness(t, 3)
	defer h.Shutdown()

	origLeaderId, _ := h.CheckSingleLeader()
	h.SubmitToServer(origLeaderId, 5)
	h.SubmitToServer(origLeaderId, 6)
	h.SubmitToServer(origLeaderId, 7)

	vals := []int{5, 6, 7}

	sleepMs(350)
	for _, v := range vals {
		h.CheckCommittedN(v, 3)
	}

	for i := 0; i < 3; i++ {
		h.CrashPeer(str((toInt(origLeaderId) + i) % 3))
	}

	sleepMs(350)

	for i := 0; i < 3; i++ {
		h.RestartPeer(str((toInt(origLeaderId) + i) % 3))
	}

	sleepMs(150)
	newLeaderId, _ := h.CheckSingleLeader()

	h.SubmitToServer(newLeaderId, 8)
	sleepMs(250)

	vals = []int{5, 6, 7, 8}
	for _, v := range vals {
		h.CheckCommittedN(v, 3)
	}
}

func TestReplaceMultipleLogEntries(t *testing.T) {
	h := NewHarness(t, 3)
	defer h.Shutdown()

	// Submit a couple of values to a fully connected cluster.
	origLeaderId, _ := h.CheckSingleLeader()
	h.SubmitToServer(origLeaderId, 5)
	h.SubmitToServer(origLeaderId, 6)

	sleepMs(250)
	h.CheckCommittedN(6, 3)

	// Leader disconnected...
	h.DisconnectPeer(origLeaderId)
	sleepMs(10)

	// Submit a few entries to the original leader; it's disconnected, so they
	// won't be replicated.
	h.SubmitToServer(origLeaderId, 21)
	sleepMs(5)
	h.SubmitToServer(origLeaderId, 22)
	sleepMs(5)
	h.SubmitToServer(origLeaderId, 23)
	sleepMs(5)
	h.SubmitToServer(origLeaderId, 24)
	sleepMs(5)

	newLeaderId, _ := h.CheckSingleLeader()

	// Submit entries to new leader -- these will be replicated.
	h.SubmitToServer(newLeaderId, 8)
	sleepMs(5)
	h.SubmitToServer(newLeaderId, 9)
	sleepMs(5)
	h.SubmitToServer(newLeaderId, 10)
	sleepMs(250)
	h.CheckNotCommitted(21)
	h.CheckCommittedN(10, 2)

	// Crash/restart new leader to reset its nextIndex, to ensure that the new
	// leader of the cluster (could be the third server after elections) tries
	// to replace the original's servers unreplicated entries from the very end.
	h.CrashPeer(newLeaderId)
	sleepMs(60)
	h.RestartPeer(newLeaderId)

	sleepMs(200)

	finalLeaderId, _ := h.CheckSingleLeader()

	h.ReconnectPeer(origLeaderId)
	sleepMs(400)

	sleepMs(100)

	// Submit another entry; this is because leaders won't commit entries from
	// previous terms (paper 5.4.2) so the 8,9,10 may not be committed everywhere
	// after the restart before a new command comes it.
	h.SubmitToServer(finalLeaderId, 11)
	sleepMs(250)

	// At this point, 11 and 10 should be replicated everywhere; 21 won't be.
	h.CheckNotCommitted(21)

	h.CheckCommittedN(11, 3)
	h.CheckCommittedN(10, 3)
}
