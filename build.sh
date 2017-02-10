#!/bin/bash

APPNAME=harbour-whisperfish
QT_VERSION=5.2.0
QT_DIR=$HOME/Qt5.2.0
QT_DOC_DIR=$QT_DIR/5.2.0/gcc_64/doc/
VERSION=$(git describe --long --tags --dirty --always 2>/dev/null | cut -f2 -d'v')

case "$1" in
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
        setup-qt)
            QT_VERSION=$QT_VERSION QT_DIR=$QT_DIR QT_DOC_DIR=$QT_DOC_DIR QT_PKG_CONFIG=true $GOPATH/bin/qtsetup generate sailfish-emulator
            ;;
        setup-mer)
            GOARCH=386 $GOROOT/src/run.bash
            ;;
        prep)
            qtmoc $PWD/settings
            qtmoc $PWD/model
            qtmoc $PWD/worker
            #QT_VERSION=$QT_VERSION QT_DIR=$QT_DIR QT_DOC_DIR=$QT_DOC_DIR QT_PKG_CONFIG=true $GOPATH/bin/qtminimal sailfish-emulator $PWD
            $GOPATH/bin/qtminimal sailfish-emulator $PWD
            ;;
        prep-arm)
            qtmoc $PWD/client
            qtmoc $PWD/model
            QT_VERSION=$QT_VERSION QT_DIR=$QT_DIR QT_DOC_DIR=$QT_DOC_DIR QT_PKG_CONFIG=true $GOPATH/bin/qtminimal sailfish $PWD
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
        *)
            echo $"Usage: $0 {prep|prep-arm|compile|rpm|deploy}"
            exit 1
 
esac
