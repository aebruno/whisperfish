import QtQuick 2.0
import Sailfish.Silica 1.0

CoverBackground {
	Column {
		anchors.centerIn: parent
		width: parent.width
		spacing: Theme.paddingMedium

		Image {
			anchors.horizontalCenter: parent.horizontalCenter
			source: "/usr/share/icons/hicolor/86x86/apps/harbour-whisperfish.png"
		}

		Label {
			id: coverdata
			anchors.horizontalCenter: parent.horizontalCenter
			color: Theme.highlightColor
			font.pixelSize: Theme.fontSizeLarge
			text:  qsTr("Whisperfish")
		}
	}
}
