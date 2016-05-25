# psem
--
    import "github.com/avezila/psem"

Package psem is a warp for POSIX named semaphore functions in C.

http://linux.die.net/man/7/sem_overview

POSIX semaphores allow processes and threads to synchronize their actions. A
semaphore is an integer whose value is never allowed to fall below zero. Two
operations can be performed on semaphores: increment the semaphore value by one
Post() and decrement the semaphore value by one Wait(). If the value of a
semaphore is currently zero, then a Wait() operation will block until the value
becomes greater than zero.

Benchmark on my laptop thinkpad w520: BenchmarkPostWait-8 3000000 416 ns/op
BenchmarkPostWaitParallel-8 10000000 175 ns/op

## Usage

```go
const (
	// Creat - then the semaphore is created if it does not already exist.
	Creat = int(C.O_CREAT)

	// Excl - if both Creat and Excl are specified in oflag,
	// then an error is returned if a semaphore with the given name already exists.
	Excl = int(C.O_EXCL)

	// EAcces - The semaphore exists, but the caller does not have permission to open it.
	// or - Unlink(): The caller does not have permission to unlink this semaphore.
	EAcces = syscall.Errno(C.EACCES)
	// EExist - Both O_CREAT and O_EXCL were specified in oflag,
	// but a semaphore with this name already exists.
	EExist = syscall.Errno(C.EEXIST)
	// EInVal - value was greater than SEM_VALUE_MAX.
	// or - name consists of just "/", followed by no other characters.
	// or - sem is not a valid semaphore.
	// or - TimedWait(): The value of timout is less than 0, or greater than or equal to 1000 million.
	EInVal = syscall.Errno(C.EINVAL)
	// EMFile - The process already has the maximum number of files and open.
	EMFile = syscall.Errno(C.EMFILE)
	// ENameTooLong - name was too long.
	ENameTooLong = syscall.Errno(C.ENAMETOOLONG)
	// ENFile - The system limit on the total number of open files has been reached.
	ENFile = syscall.Errno(C.ENFILE)
	// ENoEnt - Open(): The Creat flag was not specified and no semaphore with this name exists;
	// or, Creat was specified, but name wasn't well formed.
	// or - Unlink(): There is no semaphore with the given name.
	ENoEnt = syscall.Errno(C.ENOENT)
	// ENoMem - Insufficient memory.
	ENoMem = syscall.Errno(C.ENOMEM)
	// EOverflow - The maximum allowable value for a semaphore would be exceeded.
	EOverflow = syscall.Errno(C.EOVERFLOW)
	// EIntr - The call was interrupted by a signal handler
	EIntr = syscall.Errno(C.EINTR)
	// EAgain - TryWait(): The operation could not be performed without blocking
	// (i.e., the semaphore currently has the value zero).
	EAgain = syscall.Errno(C.EAGAIN)
	// ETimedOut - TimedWait(): The call timed out before the semaphore could be locked.
	ETimedOut = syscall.Errno(C.ETIMEDOUT)
)
```

#### func  Unlink

```go
func Unlink(name string) error
```
Unlink - removes the named semaphore referred to by name. The semaphore name is
removed immediately. The semaphore is destroyed once all other processes that
have the semaphore open close it.

#### type Sem

```go
type Sem struct {
}
```

Sem - struct for save pointer to the semaphore

#### func  Open

```go
func Open(name string, flag int, mode int, value uint) (*Sem, error)
```
Open - creates a new POSIX semaphore or opens an existing semaphore. The
semaphore is identified by name. If Creat is specified, and a semaphore with the
given name already exists, then mode and value are ignored. flag: Creat|Excl
mode: for example 0600 value: initial value

#### func (*Sem) Close

```go
func (sem *Sem) Close() error
```
Close - closes the named semaphore, allowing any resources that the system has
allocated to the calling process for this semaphore to be freed.

#### func (*Sem) Get

```go
func (sem *Sem) Get() (uint, error)
```
Get - return current value of the semaphore. If one or more processes or threads
are blocked waiting to lock the semaphore with Wait() POSIX.1-2001 permits two
possibilities for the value returned in sval: either 0 is returned; or a
negative number whose absolute value is the count of the number of processes and
threads currently blocked in Wait()

#### func (*Sem) Post

```go
func (sem *Sem) Post() error
```
Post - increments (unlocks) the semaphore. If the semaphore's value consequently
becomes greater than zero, then another process or thread blocked in a psem.Wait
call will be woken up and proceed to lock the semaphore.

#### func (*Sem) TimedWait

```go
func (sem *Sem) TimedWait(timeout time.Duration) error
```
TimedWait - is the same as Wait(), except that timeout specifies a limit on the
amount of time that the call should block if the decrement cannot be immediately
performed. If the timeout has already expired by the time of the call, and the
semaphore could not be locked immediately, then TimedWait() fails with a timeout
error. If the operation can be performed immediately, then TimedWait() never
fails with a timeout error, regardless of the value of timeout. Furthermore, the
validity of timeout is not checked in this case.

#### func (*Sem) TryWait

```go
func (sem *Sem) TryWait() error
```
TryWait - is the same as Wait(), except that if the decrement cannot be
immediately performed, then call returns an error instead of blocking.

#### func (*Sem) Wait

```go
func (sem *Sem) Wait() error
```
Wait - decrements (locks) the semaphore. If the semaphore's value is greater
than zero, then the decrement proceeds, and the function returns, immediately.
If the semaphore currently has the value zero, then the call blocks until either
it becomes possible to perform the decrement (i.e., the semaphore value rises
above zero), or a signal handler interrupts the call.
