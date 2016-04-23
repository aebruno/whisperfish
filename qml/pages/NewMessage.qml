import QtQuick 2.0
import Sailfish.Silica 1.0
import Sailfish.Contacts 1.0
import org.nemomobile.contacts 1.0
import org.nemomobile.commhistory 1.0

// This was adapted from jolla-messages

Page {
    id: newMessagePage
    property Label errorLabel

    _clickablePageIndicators: !(isLandscape && recipientField.activeFocus)

    onStatusChanged: {
        if (status === PageStatus.Active) {
            recipientField.forceActiveFocus()
        }
    }

    SilicaFlickable {
        id: messages
        focus: true
        contentHeight: content.y + content.height
        anchors.fill: parent

        RemorsePopup { id: remorse }

        Column {
            id: content
            y: newMessagePage.isLandscape ? Theme.paddingMedium : 0
            width: messages.width
            Item {
                width: messages.width
                height: Math.max(recipientHeader.height + (errorLabel.visible ? Theme.paddingLarge + errorLabel.height : 0), messages.height - textInput.height - content.y)

                Column {
                    id: recipientHeader
                    width: parent.width
                    PageHeader {
                        //% "New message"
                        title: "New Message"
                        visible: newMessagePage.isPortrait
                    }
                    RecipientField {
                        id: recipientField
                        property bool hasValidContact
                        property QtObject recipient
                        width: parent.width
                        requiredProperty: PeopleModel.PhoneNumberRequired
                        showLabel: newMessagePage.isPortrait
                        multipleAllowed: false

                        onEmptyChanged: if (empty) errorLabel.text = ""

                        function updateContact() {
                            for (var i = 0; i < selectedContacts.count; i++) {
                                var contact = selectedContacts.get(i)
                                if (contact.property !== undefined && contact.propertyType === "phoneNumber") {
                                    console.log("PHONE: "+contact.property.number)
                                    var c = contactsModel.find(contact.property.number)
                                    if(c.name.length != 0){
                                        hasValidContact = true
                                        recipient = c
                                        textInput.contactName = c.name
                                    } else {
                                        hasValidContact = false
                                        errorLabel.text = "Invalid recipient"
                                        remorse.execute("Contact not registered with Signal!")
                                        return
                                    }
                                } else {
                                    continue
                                }
                            }
                        }

                        //: A single recipient
                        //% "recipient"
                        placeholderText: qsTr("Recipient")

                        //: Summary of all selected recipients, e.g. "Bob, Jane, 75553243"
                        //% "Recipients"
                        summaryPlaceholderText: selectedContacts.count > 0 && recipient ? recipient.tel : ""

                        onFinishedEditing: {
                            textInput.forceActiveFocus()
                        }

                        onSelectionChanged: {
                            updateContact()
                        }
                    }
                }
                Label {
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
                enabled: recipientField && !recipientField.empty
                clearAfterSend: recipientField.hasValidContact

                onSendMessage: {
                    if (recipientField.hasValidContact) {
                        whisperfish.sendMessage(recipientField.recipient.tel, text)
                        whisperfish.refreshConversation()
                        pageStack.replaceAbove(pageStack.previousPage(), Qt.resolvedUrl("../pages/Conversation.qml"));
                    } else {
                        //: Invalid recipient error
                        //% "Invalid recipient"
                        errorLabel.text = qsTrId("Invalid recipient")
                        remorse.execute("Contact not registered with Signal!")
                    }
                }
            }
        }
        VerticalScrollDecorator {}
    }
}
