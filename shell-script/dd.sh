i=1
while [ $i -lt 2 ]
do
	dd if=/dev/zero of=/root/mount2/file bs=1024 count=1000000000
	((i ++))
done
