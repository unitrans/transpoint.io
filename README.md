# tpoint.io

# Load test
echo "GET http://46.101.248.60:8088/ping" | vegeta attack -duration=30s -rate=6000 -keepalive=true | tee results.bin | vegeta report -reporter=plot -output=index.html 

ab -n 10000 -c200 -k http://46.101.248.60:8088/ping

 - DO 512 Mb [http]: latency 20ms, 9988rps
 - Heroku [https]: latency 40ms, 2100rps

 
ab -n10000 -c200 -k https://transpoint.cleverapps.io/v1/translations/123

 - CleverCloud [https]: latency 1000ms -> 30ms, 1073rps
 - Heroku [https]: latency 300ms -> 50-70ms, 1172rps
 - Google CE [http]: latency 20-50ms, 1700rps (internal worktime 0.08-0.16ms) (n1 & f1 are same ... f1 seems to be faster)
 - DO [http]: latency 20-50ms (spikes 600ms), 1600rps
 
### From scratch

```
gcloud compute --project sodium-platform-99915 instances create "packer" --zone "europe-west1-b" --machine-type "f1-micro" \
--scopes "https://www.googleapis.com/auth/devstorage.full_control,https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/cloud-platform" \
--tags "http-server" \
--image "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/backports-debian-7-wheezy-v20150603" 
```

### build packer image

```
packer build packer.json 
```

### from Image

```
gcloud compute --project sodium-platform-99915 instances create "packer-auto" --zone "europe-west1-b" --machine-type "f1-micro" \
--scopes "https://www.googleapis.com/auth/devstorage.full_control,https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/cloud-platform" \
--tags "http-server" \
--image "packer-trio" 
```