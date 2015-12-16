# unitrans.me

http://unitrans.me

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


## Docker

```
docker rmi $(docker images -qf "dangling=true")

docker ps -a

docker build -t unitrans_image .

docker run --publish 6060:8088 --name unitrans_container --rm unitrans_image


# GCloud
docker tag unitrans_image eu.gcr.io/unitrans-1107/unitrans_image

gcloud config set compute/zone europe-west1-c
gcloud docker push eu.gcr.io/unitrans-1107/unitrans_image

gcloud container clusters create guestbook --num-nodes 1 --machine-type g1-small 
gcloud container clusters list
gcloud container clusters describe guestbook

#kubernetes
kubectl create -f redis-service.yaml
kubectl create -f redis-controller.yaml
kubectl create -f service-controller.yaml
kubectl create -f service-service.yaml


kubectl get services
kubectl describe services frontend

#docker
kubectl get pods -o wide
gcloud compute ssh gke-guestbook-a2f9068e-node-age1

sudo docker ps
sudo docker logs -f 54e9df6efa7d
```
