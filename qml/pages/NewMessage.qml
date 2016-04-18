import QtQuick 2.0
import Sailfish.Silica 1.0

Page {
    id: newMessage
    objectName: "newMessage"

	SilicaFlickable {
		anchors.fill: parent
		contentWidth: parent.width
		contentHeight: col.height + Theme.paddingLarge

        Column {
			id: col
            width: parent.width

            PageHeader { title: "New Message" }

            TextArea {
                width: parent.width
                height: Math.max(newMessage.width/3, implicitHeight)
                placeholderText: "Type multi-line text here"
                label: "Expanding text area"
            }
        }
    }
}
