#! /bin/bash
p=`pidof chi-auto-ocr`
if [[ $p -gt 0 ]]; then
    echo "killing $p"
    kill -9 $p
fi
go build -o chi-auto-ocr .
mv chi-auto-ocr /usr/local/bin/
nohup chi-auto-ocr > chi-auto-ocr.log 2>&1 &
p2=`pidof chi-auto-ocr`
echo "$p2 is running"
