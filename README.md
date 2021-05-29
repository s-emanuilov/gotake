# gotake
## Fast, reliable and easy file downloads


CLI for blazing fast downloads, available for Mac, Linux and Windows.

## Install
### macOS
```bash
sudo curl -L "https://github.com/simeonemanuilov/gotake/releases/download/0.5/gotake-darwin-x86_64" -o /usr/local/bin/gotake && sudo chmod +x /usr/local/bin/gotake
```
### Linux
```bash
sudo curl -L "https://github.com/simeonemanuilov/gotake/releases/download/0.5/gotake-linux-x86_64" -o /usr/local/bin/gotake && sudo chmod +x /usr/local/bin/gotake
```

## Examples

### Quick download (auto find optimal connection number)
```bash
gotake http://sample.li/face.png
```

### Download and print summary
```bash
gotake http://sample.li/boat.jpg -i
```

### Download in verbose mode
```bash
gotake http://sample.li/tesla.jpg -v
```

### Download with different type of connections
```bash
gotake http://sample.li/tesla.jpg -c=10 -a=False
```

## Benchmarks
![Quick Benchmark](/docs/images/quick-benchmark.png)

## Documentation
Check the documentation and available flags
```bash
gotake -h
```
![gotake Documentation](/docs/images/help.png)

## ⚒ Compile from source
The other way to install **gotake** is to clone its GitHub repository and build it from source. 
That is the common way if you want to make changes to the code base.

```bash
git clone https://github.com/simeonemanuilov/gotake
cd gotake
go build
```

## License
**gotake** is an Open-Source Project, and you can contribute to it in many ways. 

