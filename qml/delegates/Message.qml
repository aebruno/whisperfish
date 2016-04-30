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
    id: message
    contentHeight: Math.max(timestamp.y + timestamp.height, retryIcon.height) + Theme.paddingMedium
    menu: messageContextMenu

    property QtObject msg: messageModel.get(index)
    property bool inbound: msg.outgoing ? false : true
    property bool hasAttachments: msg.hasAttachment
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
        id: attachments
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
            model: msg.hasAttachment ? 1 : 0

            Thumbnail {
                id: attachment
                opacity: 1.0
                width: opacity == 1.0 ? size : 0
                height: width
                sourceSize {
                    width: size * 2
                    height: size * 2
                }

                property int size: Theme.itemSizeLarge
                property bool highlighted
                property bool isThumbnail: msg.mimeType.substr(0, 6) === "image/"
                property bool isVCard: {
                    var type = msg.mimeType.toLowerCase()
                    return type.substr(0, 10) === "text/vcard" || type.substr(0, 12) === "text/x-vcard"
                }

                source: isThumbnail ? msg.attachment : ""

                Image {
                    id: icon
                    anchors.fill: parent
                    fillMode: Image.Pad
                    source: iconSource()

                    function iconSource() {
                        if (msg === undefined ||
                            msg.mimeType.substr(0, 16) === "application/smil" ||
                            msg.mimeType.substr(1, 10) === "text/plain")
                            return ""
                        else if (isThumbnail && attachment.status !== Thumbnail.Error)
                            return ""
                        else if (isVCard)
                            return "image://theme/icon-m-person" + (highlighted ? "?" + Theme.highlightColor : "")
                        else
                            return "image://theme/icon-m-attach" + (highlighted ? "?" + Theme.highlightColor : "")
                    }

                    Rectangle {
                        anchors.fill: parent
                        z: -1
                        color: Theme.highlightColor
                        opacity: 0.1
                        visible: true
                    }
                }
            }
        }
    }

    LinkedText {
        id: messageText
        anchors {
            left: inbound ? parent.left : attachments.right
            right: inbound ? attachments.left : parent.right
            leftMargin: inbound ? sidePadding : (attachments.height ? Theme.paddingMedium : (retryIcon.width ? Theme.paddingMedium : Theme.horizontalPageMargin))
            rightMargin: !inbound ? sidePadding : (attachments.height ? Theme.paddingMedium : (retryIcon.width ? Theme.paddingMedium : Theme.horizontalPageMargin))
        }

        property int sidePadding: Theme.itemSizeSmall + Theme.horizontalPageMargin
        y: Theme.paddingMedium / 2
        height: Math.max(implicitHeight, attachments.height)
        wrapMode: Text.Wrap

        plainText: {
            if (!msg) {
                hasText = false
                return ""
            } else if (msg.message != "") {
                hasText = true
                return msg.message
            } else if (msg.hasAttachment) {
                hasText = false
                return qsTr("Multimedia Message")
            } else {
                hasText = false
                return ""
            }
        }

        color: (message.highlighted || !inbound) ? Theme.highlightColor : Theme.primaryColor
        font.pixelSize: inbound ? Theme.fontSizeMedium : Theme.fontSizeSmall
        horizontalAlignment: inbound ? Qt.AlignRight : Qt.AlignLeft
        verticalAlignment: Qt.AlignBottom
    }

    Label {
        id: timestamp
        anchors {
            left: parent.left
            leftMargin: Theme.horizontalPageMargin
            right: parent.right
            rightMargin: Theme.horizontalPageMargin
            top: messageText.bottom
            topMargin: Theme.paddingSmall
        }

        color: messageText.color
        opacity: 0.6
        font.pixelSize: Theme.fontSizeExtraSmall
        horizontalAlignment: messageText.horizontalAlignment

        text: {
            if (!msg) {
                return ""
            } else {
                var re = msg.date
                if (msg.received) {
                    re += " | " + qsTrId("Received")
                } else if (msg.sent) {
                    re += " | " + qsTr("Sent")
                }
                return re
            }
        }
    }

    onClicked: {
        if (msg.hasAttachment && attachments.height > 0) {
            pageStack.push(Qt.resolvedUrl("../pages/AttachmentPage.qml"), { 'source': msg.attachment, 'message': msg })
        }
    }
}

