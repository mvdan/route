// Copyright (c) 2019, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package route

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	if os.Getenv("ROUTE_DUMP") != "" {
		fmt.Print(`to("1.2.3.4") = `)
		fmt.Println(to("1.2.3.4"))
		fmt.Print(`Offline() = `)
		fmt.Println(Offline())
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func TestLinuxViaDocker(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skipf("skipping linux test on non-linux")
	}
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skipf("skipping docker test since it's not installed")
	}
	t.Parallel()

	const image = "alpine:3.9"
	binPath, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	dir, name := filepath.Split(binPath)
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			"NoNetwork",
			[]string{"--network=none", image},
			`
				to("1.2.3.4") =  <nil>
				Offline() = true
			`,
		},
		{
			"BridgeNetworkNoBinary",
			[]string{"--network=bridge", "--env=PATH=", image},
			`
				to("1.2.3.4") =  unsupported platform
				Offline() = false
			`,
		},
		{
			"BridgeNetwork",
			[]string{"--network=bridge", image},
			// TODO: is the bridge IP always the same? seems to be.
			`
				to("1.2.3.4") = 172.17.0.1 <nil>
				Offline() = false
			`,
		},
		{
			"BridgeNetworkGuestuser",
			[]string{"--network=bridge", "--user=guest", image},
			`
				to("1.2.3.4") = 172.17.0.1 <nil>
				Offline() = false
			`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test := test
			t.Parallel()

			args := append([]string{
				"run",
				"--volume=" + dir + ":/bindir",
				"--workdir=/bindir",
				"--entrypoint=/bindir/" + name,
				"--env=ROUTE_DUMP=true",
			}, test.args...)
			cmd := exec.Command("docker", args...)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatal(err)
			}
			want := strings.Replace(test.want, "\t", "", -1)
			want = strings.TrimPrefix(want, "\n")
			if got := string(out); got != want {
				t.Fatalf("\ngot:\n%s\nwant:\n%s", got, want)
			}
		})
	}
}
