/*
 * Copyright (C) 2012-2015 Jolla Ltd.
 *
 * The code in this file is distributed under multiple licenses, and as such,
 * may be used under any one of the following licenses:
 *
 *   - GNU General Public License as published by the Free Software Foundation;
 *     either version 2 of the License (see LICENSE.GPLv2 in the root directory
 *     for full terms), or (at your option) any later version.
 *   - GNU Lesser General Public License as published by the Free Software
 *     Foundation; either version 2.1 of the License (see LICENSE.LGPLv21 in the
 *     root directory for full terms), or (at your option) any later version.
 *   - Alternatively, if you have a commercial license agreement with Jolla Ltd,
 *     you may use the code under the terms of that license instead.
 *
 * You can visit <https://sailfishos.org/legal/> for more information
 */

import QtQuick 2.0
import Sailfish.Silica 1.0
import Sailfish.Gallery 1.0
import Sailfish.TransferEngine 1.0
import Sailfish.TextLinking 1.0

SplitViewPage {
    id: root
    property alias source: viewer.source
    property var message

    signal copy()

    allowedOrientations: Orientation.All
    backNavigation: true
    open: true

    // This is the share method list, but it also
    // includes the pulley menu
    background: ShareMethodList {
        id: menuList

        source: root.source
        anchors.fill: parent
        filter: root.message.mimeType
        content: QtObject { property string type: root.message.mimeType }

        PullDownMenu {
            id: pullDownMenu
            MenuItem {
                //% "Copy to gallery"
                text: qsTr("Copy to gallery")
                onClicked: root.copy()
            }
        }


        header: PageHeader {
            function msgDate() {
                var dt = new Date(root.message.timestamp)
                return Format.formatDate(dt, Formatter.Timepoint)
            }
            title: root.message.outgoing ? qsTr("Me") : MessageModel.peerName
            description: msgDate()
        }
    }

    ImageViewer {
        id: viewer
        anchors.fill: parent
        enableZoom: !root.open
        fit: root.isPortrait ? Fit.Width : Fit.Height
        onClicked: root.open = !root.open
    }
}
