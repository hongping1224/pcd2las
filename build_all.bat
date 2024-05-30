#windows
go build

# linux
set GOOS=linux
set GOARCH=amd64
go build -o pcd2las_linux_amd64
