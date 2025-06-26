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

