package agent

import (
	"sync"
	"testing"

	"github.com/satori/go.uuid"
)

func TestAgent_Cgroup(t *testing.T) {
	ag := agent{agentId: uuid.NewV4()}
	if err := ag.installCgroup(); err != nil {
		t.Fatal(err)
	}
	t.Logf("GetInitStub --> %v", ag.cg.GetInitStub())
	if err := ag.updateCgroup(0, 0); err != nil {
		t.Fatal(err)
	}
	if err := ag.unInstallCgroup(); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkAgent_Cgroup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ag := agent{agentId: uuid.NewV4()}
		if err := ag.installCgroup(); err != nil {
			b.Fatal(err)
		}
		if err := ag.updateCgroup(0, 0); err != nil {
			b.Fatal(err)
		}
		if err := ag.unInstallCgroup(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAgent_Cgroup_Goroutine(b *testing.B) {
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			ag := agent{agentId: uuid.NewV4()}
			if err := ag.installCgroup(); err != nil {
				b.Fatal(err)
			}
			if err := ag.updateCgroup(0, 0); err != nil {
				b.Fatal(err)
			}
			if err := ag.unInstallCgroup(); err != nil {
				b.Fatal(err)
			}
		}()
	}
	wg.Wait()
}
