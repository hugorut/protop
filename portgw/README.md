# Port GW

The gateway for interactions with port services.

## How to run

The simplest way to run portgw on your machine is through docker.

First build the image:

```bash
$ docker build -t portgw . -f Dockerfile
```

Then run the image, binding your local port to the container:

```bash
$ docker run -p 8080:8080 portgw
```

You should see an output similar to below:

```bash
2021/01/23 18:36:59 server starting on port 8080
```

## Endpoints

### POST /ports/file/{provider}/upload

Provides a simple way to upload a file into the port system. The supported providers at the time of writing this readme:

```bash
> OS - machine file system
```

Future supported providers:

```bash
> S3
> Azure Blob
> GCP Cloud Storage
```

Currently the upload only supports JSON format files.

#### Request parameters

| name     	| type 	| description                                                                 	| examples                      	|
|----------	|------	|-----------------------------------------------------------------------------	|-------------------------------	|
| provider 	| path 	| the file processor provider you wish to use, currently only OS is supported 	| os                            	|
| location 	| body 	| the absolute location of the file you wish to persist to the system         	| /Path/to/your/file/ports.json 	|

#### Example request

```bash
$ curl -X "POST" -d '{"location": "/Path/to/your/file/ports.json"}' localhost:8080/ports/file/os/upload
```

If all goes well, tailing your docker logs should reveal an output similar to below:

```bash
INFO[0045] storing port Lusaka                          
INFO[0045] storing port Bulawayo                        
INFO[0045] storing port Harare                          
INFO[0045] storing port Mutare                          
INFO[0045] storing port Qui Nhon                        
INFO[0045] storing port Ho Chi Minh, VICT               
INFO[0045] storing port Vung Tau                        
INFO[0045] storing port Espiritu Santo  
```



 