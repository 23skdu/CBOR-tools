# perl-tool-template
#Build the tool into a docker and run from there
```
$ docker build -t tool:0.1 .
```
## then
```
perl-tool-template $ docker run tool:0.1 tool
ERR: Must specify action(-a)
Actions: SomeFunction 
````
