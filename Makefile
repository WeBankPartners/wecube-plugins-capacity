current_dir=$(shell pwd)
version=$(PLUGIN_VERSION)
project_name=$(shell basename "${current_dir}")

clean:
	rm -rf server/server
	rm -rf server/public/*
	rm -rf ui/dist
	rm -rf ui/plugin

build: clean
	chmod +x ./build/*.sh
	docker run --rm -v $(current_dir):/go/src/github.com/WeBankPartners/$(project_name) --name build_$(project_name) golang:1.12.5 /bin/bash /go/src/github.com/WeBankPartners/$(project_name)/build/build-server.sh
	docker run --rm -v $(current_dir):/app/wecube-plugins-capacity --name node-build node /bin/bash /app/wecube-plugins-capacity/build/build-ui.sh
	mv ui/dist/index.html server/public/
	mv ui/dist/favicon.ico server/public/
	mv ui/dist/capacity/* server/public/

image: build
	docker build -t $(project_name):$(version) .

package: image
	mkdir -p plugin
	cp -r ui/plugin/* plugin/
	zip -r ui.zip plugin
	rm -rf plugin
	cp build/register.xml ./
	cp doc/init.sql ./
	sed -i "s~{{PLUGIN_VERSION}}~$(version)~g" ./register.xml
	docker save -o image.tar $(project_name):$(version)
	zip  $(project_name)-$(version).zip image.tar init.sql register.xml ui.zip
	rm -f register.xml
	rm -f init.sql
	rm -f ui.zip
	rm -rf ./*.tar
	docker rmi $(project_name):$(version)

upload: package
	$(eval container_id:=$(shell docker run -v $(current_dir):/package -itd --entrypoint=/bin/sh minio/mc))
	docker exec $(container_id) mc config host add wecubeS3 $(s3_server_url) $(s3_access_key) $(s3_secret_key) wecubeS3
	docker exec $(container_id) mc cp /package/$(project_name)-$(version).zip wecubeS3/wecube-plugin-package-bucket
	docker rm -f $(container_id)
	rm -rf $(project_name)-$(version).zip