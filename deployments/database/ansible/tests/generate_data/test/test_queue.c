#include "unity.h"
#include "queue.h"

void setUp(void) {}
void tearDown(void) {}

void test_init_queue(void) {
    char errstr[1024];
    Queue *q = init_queue(errstr); 

    TEST_ASSERT_NOT_NULL(q);
    TEST_ASSERT_EQUAL_INT(q->length, 0);
    TEST_ASSERT_NULL(q->front);
    TEST_ASSERT_NULL(q->rear);

    free_queue(q);
}

void test_enqueue(void) {
    char errstr[1024];
    Queue *q = init_queue(errstr);

    TEST_ASSERT_NOT_NULL(q);
    TEST_ASSERT_EQUAL_INT(q->length, 0);
    TEST_ASSERT_NULL(q->front);
    TEST_ASSERT_NULL(q->rear);

    int *i;
    bool result;
    int *data;

    i = malloc(sizeof(int));;
    *i = 1;
    result = enqueue(q, i, errstr);
    TEST_ASSERT_EQUAL_INT(q->length, 1);
    data = (int *)(q->front->data);
    TEST_ASSERT_EQUAL_INT(*data, 1);
    data = (int *)(q->rear->data);
    TEST_ASSERT_EQUAL_INT(*data, 1);

    i = malloc(sizeof(int));;
    *i = 2;
    result = enqueue(q, i, errstr);
    TEST_ASSERT_EQUAL_INT(q->length, 2);
    data = (int *)(q->front->data);
    TEST_ASSERT_EQUAL_INT(*data, 1);
    data = (int *)(q->rear->data);
    TEST_ASSERT_EQUAL_INT(*data, 2);

    i = malloc(sizeof(int));;
    *i = 3;
    result = enqueue(q, i, errstr);
    TEST_ASSERT_EQUAL_INT(q->length, 3);
    data = (int *)(q->front->data);
    TEST_ASSERT_EQUAL_INT(*data, 1);
    data = (int *)(q->rear->data);
    TEST_ASSERT_EQUAL_INT(*data, 3);
    data = (int *)(q->front->next->data);
    TEST_ASSERT_EQUAL_INT(*data, 2);
}

void test_dequeue(void) {
    char errstr[1024];
    Queue *q = init_queue(errstr);

    TEST_ASSERT_NOT_NULL(q);
    TEST_ASSERT_EQUAL_INT(q->length, 0);


    int *i;
    bool result;
    int *data;
    Node *node;

    i = malloc(sizeof(int));;
    *i = 1;
    result = enqueue(q, i, errstr);
    TEST_ASSERT_EQUAL_INT(q->length, 1);

    i = malloc(sizeof(int));;
    *i = 2;
    result = enqueue(q, i, errstr);
    TEST_ASSERT_EQUAL_INT(q->length, 2);

    i = malloc(sizeof(int));;
    *i = 3;
    result = enqueue(q, i, errstr);
    TEST_ASSERT_EQUAL_INT(q->length, 3);

    node = dequeue(q, errstr);
    TEST_ASSERT_EQUAL_INT(q->length, 2);
    data = (int *) node->data;
    TEST_ASSERT_EQUAL_INT(*data, 1);

    node = dequeue(q, errstr);
    data = (int *) node->data;
    TEST_ASSERT_EQUAL_INT(*data, 2);

    node = dequeue(q, errstr);
    data = (int *) node->data;
    TEST_ASSERT_EQUAL_INT(*data, 3);
}
