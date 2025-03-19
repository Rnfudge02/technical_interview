#!/bin/bash
#Robert Fudge
#Technical Interview Question 1 Script

#Creates and configures the linux kernel
make_kernel() {
    #Install needed dependencies
    sudo apt install wget build-essential gcc g++ libssl-dev bc qemu-system-x86 bison flex

    mkdir -p build
    cd build

    #Retrieve the linux kernel using wget, decompress, dont need to redownload
    if [ ! -f linux-6.1.38.tar.xz ]; then
        wget https://cdn.kernel.org/pub/linux/kernel/v6.x/linux-6.1.38.tar.xz
    fi
    tar -xvf linux-6.1.38.tar.xz
    cd linux-6.1.38

    #Kernel configuration using make
    make defconfig  #Start with default config
    
    #Enable ext4 support and serial console
    ./scripts/config --set-val CONFIG_EXT4_FS y
    ./scripts/config --set-val CONFIG_SERIAL_8250 y
    ./scripts/config --set-val CONFIG_SERIAL_8250_CONSOLE y

    make -j$(nproc)  #Compile with all CPU cores

    cd ../..
}

#Compile the minimal program
make_hello() {
    mkdir -p build/hello_world
    cd build/hello_world
    gcc -static -o hello_world ../../hello_world.c
    cd ../..
}

#Creates the ext4 filesystem, installing programs would go here
create_ext4_fs() {
    mkdir -p build
    cd build

    dd if=/dev/zero of=rootfs.ext4 bs=1M count=2048  #2GB image
    mkfs.ext4 -F rootfs.ext4  #Make ext4 fs

    #Create directory and mount it, this is how files will be injected
    mkdir ./rootfs
    sudo mount -o loop rootfs.ext4 ./rootfs

    #Use {} to avoid writing ./rootfs/... as many times
    sudo mkdir -p ./rootfs/{bin,dev,proc,sys,etc,home,mnt,root,tmp,var}

    sudo cp hello_world/hello_world ./rootfs/bin/

    #Create inodes (for console I/O)
    sudo mknod ./rootfs/dev/console c 5 1  #Console device
    sudo mknod -m 666 rootfs/dev/null c 1 3 #Null

    #Unmount the disk image
    sudo umount ./rootfs

    cd ..
}

#Runs the created image
run_image() {
   echo -e "[Image Creator] Starting QEMU..."
    cd build
    qemu-system-x86_64 \
        -kernel linux-6.1.38/arch/x86/boot/bzImage \
        -drive file=rootfs.ext4,format=raw \
        -append "root=/dev/sda rw init=/bin/hello_world console=ttyS0" \
        -nographic \
        -serial mon:stdio
    cd ..
}

#Parse command line arguments
while getopts "brh" options; do
    case ${options} in
        #Build - controls building of image
        b)
            echo "[Image Creator] Building minimal linux kernel with ext4 filesystem"

            make_kernel
            make_hello
            create_ext4_fs

            echo "[Image Creator] Kernel and FS build script done."
        ;;

        #Cross-compile - build for opposite architecture as target - NOT WORKING, issue with transferring build stages
        r)
            run_image
        ;;

        #Help command
        h)
            echo "Image Builder v0.1.0 - Developed by Robert Fudge"
            echo "Valid commands are listed below:"

            echo "ARGUMENT       NAME            INFO"
            echo "-b             Build           Builds the linux kernel with ext4 support and rootfs."
            echo "-h             Help            Displays the help menu."
            echo "-r             Run             Run the built image"
        ;;
    esac
done