import QtQuick 2.0
import Sailfish.Silica 1.0

Page {
    id: root

    SilicaFlickable {
        anchors.fill: parent
        contentHeight: column.height + Theme.paddingLarge
        contentWidth: parent.width

        PullDownMenu {
            MenuItem {
                text: qsTr("About Whisperfish")
                onClicked: pageStack.push(Qt.resolvedUrl("About.qml"))
            }
        }

        VerticalScrollDecorator {}

		Column {
			id: col
			spacing: Theme.paddingLarge
			width: parent.width
			PageHeader {
				title: qsTr("Hello World")
			}

            Label {
                anchors.horizontalCenter: parent.horizontalCenter
                font.bold: true
                text: qsTr("Whisperfish")
            }
		}
    }
}
