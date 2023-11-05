# gotake: High-Speed File Downloader

`gotake` is a CLI tool designed to expedite file downloading processes. Engineered for reliability and ease of use, it's compatible across Mac, Linux, and Windows platforms.  
Under the hood, `gotake` employs either standard HTTP requests or range requests to fetch files, the latter enabling concurrent downloads of file chunks which are subsequently assembled to form the complete file.

## Install
Install `gotake` with a simple one-liner. Below are the instructions for macOS and Linux:

### macOS

```bash
sudo curl -L "https://github.com/simeonemanuilov/gotake/releases/download/0.5/gotake-darwin-x86_64" -o /usr/local/bin/gotake && sudo chmod +x /usr/local/bin/gotake
```

### Linux

```bash
sudo curl -L "https://github.com/simeonemanuilov/gotake/releases/download/0.5/gotake-linux-x86_64" -o /usr/local/bin/gotake && sudo chmod +x /usr/local/bin/gotake
```

## Examples
Launching a download is straightforward with gotake. Here are some examples to get you started:

### Quick download (auto find optimal connections amount)

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

## How it works

<p align="center">
    <img src="/docs/images/schema.png" alt="Schema for file downloads with gotake">
</p>

## Benchmarks

<p align="center">
    <img src="/docs/images/quick-benchmark.png" alt="Quick Benchmark">
</p>

## Documentation

Check the documentation and available flags

```bash
gotake -h
```

<p align="center">
    <img src="/docs/images/docs.png" alt="gotake Documentation" style="width: 60%">
</p>

## Flags
### -a, --auto (bool) 
default: **True**  
Choose the number of connections automatically, depending on the file size.

### -c, --connections (int)
default: **50**  
Number of connections. You need to disable **auto** mode for using a specific number.

### -f, --filename (string)
default: **extracted from URL**  
Define different name to the target file if you prefer.

### -i, --info (bool)
default: **False**  
Print a summary for the download after finish.

### -s, --standard (bool)
default: **False**  
Force the standard method of download (not using Range-Request).

### -v, --verbose (bool)
default: **False**  
Show information about the current state of download process. Useful for debug.

### -h, --help
Print a documentation with all available flags.

## ⚒ Compile from source

The other way to install **gotake** is to clone its GitHub repository and build it from source. That is the common way
if you want to make changes to the code base. The tool is made with [Golang](http://golang.bg/).


```bash
git clone https://github.com/simeonemanuilov/gotake
cd gotake
go build
```

## License

**gotake** is an Open-Source Project, and you can contribute to it in many ways. 

