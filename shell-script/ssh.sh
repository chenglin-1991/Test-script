#!/bin/bash

ssh node22 "/root/sh/2-fdisk.sh"
if [ $? -eq 0 ]
then
	echo "====================ssh  node22 ok!==========================="
fi

ssh node23 "/root/sh/2-fdisk.sh"
if [ $? -eq 0 ]
then 
	echo "ssh node23 ok!"
	echo "====================ssh  node23 ok!==========================="
fi

#ssh node101 "/root/sh/2-fdisk.sh"
#if [ $? -eq 0 ]
#then 
	#echo "====================ssh  node101 ok!==========================="
#fi

#ssh node102 "/root/sh/2-fdisk.sh"
#if [ $? -eq 0 ]
#then 
	#echo "====================ssh  node102 ok!==========================="
#fi

exit 0
