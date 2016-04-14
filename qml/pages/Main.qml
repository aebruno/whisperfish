import QtQuick 2.0
import Sailfish.Silica 1.0

Page {
    id: main

    SilicaListView {
        id: listView
        model: contactsModel.len
        anchors.fill: parent

        PullDownMenu {
            MenuItem {
                text: qsTr("About Whisperfish")
                onClicked: pageStack.push(Qt.resolvedUrl("About.qml"))
            }
        }

        VerticalScrollDecorator {}

        delegate: BackgroundItem {
            id: delegate

            Label {
                x: Theme.paddingLarge
                text: contactsModel.contact(index).name
                anchors.verticalCenter: parent.verticalCenter
                color: delegate.highlighted ? Theme.highlightColor : Theme.primaryColor
            }
            onClicked: console.log("Clicked " + index)
        }

    }
}
