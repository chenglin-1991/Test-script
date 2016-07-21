line1=`ps -ef | grep "/usr/local/sbin/glusterfsd"`
## | cut -d " " -f 7`
echo $line1
#line2=`echo $line1 | cut -d " " -f 1`

#echo "$line2"
#if [ $line2 -ne 0 ]
#then
        #kill -9 $line2
#fi
