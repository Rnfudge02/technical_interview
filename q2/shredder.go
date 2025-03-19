package main
//Program developed by Robert Fudge 2025
//Licensed under the Apache 2.0 license

//Import dependecies
import (
    "fmt"
	"crypto/rand"
    "os"
	"path/filepath"
)

//Main function, runs shred
func main() {
	//Ensure filename is passed
	if len(os.Args) < 2 {
		fmt.Println("Usage: program <file_path>")
		os.Exit(1)
	}

	//Get command line arguments
	file_path := os.Args[1]

	shred(file_path)

	os.Exit(0)
}

//Main function
func shred(path string) {
	//Resolve any symlinks to the actual file
    resolvedPath, err0 := filepath.EvalSymlinks(path)

    if err0 != nil {
        fmt.Printf("Error resolving path: %v\n", err0)
        os.Exit(2)
    }
    path = resolvedPath

	//Attempt to get file info
	file_info, err1 :=os.Stat(path)

	//If the error is not a pointer to 0, then the stat failed
    if err1 != nil {
        fmt.Println("Error reading file", err1)
        os.Exit(2)
    }

	//Chose to not shred directories
	if file_info.IsDir() == true {
		fmt.Printf("%s is a directory\n", file_info.Name())
        os.Exit(3)
	}

	//If file is empty, there is nothing to do, exit with return code 0
	if file_info.Size() == 0 {
		fmt.Printf("%s is empty\n", file_info.Name())
        os.Exit(4)
	}

	//Attempt to open first to verify permissions (read and write)
	//Use OpenFile to check permissions when opening
    file, err_open := os.OpenFile(path, os.O_RDWR, 0666)

    if err_open != nil {
        if os.IsPermission(err_open) {
            fmt.Printf("Insufficient permissions to read/write %s\n", path)
            os.Exit(5)
        } else {
            fmt.Printf("Error opening file: %v\n", err_open)
            os.Exit(6)
        }
    }

	fmt.Printf("Shredding file %s.\n", path)
	fmt.Printf("Name: %s\nSize %dB.\n", file_info.Name(), file_info.Size())

	//Confirm user wants to delete the file
	var conf string
	print("Are you sure you want to do this? (Y/N) ")
	fmt.Scan(&conf)

	//If we get either of these characters (technically strings)
	//We can continue to shred the file
	if conf == "Y" || conf == "y" {
		println("Approval confirmed, shredding...")

		size := file_info.Size()

		//"Shred" the file X amount of times
		for i := 0; i < 3; i++ {
			//Was going to redirect from /dev/random but that is not cross-platform
			rand_arr := make([]byte, size)
			_, err_rand := rand.Read(rand_arr)
			if err_rand != nil {
				fmt.Println("Error generating random data:", err_rand)
				os.Exit(8)
			}

			var n_w int64 = 0

			//
			for n_w < size {
				rem := size - n_w
				slice := rand_arr[n_w : n_w+rem]
				_, err := file.Seek(n_w, 0)
				if err != nil {
    				println("Error, failed to seek to next write location.")
					os.Exit(9)
				}

				n_wr_rec , err_wr := file.Write(slice)

				//Ensure only full rewrites are counted
				if err_wr != nil  {
					println("Overwrite error", err_wr)
					os.Exit(10)
				
				}

				err_sync := file.Sync()
				if err_sync != nil {
    				println("Error: could not sync file.")
					os.Exit(11)
				}

				n_w += int64(n_wr_rec)
				fmt.Printf("Bytes written till now %d.\n", n_w)
			}
		}

		//Close file and delete
		file.Close()
		err_rem := os.Remove(path)

		//If the file couldn't be unlinked notify the user and exit
		if err_rem != nil { 
			fmt.Printf("Error: Couldn't remove file %s.\n", path)
			os.Exit(12)
		} else {
			fmt.Printf("File %s shredded and destroyed.\n", path)
			return
		}

	//If the user does not want the files destroyed, or input can't be understood, abort immediately
	} else {
		println("Aborting.")
		os.Exit(7)
	}
}