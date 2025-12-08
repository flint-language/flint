#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <stdint.h>

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

extern char *to_string(int64_t i)
{
	static char buffer[32];
	sprintf(buffer, "%lld", (long long)i);
	return buffer;
}