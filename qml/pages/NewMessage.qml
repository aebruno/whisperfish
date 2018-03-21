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

import QtQuick 2.2
import Sailfish.Silica 1.0

Page {
    id: newMessagePage
    property Label errorLabel
    property string recipientNumber
    property string recipientName

    _clickablePageIndicators: !(isLandscape && recipientField.activeFocus)

    SilicaFlickable {
        id: newMessage
        focus: true
        contentHeight: content.y + content.height
        anchors.fill: parent

        RemorsePopup { id: remorse }

        Column {
            id: content
            y: newMessagePage.isLandscape ? Theme.paddingMedium : 0
            width: newMessage.width
            Item {
                width: newMessage.width
                height: Math.max(recipientField.height, newMessage.height - textInput.height - content.y)

                Column {
                    id: recipientHeader
                    width: parent.width
                    PageHeader {
                        //: New message page title
                        //% "New message"
                        title: qsTrId("whisperfish-new-message-title")
                        visible: newMessagePage.isPortrait
                    }

                    ValueButton {
                        id: recipientField
                        //: New message recipient label
                        //% "Recipient"
                        label: qsTrId("whisperfish-new-message-recipient")
                        value: {
                            if (recipientName != "") {
                                return recipientName
                            } else if (recipientNumber != "") {
                                return recipientNumber
                            } else {
                                //: New message recipient select default label
                                //% "Select"
                                return qsTrId("whisperfish-new-message-recipient-select-default")
                            }
                        }
                        onClicked: {
                            contactList.refresh()
                            var c = pageStack.push(Qt.resolvedUrl("SelectContact.qml"), {contactList: contactList})
                            c.selected.connect(function(name, tel) {
                                console.log("Selected contact: "+name+' '+tel)
                                recipientNumber = tel
                                recipientName = name
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
                enabled: recipientNumber.length != 0
                clearAfterSend: recipientNumber.length != 0

                onSendMessage: {
                    if (recipientNumber.length != 0) {
                        var source = recipientNumber
                        var sid = MessageModel.createMessage(source, text, "", "", false)
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
