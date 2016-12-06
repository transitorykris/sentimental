all:
	docker build -t timberslide/sentimental .

push:
	docker push timberslide/sentimental

clean:
	docker rmi -f timberslide/sentimental
