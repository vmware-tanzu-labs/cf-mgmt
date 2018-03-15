covfile=$(mktemp)
echo "mode: set" > $covfile
for pkg in `glide nv`; do
  covout=$(mktemp)
  go test -coverprofile=$covout ${pkg/.../.}
  tail -n +2 ${covout} >> ${covfile}
done
go tool cover -html=${covfile}
