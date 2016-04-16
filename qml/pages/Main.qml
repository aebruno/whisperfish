import QtQuick 2.0
import Sailfish.Silica 1.0

Page {
    id: main
    objectName: "main"

    property QtObject currentPage: pageStack.currentPage

    function getPhoneNumber() {
        pageStack.push(Qt.resolvedUrl("Register.qml"))
    }

    function getVerificationCode() {
        pageStack.push(Qt.resolvedUrl("Verify.qml"))
    }

    function getStoragePassword() {
        pageStack.push(Qt.resolvedUrl("Password.qml"))
    }

    function registered() {
        registeredRemorse.execute("Registration complete!", function() { console.log("Registration complete") })
    }

    RemorsePopup { id: registeredRemorse }

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

        ViewPlaceholder {
            enabled: listView.count == 0
            text: "No contacts found"
            hintText: "None of our contacts appear to be in Signal"
        }

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
