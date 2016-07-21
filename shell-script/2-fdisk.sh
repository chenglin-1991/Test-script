#!/bin/bash 
#j=1
#while [ $j -lt 3 ]
#do
	#for i in `attr -lq ~/node22-0$j/brick22-0$j/`
	#do
		#setfattr -x trusted.$i ~/node22-0$j/brick22-0$j
	#done
		#((j ++))
#done
umount ~/node21-01/
umount ~/node21-02/
umount ~/node21-03/

mkfs.xfs -f /dev/sdb1
mkfs.xfs -f /dev/sdc1
mkfs.xfs -f /dev/sdd1

mount -a

#rm -rf ~/node21-*

mkdir -p ~/node21-01/brick21-01
mkdir -p ~/node21-02/brick21-02
mkdir -p ~/node21-03/brick21-03
