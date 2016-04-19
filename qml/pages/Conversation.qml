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

        PullDownMenu {
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

        Component.onCompleted: {
            messageView.positionViewAtEnd()
        }
    }
}
