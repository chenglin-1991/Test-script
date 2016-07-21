function shell_make()
{
	path=`pwd`
	cd $path
	make && make install
}

for i in 1 2 3
do
	{
		shell_make
    }&
done
