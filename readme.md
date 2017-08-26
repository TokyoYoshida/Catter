# Catter

Catter is a container based virtualization software.


- Install

```
Install the go langage.

sudo apt-get install libcgroup-dev

apt-get install cgroup-bin

sudo apt-get install debootstrap
sudo debootstrap --arch i386 xenial ./rootfs http://archive.ubuntu.com/ubuntu

mkdir upper workdir containerRoot

sudo mount -t overlayfs -o lowerdir=rootfs,upperdir=upper,workdir=workdir overlayfs containerRoot

mkdir homeShare
cd homeShare;touch a_shared_file
```

- Build

```
go build catter.go
```

- Run

```
sudo ./catter
```

- License

MIT
