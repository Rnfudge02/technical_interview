# Technical Interview Questions
This is my submission to the technical interview component of an open-source software company.

The code can be found on the Github repository [here](https://github.com/Rnfudge02/technical_interview.git)

# Use
This code is licensed under the Apache 2.0 license, feel free to use these snippets as a building point for larger programs.

Feedback on design improvements would be greatly appreciated :)

# Question 1
The challenge was to create a minimal linux image that has an ext4 filesystem and prints "Hello World!" to the screen.

The program should be bootable via QEMU.

## Use cases
Building containerized images with linux kernel and FS.

## Methodology
1. Create a script function that compiles the linux kernel with ext4 support and serial console for terminal output.
2. Create a script function that properly compiles the hello_world program
3. Create a script function which creates the file system in the project directory, mounts it, makes nessecary changes, then unmounts.
4. Create a script function that will launch the image using qemu-system-x86_64.
5. Create a minimal terminal UI for interacting with the project

## To use
1. Run ./build_image.sh -b from q1 directory.
2. Run ./build_image.sh -r from q1 directory.

# Question 2
This challenge was to create a file shredder using Go.

Automating testing is worth extra points.

## Use cases
Since OS's by default only unlink files when a user "removes" or "deletes" them, the raw data in the deleted file is still located on the disk, and can be retrieved by a skilled user. This program attempts to enable secure deletion of files by overwriting the program with random data generated and stored in RAM, by flushing the changes to disk, the data should be overwritten a total of three times, and then unlinked. This could be used as a more sensitive deletion tool for documents that contain personal or sensitive information.

The code could be incorporated into a larger program, or shell script that could be used to securely erase all files on a disk.

Due to the destructive nature of this program, it should be used with care and caution. The possibility exists for misuse, and as such I am not liable for any damages caused by developing with / running of the following code.

## Advantages
1. The program checks permissions, and if the appropriate permissions are not available, the program will exit gracefully.
2. If the program cannot find the file, it will exit gracefully.
3. The program ensures all bytes have been written in every iteraton before continuing, ensuring that the file is completely overwritten 3 times, or else an error will occur, and the program will exit gracefully.
4. The program will produce user-readable error messages for issues like seeking, writing, syncing, etc...
5. Evaluates symlinks, ensuring desired targeted data is deleted.

## Drawbacks
1. The program currently is really inefficient for large files, due to the way that the random array is created in one sweep

2. On newer systems, like SSD's, because the hardware has wear-levelling algorithms, which control which physical addresses bits are written to, we cannot be certain that all the data has been physically overwritten.

## Edge Cases
1. No input - Handled by returning 1 to caller.
2. Error evaluating symlink - Handled by returning 2 to caller
3. File existence - Handled by getting file info, if this fails, likely because the file doesn't exist. returns 2 to caller
4. Directories - Handled by not intercating with directories, aborting if directory is detected. Returns code 3 to caller.
5. Empty - Handled by checking if file is empty. Returns code 4 to caller.
6. File permissions - Handled by exiting if insufficient permissions are detected. Returns code 5 to caller.
7. File opening error - Handled by returning code 6 to caller.
8. Ambigious/Negative user input - Handled by returning code 7 to caller.
9. Random data generation failure - Handled by returning code 8 to caller.
10. Seeking errors - Handled by returning code 9 to caller.
11. Writing errors - Handled by returning code 10 to caller.
12. Syncing errors - Handled by returning code 11 to caller
13. Unlinking errors - Handled by returning code 12 to caller.
14. Full write isn't completed - Partially handled by re-attempting iteration indefinetly, added support for seeking to position of last write, and continuing write, once all bytes have been written, return.
15. File contents/length changing during execution - Not handled, would need to check at the start of every iteration, or poitentially before every write.
16. OS level abstraction - Not mitigated. Certain filesystem features, such as journaling and caching may prevent the proper overwriting of the program on the disk itself.

## Methodology
1. Import required libraries and check if args were recieved.
- fmt for string formatting.
- crypto/rand for secure, reliable random number generation.
- os for interacting with system in platform invariant way.
- If no arg was given return with error code 1.

2. Get the file path from args and attempt to retrieve the file info.
- In this step, need to evaluate symlink prior to stat.
- If this step fails, the program will exit with a return code of 2, likely because the directory cannot exist.

3. Check if the argument is a file or directory.
- If this step finds a directory, the program will exit with a return code of 3.

4. Ensure the file isn't empty.
- If it is, return error code 4.

4. Attempt to open the file with read/write permissions.
- Exits with return code of 5 if permissions are insufficient.
- Exits with a return code of 6 if there is another unspecified error.

5. Get the user to confirm that file deletion is okay.
- Y or y to confirm
- N or n to deny, will lead to the program returning with exit code 5.
- If any other pattern is detected, cannot be sure of user intention, abort immediately. Will also return exit code 7.

6. Loop 3 times, generate a random array of bytes the same length of the file, seek to the start of the file, and attempt to overwrite the file with the random array.
- If a random number generation error occurs program will return with exit code 8.
- If a seek error occurs, program will return with exit code 9.
- If an overwrite error occurs, program will return with exit code 10.
- If a sync error occurs , program will return with exit code 11.

7. Close the file, and attempt to unlink it.
- If an error occurs with unlinking, program will return with exit code 12.

8. If all above steps complete successfully, the program will return with error code 0.

## Bugs
The testing harness is currently not fully functional. Despite the fact that all tests pass, the harness is not able to verify that all the tests pass.

## Testing Automation (WIP)
Run 'make all', and interact with user prompts to complete the full testing suite. Next step is to interact with stdin via a script to feed data to program.

For now, the makefile can be used to perform some tests. The tests do the following:
1. make all - Runs the subsequent commands in proper order.
2. make - Runs the program on both a file with user read/write permissions (Should pass), and on a file with superuser read/write (Should fail).
3. make permissions - Re-runs the superuser permission file in an elevated go executable (Should work).
4. make invalid - Attempts to run the program on an invalid file (Should fail)
5. make changing - Undefined behavior. Not implemented yet.
6. make restore - Helper routine to restore sample files from backups with appropriate permissions
