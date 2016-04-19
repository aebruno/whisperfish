import QtQuick 2.0
import Sailfish.Silica 1.0
import "../delegates"

Page {
    id: conversation
    objectName: "conversation"

    SilicaListView {
        id: messageView
        model: messageModel.length
        anchors.fill: parent
        spacing: Theme.paddingMedium
        property int len: messageModel.length

        onLenChanged: {
            messageView.positionViewAtEnd()
        }

        PushUpMenu {
            MenuItem {
                text: qsTr("Delete All")
                onClicked: console.log("TODO: implement me")
            }
        }

        VerticalScrollDecorator {}

        ViewPlaceholder {
            enabled: messageView.count == 0
            text: "No messages"
            hintText: ""
        }

        delegate: Message{}

        footer: Row {
                width: conversation.width
                TextArea {
                    objectName: "sendTextArea"
                    id: sendBox
                    anchors.bottom: parent.bottom
                    width: parent.width - sendButton.width
                    focus: true
                    color: Theme.highlightColor
                    font.family: "cursive"
                    placeholderText: qsTr("Hi ") + messageModel.name

                    EnterKey.enabled: true
                    EnterKey.highlighted: true
                    EnterKey.iconSource: "image://theme/icon-m-enter-next"
                    EnterKey.onClicked: {
                        sendMessage()
                        sendBox.select()
                    }

                    function sendMessage() {
                        if (sendBox.text.trim()==="") {
                            console.log("Empty message")
                            sendBox.text = "";
                            return
                        }

                        whisperfish.sendMessage(messageModel.tel, sendBox.text.replace(/(\n)/gm,"").trim())
                        sendBox.deselect();
                        sendBox.text = "";
                        sendBox.placeholderText = qsTr("Sending message...");
                    }
                }
                IconButton {
                    id: sendButton
                    width: 140
                    icon.source: "/usr/share/harbour-whisperfish/icons/ic_send_push_white_24dp.png"
                    icon.width: Theme.iconSizeMedium
                    icon.height: Theme.iconSizeMedium
                    onClicked: sendBox.sendMessage()
                    onPressAndHold: {
                        console.log("TODO: implement image picker")
                    }
                 }
                function openKeyboard() {
                    sendBox.forceActiveFocus()
                }
                function openText(text) {
                    sendBox.text = text;
                }

                Connections {
                    target: messageView
                    onLenChanged: { sendBox.placeholderText = qsTr("Hi ") + messageModel.name }
                }

}

        Component.onCompleted: {
            messageView.positionViewAtEnd()
        }
    }
}
