/*
 * This was adapted from jolla-messages for use with Whisperfish
 *
 * Copyright (C) 2012-2015 Jolla Ltd.
 *
 * The code in this file is distributed under multiple licenses, and as such,
 * may be used under any one of the following licenses:
 *
 *   - GNU General Public License as published by the Free Software Foundation;
 *     either version 2 of the License (see LICENSE.GPLv2 in the root directory
 *     for full terms), or (at your option) any later version.
 *   - GNU Lesser General Public License as published by the Free Software
 *     Foundation; either version 2.1 of the License (see LICENSE.LGPLv21 in the
 *     root directory for full terms), or (at your option) any later version.
 *   - Alternatively, if you have a commercial license agreement with Jolla Ltd,
 *     you may use the code under the terms of that license instead.
 *
 * You can visit <https://sailfishos.org/legal/> for more information
 */

import QtQuick 2.0
import Sailfish.Silica 1.0
import Sailfish.Contacts 1.0
import org.nemomobile.contacts 1.0
import org.nemomobile.commhistory 1.0


Page {
    id: newMessagePage
    property Label errorLabel

    _clickablePageIndicators: !(isLandscape && recipientField.activeFocus)

    onStatusChanged: {
        if (status === PageStatus.Active) {
            recipientField.forceActiveFocus()
        }
    }

    Component {
        id: personComponent
        Person {}
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
                        //: New message page title
                        //% "New message"
                        title: qsTrId("whisperfish-new-message-title")
                        visible: newMessagePage.isPortrait
                    }
                    RecipientField {
                        id: recipientField
                        property int validContacts
                        property var recipients: new Object()
                        width: parent.width
                        requiredProperty: PeopleModel.PhoneNumberRequired
                        showLabel: newMessagePage.isPortrait
                        contactSearchModel: peopleModel
                        multipleAllowed: true

                        onEmptyChanged: if (empty) errorLabel.text = ""

                        function updateConversation() {
                            var invalidContactFound = false
                            var exists = true
                            for (var i = 0; i < selectedContacts.count; i++) {
                                var contact = selectedContacts.get(i)
                                if (contact.property !== undefined && contact.propertyType === "phoneNumber") {
                                    var tel = ContactModel.format(contact.property.number)
                                    exists = ContactModel.registered(contact.property.number)
                                    if(tel.length != 0){
                                        recipients[tel] = true
                                    } else {
                                        invalidContactFound = true
                                        var p = personComponent.createObject(null)
                                        p.resolvePhoneNumber(contact.property.number, true)
                                        if (p.id) {
                                            //: Could not format contact phone number error message
                                            //% "Error: invalid phone number for %1"
                                            errorLabel.text = qsTrId("whisperfish-error-contact-number-format").arg(p.firstName)
                                        }
                                    }
                                }
                            }

                            if(!exists && invalidContactFound == false && Object.keys(recipients).length > 0){
                                //: Could verify contact is registered with signal
                                //% "Warning: could not verify contact in Signal"
                                errorLabel.text = qsTrId("whisperfish-error-verify-contact")
                                validContacts = Object.keys(recipients).length
                            } else if(invalidContactFound == false && Object.keys(recipients).length > 0){
                                validContacts = Object.keys(recipients).length
                                errorLabel.text = ""
                            }
                        }

                        //: A single recipient
                        //% "Recipient"
                        placeholderText: qsTrId("whisperfish-recipient")

                        //: Summary of all selected recipients, e.g. "Bob, Jane, 75553243"
                        //% "Recipients"
                        summaryPlaceholderText: qsTrId("whisperfish-recipients")

                        onFinishedEditing: {
                            textInput.forceActiveFocus()
                        }

                        onSelectionChanged: {
                            updateConversation()
                        }
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
                        visible: recipientField.validContacts > 1
                        horizontalAlignment: TextInput.AlignLeft
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
                enabled: recipientField && !recipientField.empty
                clearAfterSend: recipientField.validContacts > 0

                onSendMessage: {
                    if (recipientField.validContacts > 0) {
                        var source = Object.keys(recipientField.recipients).join(",")
                        var sid = MessageModel.createMessage(source, text, groupName.text, attachmentPath, false)
                        if(sid > 0) {
                            pageStack.replaceAbove(pageStack.previousPage(), Qt.resolvedUrl("../pages/Conversation.qml"));
                            MessageModel.load(sid, ContactModel.name(source))
                            SessionModel.add(sid, true)
                        } else {
                            //: Failed to create message
                            //% "Failed to create message"
                            errorLabel.text = qsTrId("whisperfish-error-message-create")
                        }
                    } else {
                        //: Invalid recipient error
                        //% "Invalid recipient"
                        errorLabel.text = qsTrId("whisperfish-error-invalid-recipient")
                    }
                }
            }
        }
        VerticalScrollDecorator {}
    }
}
