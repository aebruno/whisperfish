#!/bin/bash

#==============================================================================
# Whisperfish build script
#==============================================================================
#
# This script performs the following functions:
#
# - Bootstrap github.com/therecipe/qt Go bindings
# - Create i386 go binary for use on mersdk
# - Runs qtmoc and qtminmal code generation
# - Compiles whisperfish
# - Build translation *.qm files using lrelease
# - Deploys whisperfish to Sailfish emulator
# - Creates RPM release for Jolla (armv7hl) device
#

APPNAME=harbour-whisperfish
VERSION=$(git describe --long --tags --dirty --always 2>/dev/null | cut -f2 -d'v')
GOQT_VERSION=d874b0a4b22e34a1cc253218e2f4ca09c9fe686d

case "$1" in
        bootstrap-qt)
            pushd .
            mkdir -p $GOPATH/src/github.com/therecipe
            cd $GOPATH/src/github.com/therecipe
            git clone https://github.com/therecipe/qt
            cd qt
            git checkout $GOQT_VERSION
            cd cmd/qtsetup
            go install .
            cd ../qtmoc
            go install .
            cd ../qtminimal
            go install .
            popd
            ;;
        setup-sdk)
            GOARCH=386 $GOROOT/src/run.bash
            ;;
        prep)
            qtmoc $PWD/settings
            qtmoc $PWD/model
            qtmoc $PWD/worker
            $GOPATH/bin/qtminimal sailfish-emulator $PWD
            ;;
        prep-arm)
            qtmoc $PWD/settings
            qtmoc $PWD/model
            qtmoc $PWD/worker
            $GOPATH/bin/qtminimal sailfish $PWD
            ;;
        compile)
            rm -f $APPNAME
            GOOS=linux GOARCH=386 CGO_ENABLED=1 CC=/opt/cross/bin/i486-meego-linux-gnu-gcc \
            CXX=/opt/cross/bin/i486-meego-linux-gnu-g++ CPATH=/srv/mer/targets/SailfishOS-i486/usr/include \
            LIBRARY_PATH=/srv/mer/targets/SailfishOS-i486/usr/lib:/srv/mer/targets/SailfishOS-i486/lib \
            CGO_LDFLAGS=--sysroot=/srv/mer/targets/SailfishOS-i486/ \
            go build -ldflags="-s -w -X main.Version=$VERSION" -tags="minimal sailfish_emulator" \
            -installsuffix=sailfish_emulator -o $APPNAME 2>&1 | egrep '\.go'
            #-installsuffix=sailfish_emulator -o $APPNAME
            if [ ! -f $APPNAME ]; then echo "Failed to compile."; fi
            ;;
        rpm)
            rm -f $APPNAME
            GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 CC=/opt/cross/bin/armv7hl-meego-linux-gnueabi-gcc \
            CXX=/opt/cross/bin/armv7hl-meego-linux-gnueabi-g++ CPATH=/srv/mer/targets/SailfishOS-armv7hl/usr/include \
            LIBRARY_PATH=/srv/mer/targets/SailfishOS-armv7hl/usr/lib:/srv/mer/targets/SailfishOS-armv7hl/lib \
            CGO_LDFLAGS=--sysroot=/srv/mer/targets/SailfishOS-armv7hl/ \
            go build -ldflags="-s -w -X main.Version=$VERSION" -tags="minimal sailfish" \
            -installsuffix=sailfish -o $APPNAME 2>&1 | egrep '\.go'
            if [ ! -f $APPNAME ]; then
                echo "Failed to compile."
            else
                mb2 -x -t SailfishOS-armv7hl build
            fi
            ;;
        i18n-up)
            # Update translations
            for filename in qml/i18n/whisperfish_*.ts; do
                sb2 lupdate qml/ -ts ${filename}
            done
            ;;
        i18n)
            # Compile translations
            for filename in qml/i18n/whisperfish_*.ts; do
                name="${filename%.*}"
                sb2 lrelease -idbased  ${name}.ts -qm ${name}.qm
            done
            ;;
        deploy)
            mb2 -x -s rpm/$APPNAME.spec -d "SailfishOS Emulator" deploy  --sdk
            ;;
        clean)
            rm -f settings/moc*
            rm -f model/moc*
            rm -f worker/moc*
            ;;
        clean-qt)
            rm -Rf  $GOPATH/src/github.com/therecipe/qt
            rm -Rf  $GOPATH/pkg/linux_amd64/github.com/therecipe/qt
            rm -f $GOPATH/bin/{qtmoc,qtsetup,qtminimal}
            mkdir -p $GOPATH/src/github.com/therecipe
            ;;
        rebuild-qt)
            pushd .
            cd $GOPATH/src/github.com/therecipe/qt
            cd cmd/qtsetup
            go install .
            cd ../qtmoc
            go install .
            cd ../qtminimal
            go install .
            popd
            ;;
        *)
            echo $"Usage: $0 {prep|prep-arm|compile|rpm|deploy}"
            exit 1
 
esac
