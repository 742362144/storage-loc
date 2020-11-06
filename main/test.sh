../build.sh


./server -host 133.133.133.127 -port 50051 -deepth 0 &
spid=`ps -ef | grep "server" | grep -v grep | awk '{print $2}'`
echo $spid

# load
./benchmark -host 133.133.133.127 -port 50051 -model 2 -parallel 16 -deepth 1 -size 128 -num 10000 >> result

# run
for((i=1;i<=8;i++));
do
  ./benchmark -host 133.133.133.127 -port 50051 -model 1 -parallel 16 -deepth $i -size 128 -num 10000 >> result
#  echo $(expr $i \* 3 + 1);
done

kill $spid

for((i=1;i<=8;i++));
do
  ./server -host 133.133.133.127 -port 50051 -deepth $i &
  spid=`ps -ef | grep "server" | grep -v grep | awk '{print $2}'`
  echo $spid
  ./benchmark -host 133.133.133.127 -port 50051 -model 1 -parallel 16 -deepth 0 -size 128 -num 10000 >> result
  kill $spid
done
