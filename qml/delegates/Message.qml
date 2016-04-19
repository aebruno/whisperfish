import QtQuick 2.0
import Sailfish.Silica 1.0

BackgroundItem {
    id: listItem
    height: itemId.height + 10

    property QtObject msg: messageModel.get(index)

    ListItem {
        id: itemId
        menu: ContextMenu {
            MenuItem {
                text: qsTr("Copy to clipboard")
                onClicked: Clipboard.text = msg.message
            }
            MenuItem {
                text: qsTr("Delete")
                onClicked: console.log("TODO: implement me")
            }
        }
        width: parent.width
        contentHeight: msgLabel.paintedHeight + dateLabel.paintedHeight
        Column {
            width: parent.width * 0.75
            anchors.left: msg.sent ? parent.right : undefined
            anchors.right: !msg.sent ? parent.right : undefined
            anchors.leftMargin: Theme.paddingLarge
            anchors.rightMargin: Theme.paddingLarge

            Label {
                id: msgLabel
                width: parent.width
                text: msg.message
                wrapMode: Text.Wrap
                elide: Text.ElideRight
                truncationMode: TruncationMode.Fade
                horizontalAlignment: msg.sent ? Text.AlignLeft : Text.AlignRight
                color: msg.sent ? Theme.highlightColor : Theme.primaryColor
                textFormat: Text.StyledText
                onLinkActivated: Qt.openUrlExternally(link)
            }

           Label {
                id: dateLabel
                width: parent.width
                font.pixelSize: Theme.fontSizeExtraSmall
                text: msg.date
                horizontalAlignment: msg.sent ? Text.AlignLeft : Text.AlignRight
                color: msg.sent ? Theme.secondaryHighlightColor : Theme.secondaryColor
            }

        }
    }
}
