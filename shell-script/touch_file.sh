echo -n "Input the number of file:"

read num

i=0

while [ $i -lt $num ]
do
    touch /cluster2/test/file$i
	i=`expr $i + 1`
    echo $i
    sleep 1
done

echo "files number is $num"

