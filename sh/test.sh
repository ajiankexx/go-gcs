#!/bin/bash
#tag: learn shell
firstname="nie"
secondname="tuan"

read x
if [ $x -le 5 ]; then
    echo "$x less than 5"
elif [ $x -ge 5 ]; then
    echo "$x greater than 5"
else
    echo "$x is 5"
fi

for i in {1..5}
do
    echo $i
done

function say_hello() {
    echo "hello $1"
}
say_hello "nietuan"
