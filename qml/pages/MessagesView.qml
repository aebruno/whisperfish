/*
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
import org.nemomobile.commhistory 1.0
import "../delegates"

SilicaListView {
    id: messagesView

    verticalLayoutDirection: ListView.BottomToTop
    // Necessary to avoid resetting focus every time a row is added, which breaks text input
    currentIndex: -1
    quickScroll: true

    delegate: Item {
        id: wrapper

        // This would normally be previousSection, but our model's order is inverted.
        property bool sectionBoundary: (ListView.nextSection != "" && ListView.nextSection !== ListView.section)
                                        || model.index === messagesView.count - 1
        property Item section

        height: loader.y + loader.height
        width: parent.width

        ListView.onRemove: loader.item.animateRemoval(wrapper)

        Loader {
            id: loader
            y: section ? section.y + section.height : 0
            width: parent.width
            sourceComponent: messageDelegate
        }

        onSectionBoundaryChanged: {
            if (sectionBoundary) {
                section = sectionHeader.createObject(wrapper, { 'modelData': model })
            } else {
                section.destroy()
                section = null
            }
        }

        Component {
            id: messageDelegate

            Message { 
                modelData: model
            }
        }
    }

    section.property: "localUid"

    Component {
        id: sectionHeader

        Row {
            id: header
            y: Theme.paddingMedium
            x: parent ? (parent.width - width) / 2 : 0
            height: text.implicitHeight + Theme.paddingSmall
            spacing: Theme.paddingMedium

            Label {
                id: text
                color: Theme.highlightColor
                font.pixelSize: Theme.fontSizeExtraSmall
                text: MessageModel.group ? 
                    //: Group message label
                    //% "Group: %1"
                    qsTrId("whisperfish-group-label").arg(MessageModel.peerName) : 
                    MessageModel.peerName
            }
        }
    }

    function remove(contentItem) {
        //: Deleteing message remorse
        //% "Deleteing"
        contentItem.remorseAction(qsTrId("whisperfish-delete-message"),
            function() {
                console.log("Delete message: "+contentItem.modelData.id)
                MessageModel.remove(contentItem.modelData.index)
            })
    }

    function resend(contentItem) {
        //: Resend message remorse
        //% "Resending"
        contentItem.remorseAction(qsTrId("whisperfish-resend-message"),
            function() {
                console.log("Resending message: "+contentItem.modelData.id)
                MessageModel.sendMessage(contentItem.modelData.id)
            })
    }

    function copy(contentItem) {
        MessageModel.copyToClipboard(contentItem.modelData.message)
    }

    Component {
        id: messageContextMenu

        ContextMenu {
            id: menu

            width: parent ? parent.width : Screen.width

            MenuItem {
                visible: menu.parent && menu.parent.hasText
                //: Copy message menu item
                //% "Copy"
                text: qsTrId("whisperfish-copy-message-menu")
                onClicked: copy(menu.parent)
            }
            MenuItem {
                //: Delete message menu item
                //% "Delete"
                text: qsTrId("whisperfish-delete-message-menu")
                onClicked: remove(menu.parent)
            }
            MenuItem {
                //: Resend message menu item
                //% "Resend"
                text: qsTrId("whisperfish-resend-message-menu")
                visible: menu.parent && menu.parent.modelData.queued
                onClicked: resend(menu.parent)
            }
        }
    }


    RemorsePopup { id: remorse }

    VerticalScrollDecorator {}
}

