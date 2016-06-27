rm -rf ./generated/*
go-bindata -pkg generated -o ./generated/bindata.go files/
