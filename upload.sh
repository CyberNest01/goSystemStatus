CGO_ENABLED=0 go build
rsync -avz htopNovin root@185.53.143.180:/root/systemstatus/
rsync -avz htopNovin root@185.53.143.186:/root/systemstatus/
rsync -avz htopNovin root@185.53.141.120:/root/systemstatus/ -e 'ssh -p 2212'