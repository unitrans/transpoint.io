# tpoint.io

# Load test
echo "GET http://46.101.248.60:8088/ping" | vegeta attack -duration=30s -rate=6000 -keepalive=true | tee results.bin | vegeta report -reporter=plot -output=index.html 

ab -n 10000 -c300 -k http://46.101.248.60:8088/ping

DO 512 Mb, ping - 20ms, 9988rps
