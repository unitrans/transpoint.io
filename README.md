# tpoint.io

# Load test
echo "GET http://46.101.248.60:8088/ping" | vegeta attack -duration=30s -rate=6000 -keepalive=true | tee results.bin | vegeta report -reporter=plot -output=index.html 

ab -n 10000 -c200 -k http://46.101.248.60:8088/ping

 - DO 512 Mb [http]: latency 20ms, 9988rps
 - Heroku [https]: latency 40ms, 2100rps

 
ab -n10000 -c200 -k https://transpoint.cleverapps.io/v1/translations/123

 - CleverCloud [https]: latency 1000ms -> 30ms, 1073rps
 - Heroku [https]: latency 300ms -> 50-70ms, 1172rps
 - Google CE [http]: latency 20-50ms, 1700rps (internal worktime 0.08-0.16ms)

