package libtox

// OnSelfConnectionStatusChanges This event is triggered whenever there is a change in the DHT connectionstate.
/*
 * When disconnected, a client may choose to call tox_bootstrap again, to reconnect to the DHT.
 * Note that this state may frequently change for short amounts of time. Clients should therefore not immediately bootstrap onreceiving a disconnect.
 */
type OnSelfConnectionStatusChanges func(tox *Tox, status ToxConnection)

// OnFriendNameChanges This event is triggered when a friend changes their name.
type OnFriendNameChanges func(tox *Tox, friendnumber uint32, name []byte, length uint32)

// OnFriendStatusMessageChanges This event is triggered when a friend changes their status message.
type OnFriendStatusMessageChanges func(tox *Tox, friendnumber uint32, message []byte, length uint32)

// OnFriendStatusChanges This event is triggered when a friend changes their status message.
type OnFriendStatusChanges func(tox *Tox, friendnumber uint32, userstatus ToxUserStatus)

// OnFriendConnectionStatusChanges This event is triggered when a friend goes offline after having been online, or when a friend goes online.
// This callback is not called when adding friends. It is assumed that when adding friends, their connection status is initially offline.
type OnFriendConnectionStatusChanges func(tox *Tox, friendnumber uint32, connectionstatus ToxConnection)

// OnFriendTypingChanges This event is triggered when a friend starts or stops typing.
type OnFriendTypingChanges func(tox *Tox, friendnumber uint32, istyping bool)

// OnFriendReadReceipt This event is triggered when the friend receives the message with the corresponding message ID. */
type OnFriendReadReceipt func(tox *Tox, friendnumber uint32, messageid uint32)

// OnFriendRequest This event is triggered when a friend request is received.
type OnFriendRequest func(tox *Tox, publickey []byte, message []byte, length uint32)

// OnFriendMessage This event is triggered when a message from a friend is received.
type OnFriendMessage func(tox *Tox, friendnumber uint32, messagetype ToxMessageType, message []byte, length uint32)

// OnFileRecvControl This event is triggered when a file control command is received from a friend.
type OnFileRecvControl func(tox *Tox, friendnumber uint32, filenumber uint32, filecontrol ToxFileControl)

// OnFileChunkRequest This event is triggered when Core is ready to send more file data.
type OnFileChunkRequest func(tox *Tox, friendnumber uint32, filenumber uint32, position uint64, length uint64)

// OnFileRecv This event is triggered when a file transfer request is received.
type OnFileRecv func(tox *Tox, friendnumber uint32, filenumber uint32, kind ToxFileKind, filesize uint64, filename string, length uint32)

// OnFileRecvChunk This event is first triggered when a file transfer request is received, and subsequently when a chunk of file data for an accepted request was received.
type OnFileRecvChunk func(tox *Tox, friendnumber uint32, filenumber uint32, position uint64, data []byte, length uint32)

// OnFriendLossyPacket This event is triggered when a lossy packet is received from a friend.
type OnFriendLossyPacket func(tox *Tox, friendnumber uint32, data []byte, length uint32)

// OnFriendLosslessPacket This event is triggered when a lossless packet is received from a friend.
type OnFriendLosslessPacket func(tox *Tox, friendnumber uint32, data []byte, length uint32)

/*Conference callbacks*/

// OnConferenceInvite This event is triggered when the client is invited to join a conference.
type OnConferenceInvite func(tox *Tox, friendnumber uint32, conferencetype ToxConferenceType, cookie []byte, length uint32)

// OnConferenceConnected This event is triggered when the client successfully connects to a conference after joining it with the tox_conference_join function.
type OnConferenceConnected func(tox *Tox, conferencenumber uint32)

// OnConferenceMessage This event is triggered when the client receives a conference message.
type OnConferenceMessage func(tox *Tox, conferencenumber uint32, peernumber uint32, messagetype ToxMessageType, message []byte, length uint32)
