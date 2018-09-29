#bin/bash

echo start

echo clean last_version_files
rm -f *.pb.go
echo clean success
echo start compailer
protoc --go_out=. *.proto
echo compailer success

echo done
