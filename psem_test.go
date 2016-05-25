package psem

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestSemaphore(t *testing.T) {
	name := "/tmp_test"
	sem, err := Open(name, Creat, 0600, 1)
	if err != nil {
		t.Error("failed create sem", err)
	} else if sem == nil {
		t.Error("bad pointer to created sem")
	}
	defer func() {
		Unlink(name)
		if sem.semt != nil {
			sem.Close()
		}
	}()
	sem2, err2 := Open(name, Creat|Excl, 0600, 1)
	if err2 == nil {
		t.Error("must be error cs Excl to same name", err2)
	} else if sem2 != nil {
		t.Error("sem have to be nil when open failed")
	}
	defer func() {
		if sem2 != nil && sem2.semt != nil {
			sem2.Close()
		}
	}()
	err = sem.TryWait()
	if err != nil {
		t.Error("failed TryWait", err)
	}
	err = sem.TryWait()
	if err == nil {
		t.Error("failed TryWait", err)
	}
	err = sem.Post()
	if err != nil {
		t.Error("failed Post", err)
	}
	err = sem.Wait()
	if err != nil {
		t.Error("failed Wait()", err)
	}
	sem.Post()
	sem.Post()
	val, err := sem.Get()
	if err != nil || val != 2 {
		t.Error("failed Get()", val, err)
	}
	err = Unlink(name)
	if err != nil {
		t.Error("failed Unlink()", err)
	}
	if _, err := os.Stat("/dev/shm/sem.tmp_test"); !os.IsNotExist(err) {
		t.Error("failed Unlink(), file exists /dev/shm/sem.tmp_test")
	}
	err = sem.Post()
	if err != nil {
		t.Error("Post dont work after unlink!", err)
	}
	sem.Wait()
	sem.Wait()
	tm := time.Now()
	err = sem.TimedWait(time.Microsecond)
	if err != nil {
		t.Error("failed timed Wait()", err)
	} else if time.Since(tm) > time.Millisecond {
		t.Error("too long TimedWait for free sem")
	}
	tm = time.Now()
	err = sem.TimedWait(time.Millisecond * 10)
	if err == nil {
		t.Error("failed timed Wait(): nil error", err)
	} else if time.Since(tm) > time.Millisecond*11 {
		t.Error("too long TimedWait for locked sem")
	} else if time.Since(tm) < time.Millisecond*9 {
		t.Error("too short TimedWait for locked sem")
	}
	err = sem.Close()
	if err != nil {
		t.Error("Failed Close()", err)
	}
}

func BenchmarkPostWait(b *testing.B) {
	sem, _ := Open("/benchsem", Creat, 0600, 0)
	Unlink("/benchsem")
	for i := 0; i < b.N; i++ {
		sem.Post()
		sem.Wait()
	}
}
func BenchmarkPostWaitParallel(b *testing.B) {
	sem, _ := Open("/benchsem", Creat, 0600, 0)
	Unlink("/benchsem")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sem.Post()
			sem.Wait()
		}
	})
}

func ExampleSemaphore() {
	sem, err := Open("/example", Creat|Excl, 0600, 1)
	Unlink("/example")
	defer sem.Close()
	if err != nil {
		fmt.Println(err)
	}
	sem.Wait() // Lock
	fmt.Println("safe print")
	sem.Post() // Unlock

	sem.Wait()
	val, _ := sem.Get()
	fmt.Println("semaphore value", val) // 0

	err = sem.TryWait()
	if err == EAgain {
		go func() {
			time.Sleep(time.Second)
			sem.Post()
		}()
		err = sem.TimedWait(time.Second * 5)
		fmt.Println("no error now:", err)
	}
	// Output:
	// safe print
	// semaphore value 0
	// no error now: <nil>
}
