import QtQuick 2.0
import Sailfish.Silica 1.0

Page {
	id: aboutpage
	SilicaFlickable {
		anchors.fill: parent
		contentWidth: parent.width
		contentHeight: col.height + Theme.paddingLarge

		VerticalScrollDecorator {}

		Column {
			id: col
			spacing: Theme.paddingLarge
			width: parent.width
			PageHeader {
				title: qsTr("About Whisperfish")
			}

			Image {
				anchors.horizontalCenter: parent.horizontalCenter
				source: "/usr/share/icons/hicolor/86x86/apps/harbour-whisperfish.png"
			}

            Label {
                anchors.horizontalCenter: parent.horizontalCenter
                font.bold: true
                text: qsTr("Whisperfish v") + Qt.application.version
            }

            TextArea {
                anchors.horizontalCenter: parent.horizontalCenter
                width: parent.width
                horizontalAlignment: TextEdit.Center
                readOnly: true
                text: qsTr("Signal client for Sailfish OS")
            }

            TextArea {
                anchors.horizontalCenter: parent.horizontalCenter
                width: parent.width
                horizontalAlignment: TextEdit.Center
                readOnly: true
                text: qsTr("Copyright: Andrew E. Bruno\nLicense: GPLv3")
            }

            Button {
                anchors.horizontalCenter: parent.horizontalCenter
                text: qsTr("Source Code")
                onClicked: {
                    Qt.openUrlExternally("https://github.com/aebruno/whisperfish")
                }
            }

            Button {
                anchors.horizontalCenter: parent.horizontalCenter
                text: qsTr("Report a Bug")
                onClicked: {
                    Qt.openUrlExternally("https://github.com/aebruno/whisperfish/issues")
                }
            }

            SectionHeader {
                text: qsTr("Additional Copyright")
            }

            Label {
                text: qsTr("Signal client library for Go (C) Jani Monoses.")
                anchors.horizontalCenter: parent.horizontalCenter
                wrapMode: Text.WrapAtWordBoundaryOrAnywhere
                width: (parent ? parent.width : Screen.width) - Theme.paddingLarge * 2
                verticalAlignment: Text.AlignVCenter
                horizontalAlignment: Text.AlignLeft
                x: Theme.paddingLarge
            }
		}
	}
}
