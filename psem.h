#if !defined PSEM_H
#define PSEM_H	1

#include <semaphore.h>
#include <fcntl.h>           /* For O_* constants */
#include <time.h>
#include <errno.h>


typedef struct timespec timespec_t;

void  Open(sem_t **sem, const char* name, int oflag, mode_t mode, unsigned int value);
void  Post(sem_t *sem);
void  Wait(sem_t *sem);
void  TryWait(sem_t *sem);
void  TimedWait(sem_t *sem, long int seconds, long int nanoseconds, timespec_t * ts);
void  Get(sem_t *sem, unsigned int *val);
void  Close(sem_t *sem);
void  Unlink(const char* name);

#endif