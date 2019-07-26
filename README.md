# GetIt
![loc](https://tokei.rs/b1/github/nektro/getit)
[![license](https://img.shields.io/github/license/nektro/getit.svg)](https://github.com/nektro/getit/blob/master/LICENSE)
[![discord](https://img.shields.io/discord/551971034593755159.svg)](https://discord.gg/P6Y4zQC)

Download files from an open directory into your local filesystem while preserving the folder structure

[![buymeacoffee](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/nektro)

## Getting Started

### Installing
```
$ go get github.com/nektro/getit
```

### Usage
```
$ ~/go/bin/getit -url {URL}
```

> Note: GetIt will always save files relative to the current working directory. So running `getit -url https://example.com/path/to/dir/`, if it contains a link to `file.txt`, it will be saved to `./example.com/path/to/dir/file.txt`.

## Built With
- https://github.com/PuerkitoBio/goquery
- https://github.com/deckarep/golang-set
- https://github.com/nektro/go-util
- https://github.com/schollz/progressbar

## Contributing
Issues and Pull Requests welcome!

## License
MIT. See [LICENSE](LICENSE) for more info.

## Acknowledgments
- The team behind `wget -m` and all of [GNU Wget](https://www.gnu.org/software/wget/)
- [The Eye](https://the-eye.eu/) community
