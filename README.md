# TextFinder

在特定路径查找或替换指定字符串，可设置过滤的文件名以及目录层级

## Usage

### 指定路径

`./main -p /var/www -f find`

### 遍历所有目录

`./main -p /var/www -f find -l -1`

### 遍历一级目录

`./main -p /var/www -f find -l 1`

### 只查找特定文件名

`./main -p /var/www -f find -o .go`

### 某些文件名不查找

`./main -p /var/www -f find -e .md`