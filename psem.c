#include "psem.h"


#include <sys/stat.h>        /* For mode constants */
void Open( sem_t ** sem, const char * name,int oflag, mode_t mode, unsigned int value){
  *sem = sem_open(name,oflag,mode,value);
}

void  Wait(sem_t *sem){
  sem_wait(sem);
}

void  TryWait(sem_t *sem){
  sem_trywait(sem);
}

void TimedWait(sem_t *sem,long int seconds,long int nanoseconds, timespec_t* ts){
  clock_gettime(CLOCK_REALTIME, ts);
  ts->tv_sec += seconds;
  ts->tv_nsec += nanoseconds;
  sem_timedwait(sem,ts);
}

void Get (sem_t * sem,unsigned int * val){
  sem_getvalue(sem,val);
}

void Post (sem_t * sem){
  sem_post(sem);
}

void Close (sem_t * sem){
  sem_close(sem);
}

void Unlink (const char* name) {
  sem_unlink(name);
}