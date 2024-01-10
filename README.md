# CBOR tool to encode, decode and tag strings (from input or strings in files like logs)
#Build the tool into a docker and run from there
```
$ docker build -t cbor-tool:0.1 .
```
## then
```
$ docker run cbor-tool:0.1 /tool
ERR: Must specify action(-a)
Actions: Encode Decode Tag
````
