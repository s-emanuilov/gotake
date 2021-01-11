## GoTake

The idea: make async download to a file, using range request.
At the end - we combine the file chunks.

## Links

- Make the CLI tooling, without the external library: 
  https://levelup.gitconnected.com/tutorial-how-to-create-a-cli-tool-in-golang-a0fd980264f
  
- Parallel download realization in Go:
    https://coderwall.com/p/uz2noa/fast-parallel-downloads-in-golang-with-accept-ranges-and-goroutines
  
- Alternate download if the other is not working fine:
    https://stackoverflow.com/questions/11692860/how-can-i-efficiently-download-a-large-file-using-go/33853856
  
## Benchmarks
TODO