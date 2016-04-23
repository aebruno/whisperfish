import QtQuick 2.0
import Sailfish.Silica 1.0
import org.nemomobile.contacts 1.0
import org.nemomobile.commhistory 1.0
import Sailfish.Contacts 1.0

// This was adopted from jolla-messages

InverseMouseArea {
    id: chatInputArea

    // Can't use textField height due to excessive implicit padding
    height: timestamp.y + timestamp.height + Theme.paddingMedium

    property string contactName: ""
    property alias text: textField.text
    property alias cursorPosition: textField.cursorPosition
    property alias editorFocus: textField.focus
    property bool enabled: true
    property bool clearAfterSend: true

    signal sendMessage(string text)

    function send() {
        Qt.inputMethod.commit()
        if (text.length < 1)
            return
        sendMessage(text)
        if (clearAfterSend) {
            text = ""
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
            verticalCenter: textField.middle
            verticalCenterOffset: textField.textVerticalCenterOffset + (textField._editor.height - height)
        }
        onClicked: chatInputArea.send()
        visible: true

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

}
