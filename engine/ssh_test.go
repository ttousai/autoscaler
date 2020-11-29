// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package engine

import (
	"strings"
	"testing"
)

func TestSSHClient(t *testing.T) {
	client, closer, err := testSSHClient(nil)
	if err != nil {
		t.Errorf("%s\n", err)
	}
	defer closer.Close()

	out, err := client.Run("whoami")
	if err != nil {
		t.Errorf("%s\n", err)
	}
	if !strings.Contains(out, "niinai") {
		t.Errorf("username doesn't match")
	}
}

func TestPing(t *testing.T) {
	_, closer, err := testSSHClient(nil)
	if err != nil {
		t.Errorf("%s\n", err)
	}
	defer closer.Close()

	/*
		out, err := client.Ping()
		if err != nil {
			t.Errorf("%s\n", err)
		}
	*/
}
