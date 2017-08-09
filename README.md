# httpfs

Go lang http.FileSystem packages

## http.FileSystem packages

|          package         |        description        |        data source          |
|--------------------------|---------------------------|-----------------------------|
| [zipfs](zipfs)           | zip filesystem            | zip file                    |
| [gitfs](gitfs)           | git filesystem            | git bare repository         |
| [unionfs](unionfs)       | as union filesystem       | multiple http.FileSystem(s) |
| [indexfs](indexfs)       | custom directory index    | http.FileSystem             |
| [templatefs](templatefs) | template engine execute   | http.FileSystem             |

## Examples

|        example code         |                   description                        | 
|-----------------------------|------------------------------------------------------|
| [tinyweb](examples/tinyweb) | tiny http static server (dir, zip, git)              | 

## Author

[gnue](https://github.com/gnue)

## License

[MIT](LICENSE.txt)

