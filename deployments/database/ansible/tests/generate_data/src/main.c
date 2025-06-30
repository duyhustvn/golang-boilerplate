#include "common.h"

#include "http/http_client.h"
#include "api_caller/signup_api.h"
#include "queue/queue.h"

#define NUM_THREADS 8


int main() {
    // Initialize client (call once at program start)
    http_client_init();

    char errstr[1024];
    Queue *q = init_queue(errstr);
    if (!q) {
        printf("%s\n", errstr);
    }

    int *v;
    for (int i = 0; i < 1000; i++) {
        v = malloc(sizeof(int));
        *v = i;
        bool ok = enqueue(q, v, errstr);
        if (!ok) {
            printf("%s\n", errstr);
            continue;
        }
    }

    printf("queue length: %d\n", q->length);

    if (is_empty(q)) {
        return 0;
    }

    
    int rc;
    pthread_t threads[NUM_THREADS];
    sign_up_thread_vars thread_vars[NUM_THREADS];

    for (int i = 0; i < NUM_THREADS; i++) {
        thread_vars[i].threadId = i;
        thread_vars[i].q = q;

        rc = pthread_create(&threads[i], NULL, call_signup_api_thread, (void *)&thread_vars[i]);
        if (rc) {
            printf("Error in creating new thread");
            return 0;
        }
    }

    /* Wait for all threads to finish */
    void *status;
    for (int i = 0; i < NUM_THREADS; i++) {
        rc = pthread_join(threads[i], &status);
        if (rc) {
            printf("Error in joinning the thread");
            return 0;
        }
    }
     
    // Cleanup client (call once at program end)
    http_client_cleanup();
    return 0;
}
