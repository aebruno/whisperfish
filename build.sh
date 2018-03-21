#!/bin/bash

#==============================================================================
# Whisperfish build script
#==============================================================================
#
# This script performs the following functions:
#
# - Bootstrap github.com/therecipe/qt Go bindings
# - Create i386/arm go binary for use on mersdk with sb2
# - Runs qtmoc and qtminmal code generation
# - Compiles whisperfish
# - Build translation *.qm files using lrelease
# - Deploys whisperfish to Sailfish emulator
# - Creates RPM release for Jolla (armv7hl) device
#

APPNAME=harbour-whisperfish
VERSION=$(git describe --long --tags --always 2>/dev/null | cut -f2 -d'v')
GOQT_VERSION=c6ada02b904734c7f78a8032acd2e6fee3e58dba
QT_VERSION=5.7.0

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
            GOARCH=arm GOARM=7 $GOROOT/src/run.bash
            ;;
        prep)
            # XXX hack to prevent qtmoc from searching vendor directory
            mv vendor ..
            QT_VERSION=$QT_VERSION qtmoc $PWD/settings
            QT_VERSION=$QT_VERSION qtmoc $PWD/model
            QT_VERSION=$QT_VERSION qtmoc $PWD/worker
            QT_VERSION=$QT_VERSION $GOPATH/bin/qtminimal sailfish-emulator $PWD
            mv ../vendor .
            ;;
        prep-arm)
            # XXX hack to prevent qtmoc from searching vendor directory
            mv vendor ..
            QT_VERSION=$QT_VERSION qtmoc $PWD/settings
            QT_VERSION=$QT_VERSION qtmoc $PWD/model
            QT_VERSION=$QT_VERSION qtmoc $PWD/worker
            QT_VERSION=$QT_VERSION $GOPATH/bin/qtminimal sailfish $PWD
            mv ../vendor .
            ;;
        compile)
            rm -f $APPNAME
            sb2 -O use-global-tmp -t SailfishOS-2.1.4.13-i486 -m sdk-build env \
            CGO_ENABLED=1 GOOS=linux GOARCH=386 CC=/srv/mer/toolings/SailfishOS-2.1.4.13/opt/cross/bin/i486-meego-linux-gnu-gcc \
            CXX=/srv/mer/toolings/SailfishOS-2.1.4.13/opt/cross/bin/i486-meego-linux-gnu-g++ \
            CGO_CFLAGS_ALLOW='.*' \
            CGO_CXXFLAGS_ALLOW='.*' \
            CGO_LDFLAGS_ALLOW='.*' \
            CPATH=/srv/mer/targets/SailfishOS-2.1.4.13-i486/usr/include:/srv/mer/targets/SailfishOS-2.1.4.13-i486/usr/include/qt5/QtCore:/srv/mer/targets/SailfishOS-2.1.4.13-i486/usr/include/qt5:/srv/mer/targets/SailfishOS-2.1.4.13-i486/usr/include/qt5/QtGui:/srv/mer/targets/SailfishOS-2.1.4.13-i486/usr/include/qt5/QtNetwork:/srv/mer/targets/SailfishOS-2.1.4.13-i486/usr/include/qt5/QtQuick:/srv/mer/targets/SailfishOS-2.1.4.13-i486/usr/include/qt5/QtQml:/srv/mer/targets/SailfishOS-2.1.4.13-i486/usr/include/sailfishapp \
            LIBRARY_PATH=/srv/mer/targets/SailfishOS-2.1.4.13-i486/usr/lib:/srv/mer/targets/SailfishOS-2.1.4.13-i486/lib \
            CGO_LDFLAGS="--sysroot=/srv/mer/targets/SailfishOS-2.1.4.13-i486/" \
            CGO_CFLAGS="--sysroot=/srv/mer/targets/SailfishOS-2.1.4.13-i486/" \
            /home/src1/projects/goroot/go/bin/linux_386/go build \
            -ldflags="-s -w -X main.Version=$VERSION" -tags="minimal sailfish_emulator" \
            -installsuffix=sailfish_emulator -o $APPNAME
            if [ ! -f $APPNAME ]; then echo "Failed to compile."; fi
            ;;
        rpm)
            rm -f $APPNAME
            sb2 -O use-global-tmp -t SailfishOSgo-armv7hl -m sdk-build env \
            CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 \
            CGO_CFLAGS_ALLOW='.*' \
            CGO_CXXFLAGS_ALLOW='.*' \
            CGO_LDFLAGS_ALLOW='.*' \
            CGO_FFLAGS_ALLOW='.*' \
            CPATH=/srv/mer/targets/SailfishOS-2.1.4.13-armv7hl/usr/include:/srv/mer/targets/SailfishOS-2.1.4.13-armv7hl/usr/include/qt5/QtCore:/srv/mer/targets/SailfishOS-2.1.4.13-armv7hl/usr/include/qt5:/srv/mer/targets/SailfishOS-2.1.4.13-armv7hl/usr/include/qt5/QtGui:/srv/mer/targets/SailfishOS-2.1.4.13-armv7hl/usr/include/qt5/QtNetwork:/srv/mer/targets/SailfishOS-2.1.4.13-armv7hl/usr/include/qt5/QtQuick:/srv/mer/targets/SailfishOS-2.1.4.13-armv7hl/usr/include/qt5/QtQml:/srv/mer/targets/SailfishOS-2.1.4.13-armv7hl/usr/include/sailfishapp \
            LIBRARY_PATH=/srv/mer/targets/SailfishOS-2.1.4.13-armv7hl/usr/lib:/srv/mer/targets/SailfishOS-2.1.4.13-armv7hl/lib \
            CGO_LDFLAGS="--sysroot=/srv/mer/targets/SailfishOS-2.1.4.13-armv7hl/" \
            CGO_CFLAGS="--sysroot=/srv/mer/targets/SailfishOS-2.1.4.13-armv7hl/" \
            /home/src1/projects/goroot/go/bin/linux_arm/go build -x \
            -ldflags="-v -s -w -X main.Version=$VERSION" -tags="minimal sailfish" \
            -installsuffix=sailfish -o $APPNAME
            if [ ! -f $APPNAME ]; then
                echo "Failed to compile."
            else
                mb2 -x -t SailfishOS-2.1.4.13-armv7hl build
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
            mb2 -x -s rpm/$APPNAME.spec -d "Sailfish OS Emulator" deploy  --sdk
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
