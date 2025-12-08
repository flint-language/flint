#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>

extern void print(const char *s)
{
    fputs(s, stdout);
}

extern void assert(bool cond)
{
    if (!cond)
    {
        fprintf(stderr, "Assertion failed\n");
        fflush(stderr);
        exit(1);
    }
}