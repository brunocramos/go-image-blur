# Go Image Blur Filter

Simple image blur filter in golang for studying purposes. It uses a 3x3 mask on each pixel to get the average RGB value.

Usage: 

```
go build image-blur
./image-blur [input-file] [output-file]
```

Ex.:
`
./image-blur lenna.png lenna-blurred.png
`