//Robert Fudge Q1 interview question file
//Minimal hello world program

#include <unistd.h>

//Have to use syscalls, no standard c library
int main() {
    write(1, "Hello World!\r\n", 15);  //1 = stdout,
    //15 = Number of bytes to write from buffer

    while(1);  //Prevent kernel panic (keep process alive)
    return 0;
}