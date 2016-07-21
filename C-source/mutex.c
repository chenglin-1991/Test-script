#include <stdio.h>
#include <unistd.h>
#include <pthread.h>
#include <sys/time.h>

int flag = 1;
pthread_t thread;
pthread_cond_t cond;
pthread_mutex_t mutex;

void *pthread_fun(void *arg)
{
	struct timeval now;
	struct timespec outtime;
	pthread_mutex_lock(&mutex);
	while (flag)
	{
		printf("enter pthread_fun\n");
		gettimeofday(&now,NULL);
		outtime.tv_sec = now.tv_sec + 5;
		outtime.tv_nsec = now.tv_usec * 1000;
		/*pthread_cond_wait(&cond,&mutex);*/
		pthread_cond_timedwait(&cond,&mutex,&outtime);
		printf("==============\n");
	}
	pthread_mutex_unlock(&mutex);
	printf("pthread_fun exit!\n");

	return NULL;
}

int main(int argc, const char *argv[])
{
	char c = 0;
	pthread_mutex_init(&mutex,NULL);
	pthread_cond_init(&cond,NULL);

	pthread_create(&thread,NULL,pthread_fun,NULL);

	while((c = getchar()) != 'q');
	flag = 0;
	pthread_mutex_lock(&mutex);
	pthread_cond_signal(&cond);
	pthread_mutex_unlock(&mutex);
	pthread_join(thread,NULL);

	return 0;
}
