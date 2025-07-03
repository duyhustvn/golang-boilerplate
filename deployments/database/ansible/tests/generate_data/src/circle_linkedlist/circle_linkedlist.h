#ifndef CIRCLE_LINKEDLIST_H_
#define CIRCLE_LINKEDLIST_H_

#include "../common.h"

typedef struct CllNode_ {
    void *data;
    struct CllNode_ *next;
} CllNode;

CllNode *create_cll_node(void *data, char *errstr, size_t errstr_size);
void free_cll_node(CllNode *node);

typedef struct CircleLinkedList_ {
    CllNode *head;
    int length;

    pthread_mutex_t mutex;
} CircleLinkedList;

CircleLinkedList *init_cll(char *errstr, size_t errstr_size);
bool insert_cll(CircleLinkedList *list, void *data, char *errstr, size_t errstr_size);
void free_cll(CircleLinkedList *list);

#endif // CIRCLE_LINKEDLIST_H_
