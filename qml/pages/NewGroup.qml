import QtQuick 2.2
import Sailfish.Silica 1.0

Page {
    id: newGroupPage
    property Label errorLabel
    property int selectedContacts
    property var recipients: new Object()

    _clickablePageIndicators: !(isLandscape && recipientField.activeFocus)

    SilicaFlickable {
        id: newGroup
        focus: true
        contentHeight: content.y + content.height
        anchors.fill: parent

        Column {
            id: content
            y: newGroupPage.isLandscape ? Theme.paddingMedium : 0
            width: newGroup.width
            Item {
                width: newGroup.width
                height: Math.max(recipientField.height+groupName.height, newGroup.height - textInput.height - content.y)

                Column {
                    id: recipientHeader
                    width: parent.width
                    PageHeader {
                        //: New group page title
                        //% "New Group"
                        title: qsTrId("whisperfish-new-group-title")
                        visible: newGroupPage.isPortrait
                    }

                    TextField {
                        id: groupName
                        width: parent.width
                        //: Group name label
                        //% "Group Name"
                        label: qsTrId("whisperfish-group-name-label")
                        //: Group name placeholder
                        //% "Group Name"
                        placeholderText: qsTrId("whisperfish-group-name-placeholder")
                        placeholderColor: Theme.highlightColor
                        horizontalAlignment: TextInput.AlignLeft
                     }

                    ValueButton {
                        id: recipientField
                        //: New group message members label
                        //% "Members"
                        label: qsTrId("whisperfish-new-group-message-members")
                        value: {
                            if (selectedContacts > 0) {
                                var numbers = Object.keys(recipients)
                                return numbers.map(function(v) { return recipients[v]; }).join(",")
                            } else {
                                //: New message recipient select default label
                                //% "Select"
                                return qsTrId("whisperfish-new-message-recipient-select-default")
                            }
                        }
                        onClicked: {
                            contactList.refresh()
                            var c = pageStack.push(Qt.resolvedUrl("SelectGroupContact.qml"), {contactList: contactList})
                            c.selected.connect(function(contacts) {
                                console.log("Selected contacts: "+Object.keys(contacts).length)
                                selectedContacts = Object.keys(contacts).length
                                recipients = contacts
                            })
                        }
                    }
                }
                ErrorLabel {
                    id: errorLabel
                    visible: text.length > 0
                    anchors {
                        bottom: parent.bottom
                        bottomMargin: -Theme.paddingSmall
                    }
                }
            }

            ChatTextInput {
                id: textInput
                width: parent.width
                enabled: selectedContacts > 0 && groupName.text != ""
                clearAfterSend: selectedContacts > 0 && groupName.text != ""

                onSendMessage: {
                    if (selectedContacts == 0) {
                        //: Invalid recipient error
                        //% "Please select group members"
                        errorLabel.text = qsTrId("whisperfish-error-invalid-group-members")
                    } else if(groupName.text == "") {
                        //: Invalid group name error
                        //% "Please name the group"
                        errorLabel.text = qsTrId("whisperfish-error-invalid-group-name")
                    } else {
                        var source = Object.keys(recipients).join(",")
                        var sid = MessageModel.createMessage(source, text, groupName.text, "", false)
                        if(sid > 0) {
                            pageStack.replaceAbove(pageStack.previousPage(), Qt.resolvedUrl("../pages/Conversation.qml"));
                            MessageModel.load(sid, ContactModel.name(source))
                            SessionModel.add(sid, true)
                        } else {
                            //: Failed to create message
                            //% "Failed to create message"
                            errorLabel.text = qsTrId("whisperfish-error-message-create")
                        }
                    }
                }
            }
        }
        VerticalScrollDecorator {}
    }
}
