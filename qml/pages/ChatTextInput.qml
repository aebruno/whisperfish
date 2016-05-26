/*
 * Adpoted for use with Whisperfish 
 *
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
import org.nemomobile.contacts 1.0
import org.nemomobile.commhistory 1.0
import Sailfish.Contacts 1.0

InverseMouseArea {
    id: chatInputArea

    // Can't use textField height due to excessive implicit padding
    height: timestamp.y + timestamp.height + Theme.paddingMedium

    property string contactName: ""
    property string attachmentPath: ""
    property alias text: textField.text
    property alias cursorPosition: textField.cursorPosition
    property alias editorFocus: textField.focus
    property bool enabled: true
    property bool clearAfterSend: true

    signal sendMessage(string text, string path)

    function setAttachmentPath(path) {
        attachmentPath = path
    }

    function send() {
        Qt.inputMethod.commit()
        if (text.length < 1 && attachmentPath.length < 1)
            return
        sendMessage(text, attachmentPath)
        if (clearAfterSend) {
            text = ""
            attachmentPath = ""
        }
        // Reset keyboard state
        if (textField.focus) {
            textField.focus = false
            textField.focus = true
        }
    }

    function forceActiveFocus() {
        textField.forceActiveFocus()
    }

    function reset() {
        Qt.inputMethod.commit()
        text = ""
    }

    property Page page: _findPage()
    function _findPage() {
        var parentItem = parent
        while (parentItem) {
            if (parentItem.hasOwnProperty('__silica_page')) {
                return parentItem
            }
            parentItem = parentItem.parent
        }
        return null
    }

    property bool onScreen: visible && Qt.application.active && page !== null && page.status === PageStatus.Active

    TextArea {
        id: textField
        anchors {
            left: parent.left
            right: sendButtonArea.left
            top: parent.top
            topMargin: Theme.paddingMedium
        }

        focusOutBehavior: FocusBehavior.KeepFocus
        textRightMargin: 0
        font.pixelSize: Theme.fontSizeSmall

        property bool empty: text.length === 0 && !inputMethodComposing

        placeholderText: contactName.length ?
        //: Personalized placeholder for chat input, e.g. "Hi John"
        //% "Hi %1"
                         qsTrId("Hi %1").arg(contactName) :
        //: Generic placeholder for chat input
        //% "Hi"
                         qsTrId("Hi")
    }

    onClickedOutside: textField.focus = false

    MouseArea {
        id: sendButtonArea
        anchors {
            fill: sendButtonText
            margins: -Theme.paddingLarge
        }
        enabled: !textField.empty && chatInputArea.enabled
        onClicked: chatInputArea.send()
    }

    IconButton {
        id: sendButtonText
        icon.source: "/usr/share/harbour-whisperfish/icons/ic_send_push_white_24dp.png"
        icon.width: Theme.iconSizeMedium
        icon.height: Theme.iconSizeMedium
        anchors {
            right: parent.right
            rightMargin: Theme.horizontalPageMargin
            verticalCenter: textField.top
            verticalCenterOffset: textField.textVerticalCenterOffset + (textField._editor.height - (height/2))
        }
        onClicked: chatInputArea.send()
        visible: true
        onPressAndHold: {
            //Workaround for rpm validator
            chatInputArea.attachmentPath = ""
            fileModel.searchPath = "foo"
            pageStack.push(imagepicker)
            imagepicker.selected.connect(chatInputArea.setAttachmentPath)
        }

        //% "Send"
    }

    Label {
        id: timestamp
        anchors {
            top: textField.bottom
            // Spacing underneath separator in TextArea is _labelItem.height + Theme.paddingSmall + 3
            topMargin: -textField._labelItem.height - 3
            left: textField.left
            leftMargin: Theme.horizontalPageMargin
            right: textField.right
        }

        color: Theme.highlightColor
        font.pixelSize: Theme.fontSizeTiny

        function updateTimestamp() {
            var date = new Date()
            text = Format.formatDate(date, Formatter.TimepointRelative)
            updater.interval = (60 - date.getSeconds() + 1) * 1000
        }

        Timer {
            id: updater
            repeat: true
            triggeredOnStart: true
            running: Qt.application.active && timestamp.visible
            onTriggered: timestamp.updateTimestamp()
        }
    }

    Label {
        id: messageType
        anchors {
            right: parent.right
            rightMargin: Theme.horizontalPageMargin
            top: timestamp.top
        }

        color: Theme.highlightColor
        font.pixelSize: Theme.fontSizeTiny
        horizontalAlignment: Qt.AlignRight
        text: attachmentPath.length == 0 ? "" : "(1) Attachment" 
    }

}
