//go:build !race
// +build !race

package main

import (
	"strings"
	"testing"
)

func TestRequestStatsTac(t *testing.T) {
	peer := StartTestPeer(4, 10, 10)
	PauseTestPeers(peer)

	if err := assertEq(4, len(PeerMap)); err != nil {
		t.Error(err)
	}

	oldParser := defaultParseOptimizer
	defaultParseOptimizer = ParseDefault
	query := strings.ReplaceAll(tacPageStatsQuery, "OutputFormat: json", "OutputFormat: wrapped_json")
	res, meta, err := peer.QueryString(query)
	if err != nil {
		t.Fatal(err)
	}
	if err = assertEq(int64(1), meta.Total); err != nil {
		t.Error(err)
	}
	if err = assertEq(int64(40), meta.RowsScanned); err != nil {
		t.Error(err)
	}
	if err = assertEq(float64(40), res[0][0]); err != nil {
		t.Error(err)
	}
	if err = assertEq(float64(32), res[0][7]); err != nil {
		t.Error(err)
	}
	if err = assertEq(float64(28), res[0][8]); err != nil {
		t.Error(err)
	}
	if err = assertEq(float64(4), res[0][9]); err != nil {
		t.Error(err)
	}

	defaultParseOptimizer = ParseOptimize
	res2, meta2, err2 := peer.QueryString(query)
	defaultParseOptimizer = ParseDefault
	if err2 != nil {
		t.Fatal(err2)
	}
	if err = assertEq(int64(1), meta2.Total); err != nil {
		t.Error(err)
	}
	if err = assertEq(int64(40), meta2.RowsScanned); err != nil {
		t.Error(err)
	}
	if err = assertEq(res, res2); err != nil {
		t.Error(err)
	}

	defaultParseOptimizer = oldParser
	if err := StopTestPeer(peer); err != nil {
		panic(err.Error())
	}
}
