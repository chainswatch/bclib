#!/bin/sh

COVERAGE_DIR="${COVERAGE_DIR:-coverage}"
PKG_LIST=$(go list ./... | grep -v /vendor/ | tr -d '\r')

# Create the coverage files directory
mkdir -p "$COVERAGE_DIR"

echo 'mode: count' > "${COVERAGE_DIR}"/coverage.cov
# Create a coverage file for each package
for package in ${PKG_LIST}; do
    go test -covermode=count -coverprofile "${COVERAGE_DIR}/${package##*/}.cov" "$package"
    tail -q -n +2 "${COVERAGE_DIR}/${package##*/}.cov" >> "${COVERAGE_DIR}"/coverage.cov
done

# Display the global code coverage
go tool cover -func="${COVERAGE_DIR}"/coverage.cov

# If needed, generate HTML report
if [ "$1" = "html" ]; then
    go tool cover -html="${COVERAGE_DIR}"/coverage.cov -o coverage.html
fi

# Remove the coverage files directory
rm -rf "$COVERAGE_DIR"
