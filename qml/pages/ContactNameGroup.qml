/*
 * Author: r0kk3rz
 * https://github.com/r0kk3rz/sailfish-alphabet-sort
 */

import QtQuick 2.2
import Sailfish.Silica 1.0

Item {
    id: root

    property bool active
    property string name
    property string iconSource
    property int groupIndex
    property bool hasEntries

    property real baseHeight
    property Item groupResultsList

    property alias pressed: mouseArea.pressed
    property alias containsMouse: mouseArea.containsMouse
    property bool highlighted: pressed && containsMouse || root.active

    signal clicked(var mouse)

    MouseArea {
        // This MouseArea can't be the root item, because that item is the parent of
        // ColumnView containing contact items, and therefore can't be disabled
        id: mouseArea

        anchors.fill: parent
        enabled: root.hasEntries
        onClicked: root.clicked(mouse)

        Rectangle {
            width: parent.width
            height: root.baseHeight
            color: Theme.highlightBackgroundColor
            opacity: highlighted ? 0.1 : 0.3
        }

        Item {
            width: parent.width
            height: root.baseHeight

            opacity: mouseArea.enabled ? 1.0 : 0.3
            Behavior on opacity { NumberAnimation { duration: 500 } }

            property color textColor: highlighted ? Theme.secondaryHighlightColor : Theme.secondaryColor

            Label {
                visible: text != ''
                anchors.horizontalCenter: parent.horizontalCenter
                y: root.baseHeight/2 - implicitHeight/2
                text: root.name
                font.pixelSize: Theme.fontSizeHuge
                color: highlighted ? Theme.highlightColor : parent.textColor

                width: parent.width - Theme.paddingMedium*2
                horizontalAlignment: Text.AlignHCenter
                verticalAlignment: Text.AlignVCenter
                fontSizeMode: Text.Fit
            }

            Image {
                visible: root.iconSource != ''
                anchors.horizontalCenter: parent.horizontalCenter
                y: root.baseHeight/2 - implicitHeight/2
                source: !visible ? '' : (root.iconSource + '?' + parent.textColor)
            }
        }
    }
}
