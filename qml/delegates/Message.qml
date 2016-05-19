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
import Sailfish.TextLinking 1.0
import org.nemomobile.thumbnailer 1.0

ListItem {
    id: messageItem
    contentHeight: Math.max(timestampLabel.y + timestampLabel.height, retryIcon.height) + Theme.paddingMedium
    menu: messageContextMenu

    property QtObject modelData
    property bool inbound: modelData.outgoing ? false : true
    property bool hasText
    property bool canRetry

    // Retry icon for non-attachment outbound messages
    Image {
        id: retryIcon
        anchors {
            left: inbound ? undefined : parent.left
            right: inbound ? parent.right : undefined
            bottom: parent.bottom
        }
    }

    Column {
        id: attachmentBox
        height: implicitHeight
        width: implicitWidth
        anchors {
            left: inbound ? undefined : parent.left
            right: inbound ? parent.right : undefined
            // We really want the baseline of the last line of text, but there's no way to get that
            bottom: messageText.bottom
            bottomMargin: messageText.y
        }

        Repeater {
            id: attachmentLoader
            model: modelData.hasAttachment ? 1 : 0
            property QtObject attachmentItem: modelData

            Attachment {
                anchors.right: inbound ? parent.right : undefined
                messagePart: attachmentLoader.attachmentItem
                showRetryIcon: false
                highlighted: messageItem.highlighted
            }
        }
    }

    LinkedText {
        id: messageText
        anchors {
            left: inbound ? parent.left : attachmentBox.right
            right: inbound ? attachmentBox.left : parent.right
            leftMargin: inbound ? sidePadding : (attachmentBox.height ? Theme.paddingMedium : (retryIcon.width ? Theme.paddingMedium : Theme.horizontalPageMargin))
            rightMargin: !inbound ? sidePadding : (attachmentBox.height ? Theme.paddingMedium : (retryIcon.width ? Theme.paddingMedium : Theme.horizontalPageMargin))
        }

        property int sidePadding: Theme.itemSizeSmall + Theme.horizontalPageMargin
        y: Theme.paddingMedium / 2
        height: Math.max(implicitHeight, attachmentBox.height)
        wrapMode: Text.Wrap

        plainText: {
            if (modelData.message != "") {
                hasText = true
                return modelData.message
            } else {
                hasText = false
                return ""
            }
        }

        color: (messageItem.highlighted || !inbound) ? Theme.highlightColor : Theme.primaryColor
        font.pixelSize: inbound ? Theme.fontSizeMedium : Theme.fontSizeSmall
        horizontalAlignment: inbound ? Qt.AlignRight : Qt.AlignLeft
        verticalAlignment: Qt.AlignBottom
    }

    Label {
        id: timestampLabel
        anchors {
            left: parent.left
            leftMargin: Theme.horizontalPageMargin
            right: parent.right
            rightMargin: Theme.horizontalPageMargin
            top: messageText.bottom
            topMargin: Theme.paddingSmall
        }

        function msgDate() {
            var dt = new Date(modelData.timestamp)
            var md = Format.formatDate(dt, Formatter.Timepoint)
            return md
        }

        color: messageText.color
        opacity: 0.6
        font.pixelSize: Theme.fontSizeExtraSmall
        horizontalAlignment: messageText.horizontalAlignment
        wrapMode: Text.Wrap

        text: {
           var re = msgDate()
           if (modelData.received) {
               re += qsTr("  ✓✓")
           } else if (modelData.sent) {
               re += qsTr("  ✓")
           }
           if(inbound && messageModel.isGroup) {
               re += " | " + contactsModel.name(modelData.source, whisperfish.settings().countryCode)
           }
           return re
        }
    }

    onClicked: {
        if (modelData.hasAttachment && attachmentBox.height > 0) {
            pageStack.push(Qt.resolvedUrl("../pages/AttachmentPage.qml"), { 'source': modelData.attachment, 'message': modelData })
        }
    }
}

