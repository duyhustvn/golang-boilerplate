#ifndef QUEUE_H_
#define QUEUE_H_

#include "../common.h"


typedef struct Node_ {
    void *data;
    struct Node_ *next;
} Node;

Node *create_node(void *data, char *errstr, size_t err_msg_size);
void free_node(Node *node);

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

Queue *init_queue(char *errstr, size_t err_msg_size);
bool enqueue(Queue *q, void *data, char *errstr);
Node *dequeue(Queue *q, char *errstr, size_t err_msg_size);
bool is_empty(Queue *q);
void free_queue(Queue *q);

#endif // QUEUE_H_
