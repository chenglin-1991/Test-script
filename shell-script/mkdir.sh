#!/bin/bash
no=0
i=0
while ((i < 100))
do
    for((i=1;i<=90000;i++))
    do
        for((j=1;j<=9000;j++))
        do
            for((h=1;h<=90;h++))
            do
                mkdir -p /cluster2/test/$1/$i/$j/$h
                touch /cluster2/test/$1/$i/$j/$h/file
                #echo "/cluster2/test/$i/$j/$h created"
                let no=$no+1
                echo $no
            done
        done
    done
    ((i++))
    rm -rf /cluster2/test/*
done

rm -rf  /cluster2/test/*
