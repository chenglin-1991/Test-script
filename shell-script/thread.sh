#!/bin/bash 
tmpfile=$$.fifo
mkfifo $tmpfile
exec 4<>$tmpfile
rm $tmpfile
thred=4

{
for ((i = 1;i <= $thred; i++))
do
	echo;
done
}>&4

for ((i = 0;i < 100; i++))
do
	read
	(./touch_file.sh;echo >&4)&
done <&4
wait
exec 4>&-
