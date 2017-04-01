import QtQuick 2.0
import Sailfish.Silica 1.0
import QtMultimedia 5.0
import org.nemomobile.thumbnailer 1.0

Page {
    id: videoAttachment
    property var message

    Flickable {
        id: imageFlickable
        anchors.fill: parent
        clip: true

        PageHeader {
            function msgDate() {
                var dt = new Date(message.timestamp)
                return Format.formatDate(dt, Formatter.Timepoint)
            }
            description: msgDate()
            title: message.outgoing ?
            //: Personalized placeholder showing the attachment is from oneself
            //% "Me"
                qsTrId("whisperfish-attachment-from-self") :
            //: Personalized placeholder showing the attachment is from contact
            //% "From %1"
                qsTrId("whisperfish-attachment-from-contact").arg(MessageModel.peerName)
        }

        Item {
            property alias player: videoPlayer
            width: imageFlickable.width
            height: imageFlickable.height
            visible: videoPlayer.playbackState != MediaPlayer.StoppedState

            VideoOutput {
                id: video

                property bool playing: videoPlayer.playbackState == MediaPlayer.PlayingState

                anchors.fill: parent
                source: MediaPlayer {
                    id: videoPlayer

                    autoPlay: true
                    source: message.attachment

                    onPlaybackStateChanged: {
                        if (playbackState == MediaPlayer.PlayingState && view.menuOpen) {
                            // go fullscreen for playback if triggered via Play icon.
                            view.clicked()
                        }
                    }
                }
            }
        }
    }

    VerticalScrollDecorator { flickable: imageFlickable }
}
