#ifndef QUEUE_H_
#define QUEUE_H_

#include "../common.h"


typedef struct Node_ {
    void *data;
    struct Node_ *next;
} Node;

Node *create_node(void *data, char *errstr);

/*
 * enqueue into rear 
 * dequeue from front
 * rear <- node <- node <- front
 * */
typedef struct Queue_ {
    Node *front;
    Node *rear;

    uint32_t length;

    pthread_mutex_t mutex;
} Queue;

Queue *init_queue(char *errstr);
bool enqueue(Queue *q, void *data, char *errstr);
void *dequeue(Queue *q, char *errstr);
bool is_empty(Queue *q);

#endif // QUEUE_H_
