#!/bin/bash

Xstart=400      
Ystart=100     
Xend=650       
Yend=650       
step=20

curl -X POST http://localhost:17000 -d "reset"
curl -X POST http://localhost:17000 -d "white"
curl -X POST http://localhost:17000 -d "update"
curl -X POST http://localhost:17000 -d "green"
curl -X POST http://localhost:17000 -d "figure $Xstart $Ystart"
curl -X POST http://localhost:17000 -d "update"

while true; do
    curl -X POST http://localhost:17000 -d "update"


    for ((Y=$Ystart; Y <= $Yend; Y+=step)); do
        curl -X POST http://localhost:17000 -d "move 0 $step"
        curl -X POST http://localhost:17000 -d "update"
    done


    for ((X=$Xstart; X <= $Xend; X+=step)); do
        curl -X POST http://localhost:17000 -d "move $step 0"
        curl -X POST http://localhost:17000 -d "update"
    done

    
    while [ $X -gt $Xstart ]; do
        curl -X POST http://localhost:17000 -d "move $((-step)) 0"
        curl -X POST http://localhost:17000 -d "update"
        X=$((X - step))
    done
    while [ $Y -gt $Ystart ]; do
        curl -X POST http://localhost:17000 -d "move 0 $((-step))"
        curl -X POST http://localhost:17000 -d "update"
        Y=$((Y - step))
    done

   
    for ((X=$Xstart; X >= $Xstart - 200; X-=step)); do
        curl -X POST http://localhost:17000 -d "move $((-step)) 0"
        curl -X POST http://localhost:17000 -d "update"
    done

    
    while [ $X -lt $Xstart ]; do
        curl -X POST http://localhost:17000 -d "move $step 0"
        curl -X POST http://localhost:17000 -d "update"
        X=$((X + step))
    done
done
