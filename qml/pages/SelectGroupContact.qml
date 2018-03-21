import QtQuick 2.2
import Sailfish.Silica 1.0

Dialog {
    id: page
    objectName: "selectGroupContact"
    canAccept: selectedContacts > 0
    allowedOrientations: Orientation.All

    property int selectedContacts: 0
    property var recipients: new Object()
    property alias contactList: alphaMenu.dataSource
    signal selected(var recipients)

    onDone: {
        if (result == DialogResult.Accepted && selectedContacts > 0) {
            page.selected(recipients)
        }
    }

    Column {
        width: parent.width
        spacing: Theme.paddingLarge

        DialogHeader {
            id: title
            //: Title for select group contact page
            //% "Select group members"
            title: qsTrId("whisperfish-select-group-contact")
            //: placeholder showing selected group contacts
            //% "Selected %1"
            acceptText: selectedContacts > 0 ? qsTrId("whisperfish-select-group-num-contacts").arg(selectedContacts) : ""
        }

        AlphaMenu {
            id: alphaMenu
            dataSource: ListModel{}
            listDelegate:  BackgroundItem {
                id: contactItem
                width: parent.width
                highlighted: tel in recipients ? true : false
                onClicked: {
                    if(tel in recipients) {
                         delete recipients[tel]
                         page.selectedContacts = Object.keys(recipients).length
                         highlighted = false
                    } else {
                         recipients[tel] = name
                         page.selectedContacts = Object.keys(recipients).length
                         highlighted = true
                    }
                }
                Row {
                    spacing: 20

                    Column {
                        Label {
                            text: name
                            font.pixelSize: Theme.fontSizeMedium
                            color: Theme.primaryColor
                        }
                        Label {
                            text: tel
                            font.pixelSize: Theme.fontSizeExtraSmall
                            color: Theme.secondaryColor
                        }
                    }
                }
            }
        }
    }
}
