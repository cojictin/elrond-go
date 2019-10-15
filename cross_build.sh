#!/usr/bin/env bash

package=node
package_split=(${package//\// })
package_name=${package_split[-1]}

platforms=("windows/amd64" "windows/386" "darwin/amd64" "darwin/386" "linux/amd64" "linux/386" "linux/arm" "linux/arm64")
BASEDIR=$(pwd)/cross_build


for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$BASEDIR/$GOOS/$GOARCH/$package_name
    echo "Building $output_name"
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    pushd cmd/node
    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name
    popd
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! '
    fi
done
